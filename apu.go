package main

const (
	HiPassStrong = 225574
	HiPassWeak   = 57593
)

var (
	SquareLookup = [][]int{
		[]int{0, 1, 0, 0, 0, 0, 0, 0},
		[]int{0, 1, 1, 0, 0, 0, 0, 0},
		[]int{0, 1, 1, 1, 1, 0, 0, 0},
		[]int{1, 0, 0, 1, 1, 1, 1, 1},
	}

	TriangleLookup = []int16{
		8, 9, 10, 11, 12, 13, 14,
		15, 15, 14, 13, 12, 11, 10,
		9, 8, 7, 6, 5, 4, 3, 2, 1,
		0, 0, 1, 2, 3, 4, 5, 6, 7,
	}

	NoiseLookup = []int{
		4, 8, 16, 32, 64, 96, 128, 160, 202,
		254, 380, 508, 762, 1016, 2034, 4068,
	}

	LengthTable = []Word{
		10, 254, 20, 2, 40, 4, 80, 6,
		160, 8, 60, 10, 14, 12, 26, 14,
		12, 16, 24, 18, 48, 20, 96, 22,
		192, 24, 72, 26, 16, 28, 32, 30,
	}
)

type Envelope struct {
	Volume       Word
	DecayRate    Word
	DecayCounter Word
	DecayEnabled bool
	LoopEnabled  bool
	Disabled     bool
	Reset        bool
}

type Square struct {
	Enabled       bool
	LengthEnabled bool
	DutyCycle     Word
	DutyCount     Word
	Timer         int
	TimerCount    int
	Length        Word
	LastTick      int
	SweepEnabled  bool
	Sweep         Word
	SweepMode     Word
	Shift         Word
	Negative      bool
	Sample        int16
	Envelope
}

type Triangle struct {
	ReloadValue   Word
	Control       bool
	Enabled       bool
	LengthEnabled bool
	Halt          bool
	Timer         int
	TimerCount    int
	Length        Word
	Counter       int
	LookupCounter int
	Sample        int16
}

type Noise struct {
	LengthEnabled bool
	Enabled       bool
	BaseEnvelope  Word
	Mode          bool
	Timer         int
	TimerCount    int
	Length        Word
	Shift         int
	Sample        float64
	Envelope
}

type Apu struct {
	DmcEnabled bool
	Square1    Square
	Square2    Square
	Triangle
	Noise
	IrqEnabled   bool
	IrqActive    bool
	HipassStrong int64
	HipassWeak   int64

	FrameCounter  int
	FrameTick     int
	LastFrameTick int

	Sample int16

	Output chan int16
}

func (s *Square) WriteControl(v Word) {
	// 76543210
	// ||||||||
	// ||||++++- Envelope
	// |||+----- Saw Envelope Disable (0: use internal counter for volume; 1: use Envelope for volume)
	// ||+------ Length Counter Disable (0: use Length Counter; 1: disable Length Counter)
	// ++------- Duty Cycle
	s.Envelope.Disabled = (v>>4)&0x1 == 1
	s.LengthEnabled = v&0x20 != 0x20
	s.Envelope.LoopEnabled = v&0x20 == 0x20
	s.DutyCycle = (v >> 6) & 0x3
	s.Envelope.DecayRate = v & 0xF
	s.Envelope.DecayEnabled = (v & 0x10) == 0

	if s.Envelope.DecayEnabled {
		s.Envelope.Volume = s.Envelope.DecayRate
	} else {
		s.Envelope.Volume = v & 0xF
	}
}

func (s *Square) WriteSweeps(v Word) {
	s.SweepEnabled = v&0x80 == 0x80
	s.Sweep = ((v >> 4) & 0x7)
	s.SweepMode = (v >> 3) & 1
	s.Negative = v&0x10 == 0x10
	s.Shift = v & 0x7
}

func (s *Square) WriteLow(v Word) {
	s.Timer = (s.Timer & 0x700) | int(v)
}

func (s *Square) WriteHigh(v Word) {
	s.Timer = (s.Timer & 0xFF) | (int(v&0x7) << 8)

	if s.Enabled {
		s.Length = LengthTable[v>>3]
	}

	s.DutyCount = 0
	s.Envelope.Reset = true
}

func (s *Square) UpdateSample(v int16) {
	s.Sample = v
}

