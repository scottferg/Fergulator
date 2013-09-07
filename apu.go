package main

import (
	"fmt"
)

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
		15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
	}

	NoiseLookup = []int{
		4, 8, 16, 32, 64, 96, 128, 160, 202,
		254, 380, 508, 762, 1016, 2034, 4068,
	}

	DmcFrequency = []int{
		428, 380, 340, 320, 286,
		254, 226, 214, 190, 160,
		142, 128, 106, 84, 72, 54,
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
	Counter      Word
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
	SweepCounter  Word
	SweepMode     Word
	Shift         Word
	SweepReload   bool
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
	Sample        int16
	Envelope
}

type Dmc struct {
	Enabled        bool
	IrqEnabled     bool
	LoopEnabled    bool
	RateIndex      int
	RateCounter    int
	DirectLoad     int
	DirectCounter  int
	Data           Word
	Sample         int16
	SampleAddress  int
	CurrentAddress int
	SampleLength   int
	SampleCounter  int
	ShiftCounter   int
	Frequency      int
	HasSample      bool
}

type Apu struct {
	Square1 Square
	Square2 Square
	Triangle
	Noise
	Dmc
	IrqEnabled   bool
	IrqActive    bool
	HipassStrong int64
	HipassWeak   int64

	FrameCounter  int
	FrameTick     int
	LastFrameTick int

	PulseOut []float64
	TndOut   []float64

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

	s.SweepReload = true
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

		if s.Timer < 8 {
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
	if s.SweepEnabled && s.SweepCounter > 0 {
		s.SweepCounter--

		delta := (s.Timer >> s.Shift)

		if s.SweepCounter == 0 {
			if s.Negative {
				s.Timer = s.Timer - delta
			} else if s.Timer+delta < 0x800 {
				s.Timer = s.Timer + delta
			}
		}
	}

	if s.SweepReload {
		s.SweepReload = false
		s.SweepCounter = s.Sweep
	}
}

func (e *Envelope) ClockDecay() {
	// When the divider outputs a clock, one of two actions occurs:
	// If the counter is non-zero, it is decremented, otherwise if
	// the loop flag is set, the counter is loaded with 15.

	// The envelope unit's volume output depends on the constant volume
	// flag: if set, the envelope parameter directly sets the volume,
	// otherwise the counter's value is the current volume.

	e.DecayCounter--

	if e.Reset {
		e.Reset = false
		e.DecayCounter = e.DecayRate + 1
		e.Counter = 0xF
	} else if e.DecayCounter == 0 {
		e.DecayCounter = e.DecayRate + 1

		if e.Counter > 0 {
			e.Counter--
		} else if e.LoopEnabled {
			e.Counter = 0xF
		}
	}

	if e.DecayEnabled {
		e.Volume = e.Counter
	} else {
		e.Volume = e.DecayRate
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
		n.UpdateSample(int16(n.Envelope.Volume))
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

func (n *Noise) UpdateSample(v int16) {
	n.Sample = v
}

func (d *Dmc) Clock() {
	if d.HasSample {
		if d.Data&1 == 0 {
			if d.DirectCounter > 0 {
				d.DirectCounter--
			}
		} else {
			if d.DirectCounter < 0x3F {
				d.DirectCounter++
			}
		}

		if d.Enabled {
			d.Sample = (int16((d.DirectCounter)<<1) + int16(d.DirectLoad&0x1))
		} else {
			d.Sample = 0
		}

		d.Data = d.Data >> 1
	}

	d.RateCounter--
	if d.RateCounter <= 0 {
		d.HasSample = false

		if d.SampleCounter == 0 && d.LoopEnabled {
			d.CurrentAddress = d.SampleAddress
			d.SampleCounter = d.SampleLength
		}

		if d.SampleCounter > 0 {
			d.FillSample()

			// TODO: Generate IRQ
		}

		d.RateCounter = 8
	}
}

func (d *Dmc) FillSample() {
	cpu.CyclesToWait += 4

	val, _ := Ram.Read(d.CurrentAddress)
	d.Data = val

	d.SampleCounter--
	if d.SampleCounter == 0 && d.LoopEnabled {
		d.SampleCounter = d.SampleLength
	}

	d.CurrentAddress++
	if d.CurrentAddress >= 0xFFFF {
		d.CurrentAddress = 0x8000
	}

	d.HasSample = true
}

func (a *Apu) Init() <-chan int16 {
	al := make(chan int16, 100)
	a.Output = al

	a.Noise.Shift = 1

	a.PulseOut = make([]float64, 31)
	for i := 0; i < len(a.PulseOut); i++ {
		a.PulseOut[i] = 95.52 / (8128.0/float64(i) + 100.0)
	}

	a.TndOut = make([]float64, 203)
	for i := 0; i < len(a.TndOut); i++ {
		a.TndOut[i] = 163.67 / (24329.0/float64(i) + 100.0)
	}

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

	if a.Dmc.Enabled {
		a.Dmc.ShiftCounter--
		if a.Dmc.ShiftCounter <= 0 && a.Dmc.Frequency > 0 {
			a.Dmc.ShiftCounter += a.Dmc.Frequency
			// a.Dmc.Clock()
		}
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
	if false && a.Dmc.Sample > 100 {
		fmt.Printf("DMC: %d\n", a.Dmc.Sample)
	}
	pulse := a.PulseOut[a.Square1.Sample+a.Square2.Sample]
	// tnd := a.TndOut[(3*a.Triangle.Sample)+(2*a.Noise.Sample)+a.Dmc.Sample]
	tnd := a.TndOut[(3*a.Triangle.Sample)+(2*a.Noise.Sample)]

	return int16((pulse + tnd) * 40000)
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
	case 0x10:
		a.WriteDmcFlags(v)
	case 0x11:
		a.WriteDmcDirectLoad(v)
	case 0x12:
		a.WriteDmcSampleAddress(v)
	case 0x13:
		a.WriteDmcSampleLength(v)
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
	a.Dmc.Enabled = ((v >> 4) & 0x1) != 0x0

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
		// a.Triangle.UpdateSample(0)
	}

	if !a.Noise.Enabled {
		a.Noise.Length = 0
		a.Noise.UpdateSample(0)
	}

	// If the DMC bit is clear, the DMC bytes remaining will be
	// set to 0 and the DMC will silence when it empties.
	// If the DMC bit is set, the DMC sample will be restarted
	// only if its bytes remaining is 0. Writing to this register
	// clears the DMC interrupt flag.
	if v>>4&0x1 == 0x0 {
		a.Dmc.SampleCounter = 0
	} else {
		a.Dmc.CurrentAddress = a.Dmc.SampleAddress
		a.Dmc.SampleCounter = a.Dmc.SampleLength
	}
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
	a.Noise.Envelope.Counter = v & 0x1F
	a.Noise.Envelope.LoopEnabled = v&0x20 == 0x20
	a.Noise.Envelope.DecayRate = v & 0xF
	a.Noise.BaseEnvelope = a.Noise.Envelope.DecayRate
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
}

// $4010
func (a *Apu) WriteDmcFlags(v Word) {
	a.Dmc.IrqEnabled = v&0x80 == 0x80
	a.Dmc.LoopEnabled = v&0x40 == 0x40
	a.Dmc.RateIndex = int(v & 0xF)
	a.Dmc.Frequency = DmcFrequency[v&0xF]
}

// $4011
func (a *Apu) WriteDmcDirectLoad(v Word) {
	a.Dmc.DirectLoad = int(v & 0x7F)
	a.Dmc.DirectCounter = a.Dmc.DirectLoad

	a.Dmc.Sample = (int16(a.Dmc.DirectCounter) << 1) + int16(v)&0x1
}

// $4012
func (a *Apu) WriteDmcSampleAddress(v Word) {
	a.Dmc.SampleAddress = int(v) << 6
}

// $4013
func (a *Apu) WriteDmcSampleLength(v Word) {
	a.Dmc.SampleLength = (int(v) << 4) + 1
	a.Dmc.SampleCounter = a.Dmc.SampleLength
}