func (s *Square) Clock() {
	if s.Length > 0 && s.Timer > 7 {
		if s.TimerCount == 0 {
			s.DutyCount = (s.DutyCount + 1) & 0x7

			s.TimerCount = (s.Timer + 1) << 1
		}

		if !s.Negative && (s.Timer+(s.Timer>>s.Shift)) > 0xFFF {
			s.UpdateSample(0)
		} else if s.Timer < 8 {
			s.UpdateSample(0)
		} else if SquareLookup[s.DutyCycle][s.DutyCount] == 1 {
			s.UpdateSample(int16(s.Volume))
		} else {
			s.UpdateSample(0)
		}

		s.TimerCount--
	} else {
		s.UpdateSample(0)
	}
}

func (s *Square) ClockSweep() {
	if s.SweepEnabled && s.Sweep > 0 {
		s.Sweep--

		if s.Negative {
			s.Timer = s.Timer - (s.Timer >> s.Shift)
		} else {
			s.Timer = s.Timer + (s.Timer >> s.Shift)
		}
	}
}

func (e *Envelope) ClockDecay() {
	if e.Reset {
		e.Reset = false
		e.DecayCounter = e.DecayRate + 1
		e.Volume = 0xF
	} else if e.DecayCounter-1 <= 0 {
		e.DecayCounter = e.DecayRate + 1
		if e.Volume > 0 {
			e.Volume--
		}
	} else {
		e.DecayCounter--
	}
}

func (t *Triangle) Clock() {
	if t.Length > 0 && t.Counter > 0 {
		if t.TimerCount == 0 {
			t.LookupCounter = (t.LookupCounter + 1) % 32

			t.TimerCount = (t.Timer + 1) * 2
		}

		t.UpdateSample(TriangleLookup[t.LookupCounter])

		t.TimerCount--
	} else {
		t.UpdateSample(0)
	}
}

func (t *Triangle) ClockLinearCounter() {
	if t.Halt {
		t.Counter = int(t.ReloadValue)
	} else if t.Counter > 0 {
		t.Counter--
	}

	if !t.Control {
		t.Halt = false
	}
}

func (t *Triangle) UpdateSample(v int16) {
	t.Sample = v
}

func (n *Noise) Clock() {
	var feedback, tmp int

	if n.Length == 0 {
		n.UpdateSample(0)
		return
	}

	if n.Shift&0x1 == 0x0 {
		n.UpdateSample(float64(n.Envelope.Volume))
	} else {
		n.UpdateSample(0)
	}

	n.TimerCount--

	if n.TimerCount == 0 {
		if n.Mode {
			tmp = n.Shift & 0x40 >> 6
		} else {
			tmp = n.Shift & 0x2 >> 1
		}

		feedback = n.Shift&0x1 ^ tmp

		n.Shift = (n.Shift >> 1) | (feedback << 14)

		n.TimerCount = n.Timer
	}
}

func (n *Noise) UpdateSample(v float64) {
	n.Sample = v * 0.2
}

func (a *Apu) Init() <-chan int16 {
	al := make(chan int16, 100)
	a.Output = al

	a.Noise.Shift = 1

	return al
}

func (a *Apu) Step() {
	// Square1
	if a.Square1.Enabled {
		a.Square1.Clock()
	}

	// Square2
	if a.Square2.Enabled {
		a.Square2.Clock()
	}

	// Triangle
	if a.Triangle.Enabled {
		a.Triangle.Clock()
	}

	// Noise
	if a.Noise.Enabled {
		a.Noise.Clock()
	}

	a.Sample = a.ComputeSample()
	a.Sample = a.RunHipassStrong(a.Sample)
	a.Sample = a.RunHipassWeak(a.Sample)
}

func (a *Apu) RunHipassStrong(s int16) int16 {
	a.HipassStrong += (((int64(s) << 16) - (a.HipassStrong >> 16)) * HiPassStrong) >> 16
	return int16(int64(s) - (a.HipassStrong >> 32))
}

func (a *Apu) RunHipassWeak(s int16) int16 {
	a.HipassWeak += (((int64(s) << 16) - (a.HipassWeak >> 16)) * HiPassWeak) >> 16
	return int16(int64(s) - (a.HipassWeak >> 32))
}

func (a *Apu) ComputeSample() int16 {
	pulse_out := 0.00752 * float64(a.Square1.Sample+a.Square2.Sample)
	tnd_out := 0.00851*float64(a.Triangle.Sample) + 0.00494*float64(a.Noise.Sample)

	return int16((pulse_out + tnd_out) * 40000)
}

func (a *Apu) PushSample() {
	a.Output <- a.Sample
}

func (a *Apu) FrameSequencerStep() {
	a.FrameTick++

	if a.FrameTick == 1 || a.FrameTick == a.FrameCounter-1 {
		if a.Square1.LengthEnabled && a.Square1.Length > 0 {
			a.Square1.Length--
		}

		if a.Square2.LengthEnabled && a.Square2.Length > 0 {
			a.Square2.Length--
		}

		if a.Triangle.LengthEnabled && a.Triangle.Length > 0 {
			a.Triangle.Length--
		}

		if a.Noise.LengthEnabled && a.Noise.Length > 0 {
			a.Noise.Length--
		}

		a.Square1.ClockSweep()
		a.Square2.ClockSweep()
	}

	a.Square1.Envelope.ClockDecay()
	a.Square2.Envelope.ClockDecay()
	a.Noise.Envelope.ClockDecay()
	a.Triangle.ClockLinearCounter()

	if a.FrameTick >= a.FrameCounter {
		a.FrameTick = 0
	}
}

func (a *Apu) RegRead(addr int) (Word, error) {
	switch addr {
	case 0x4015:
		// TODO: When a status read occurrs, emulate the APU up to that point.
		// http://forums.nesdev.com/viewtopic.php?t=2123
		return a.ReadStatus(), nil
	}

	return 0, nil
}

func (a *Apu) RegWrite(v Word, addr int) {
	switch addr & 0xFF {
	case 0x0:
		a.WriteSquare1Control(v)
	case 0x1:
		a.WriteSquare1Sweeps(v)
	case 0x2:
		a.WriteSquare1Low(v)
	case 0x3:
		a.WriteSquare1High(v)
	case 0x4:
		a.WriteSquare2Control(v)
	case 0x5:
		a.WriteSquare2Sweeps(v)
	case 0x6:
		a.WriteSquare2Low(v)
	case 0x7:
		a.WriteSquare2High(v)
	case 0x8:
		a.WriteTriangleControl(v)
	case 0xA:
		a.WriteTriangleLow(v)
	case 0xB:
		a.WriteTriangleHigh(v)
	case 0xC:
		a.WriteNoiseBase(v)
	case 0xE:
		a.WriteNoisePeriod(v)
	case 0xF:
		a.WriteNoiseLength(v)
	case 0x15:
		a.WriteControlFlags1(v)
	case 0x17:
		a.WriteControlFlags2(v)
	}
}

// $4015 (w)
func (a *Apu) WriteControlFlags1(v Word) {
	// 76543210
	//    |||||
	//    ||||+- Square 1 (0: disable; 1: enable)
	//    |||+-- Square 2
	//    ||+--- Triangle
	//    |+---- Noise
	//    +----- DMC
	a.Square1.Enabled = (v & 0x1) != 0x0
	a.Square2.Enabled = ((v >> 1) & 0x1) != 0x0
	a.Triangle.Enabled = ((v >> 2) & 0x1) != 0x0
	a.Noise.Enabled = ((v >> 3) & 0x1) != 0x0
	a.DmcEnabled = ((v >> 4) & 0x1) != 0x0

	if !a.Square1.Enabled {
		a.Square1.Length = 0
		a.Square1.UpdateSample(0)
	}

	if !a.Square2.Enabled {
		a.Square2.Length = 0
		a.Square2.UpdateSample(0)
	}

	if !a.Triangle.Enabled {
		a.Triangle.Length = 0
		a.Triangle.UpdateSample(0)
	}

	if !a.Noise.Enabled {
		a.Noise.Length = 0
		a.Noise.UpdateSample(0)
	}

	// TODO:
	// If the DMC bit is clear, the DMC bytes remaining will be 
	// set to 0 and the DMC will silence when it empties.
	// If the DMC bit is set, the DMC sample will be restarted 
	// only if its bytes remaining is 0. Writing to this register 
	// clears the DMC interrupt flag.
}

// $4015 (r)
func (a *Apu) ReadStatus() Word {
	// if-d nt21   DMC IRQ, frame IRQ, length counter statuses
	var status Word

	status |= a.Square1.Length
	status |= a.Square2.Length << 1
	status |= a.Triangle.Length << 2
	status |= a.Noise.Length << 3

	if a.IrqActive && a.IrqEnabled {
		status |= 1 << 6
	} else {
		status |= 0 << 6
	}

	// TODO: Noise -> 0x10

	// Reading this register clears the frame interrupt 
	// flag (but not the DMC interrupt flag).
	// If an interrupt flag was set at the same moment of 
	// the read, it will read back as 1 but it will not be cleared.

	return status & 0xFF
}

// $4017
func (a *Apu) WriteControlFlags2(v Word) {
	// fd-- ----   5-frame cycle, disable frame interrupt
	if v&0x80 == 0x80 {
		a.FrameCounter = 5
		a.FrameSequencerStep()
	} else {
		a.FrameCounter = 4
	}

	a.IrqEnabled = v&0x40 != 0x40
}

// $4000
func (a *Apu) WriteSquare1Control(v Word) {
	a.Square1.WriteControl(v)
}

// $4001
func (a *Apu) WriteSquare1Sweeps(v Word) {
	a.Square1.WriteSweeps(v)
}

// $4002
func (a *Apu) WriteSquare1Low(v Word) {
	a.Square1.WriteLow(v)
}

// $4003
func (a *Apu) WriteSquare1High(v Word) {
	a.Square1.WriteHigh(v)
}

// $4004
func (a *Apu) WriteSquare2Control(v Word) {
	a.Square2.WriteControl(v)
}

// $4005
func (a *Apu) WriteSquare2Sweeps(v Word) {
	a.Square2.WriteSweeps(v)
}

// $4006
func (a *Apu) WriteSquare2Low(v Word) {
	a.Square2.WriteLow(v)
}

// $4007
func (a *Apu) WriteSquare2High(v Word) {
	a.Square2.WriteHigh(v)
}

// $4008
func (a *Apu) WriteTriangleControl(v Word) {
	// 76543210
	// ||||||||
	// |+++++++- ReloadValue
	// +-------- Control Flag (0: use internal counters; 1: disable internal counters)
	a.Triangle.ReloadValue = v & 0x7F
	a.Triangle.Control = (v & 0x80) != 0
	a.Triangle.LengthEnabled = !a.Triangle.Control
}

// $400A
func (a *Apu) WriteTriangleLow(v Word) {
	a.Triangle.Timer = (a.Triangle.Timer & 0x700) | int(v)
}

// $400B
func (a *Apu) WriteTriangleHigh(v Word) {
	a.Triangle.Timer = (a.Triangle.Timer & 0xFF) | (int(v&0xF) << 8)

	if a.Triangle.Enabled {
		a.Triangle.Length = LengthTable[v>>3]
	}

	a.Triangle.Halt = true
}

// $400C
func (a *Apu) WriteNoiseBase(v Word) {
	// --LC NNNN	 Envelope loop / length counter disable (L), constant volume (C), volume/envelope (V)
	a.Noise.LengthEnabled = (v & 0x20) != 0x20
	a.Noise.Envelope.Volume = v & 0x1F
	a.Noise.BaseEnvelope = a.Noise.Envelope.Volume
	a.Noise.Envelope.LoopEnabled = v&0x20 == 0x20
	a.Noise.Envelope.DecayRate = v & 0xF
	a.Noise.Envelope.DecayEnabled = (v & 0x10) == 0

	if a.Noise.Envelope.DecayEnabled {
		a.Noise.Envelope.Volume = a.Noise.Envelope.DecayRate
	} else {
		a.Noise.Envelope.Volume = v & 0xF
	}
}

// $400E
func (a *Apu) WriteNoisePeriod(v Word) {
	// L--- PPPP	 Mode noise (L), noise period (P)
	a.Noise.Mode = v&0x80 == 0x80
	a.Noise.Timer = NoiseLookup[v&0xF]
	a.Noise.TimerCount = a.Noise.Timer
}

// $400F
func (a *Apu) WriteNoiseLength(v Word) {
	// LLLL L---	 Length counter load (L)
	a.Noise.Length = LengthTable[v>>3]
	a.Noise.Envelope.Volume = a.Noise.BaseEnvelope
}
