package main

var (
	SquareLookup = []int{
		0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0,
		1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	}

	TriangleLookup = []int{
		8, 9, 10, 11, 12, 13, 14,
		15, 15, 14, 13, 12, 11, 10,
		9, 8, 7, 6, 5, 4, 3, 2, 1,
		0, 0, 1, 2, 3, 4, 5, 6, 7,
	}

	LengthTable = []Word{
		10, 254, 20, 2, 40, 4, 80, 6,
		160, 8, 60, 10, 14, 12, 26, 14,
		12, 16, 24, 18, 48, 20, 96, 22,
		192, 24, 72, 26, 16, 28, 32, 30,
	}
)

type Square struct {
	Envelope         Word
	EnvelopeDisabled bool
	LengthDisabled   bool
	DutyCycle        Word
	Timer            int
	Length           Word
	DutyLength       int
	Frequency        Word
	Samples          []byte
	SampleIndex      int
}

type Triangle struct {
	Value                    Word
	InternalCountersDisabled bool
	Timer                    int
	Length                   Word
	Frequency                Word
	Counter                  int
	Sample                   int
}

type Apu struct {
	Square1Enabled  bool
	Square2Enabled  bool
	TriangleEnabled bool
	NoiseEnabled    bool
	DmcEnabled      bool
	Square1         Square
	Square2         Square
	Triangle

	Output chan []byte
}

func (s *Square) WriteLow(v Word) {
	s.Timer = (s.Timer & 0x700) | int(v)
}

func (s *Square) WriteHigh(v Word) {
	s.Timer = (s.Timer & 0xFF) | int((v&0xF)<<8)
	s.Length = LengthTable[v>>3]
}

func (s *Square) Clock() {
	if s.Length > 0 {
		if s.EnvelopeDisabled {
			s.Envelope = 0xF
		}

		if s.Length > 0 {
			dutyRow := 8 * s.DutyCycle

			if SquareLookup[int(dutyRow)+s.DutyLength] == 1 {
				s.Samples[s.SampleIndex] = byte(s.Envelope*s.Length+Word(s.Timer))
			} else {
				s.Samples[s.SampleIndex] = byte(0)
			}

			s.DutyLength = (s.DutyLength + 1) & 0xF
		}
	} else {
		s.Samples[s.SampleIndex] = byte(0)
	}

	s.SampleIndex++
	if s.SampleIndex == 0xAC44 {
		s.SampleIndex = 0
	}
}

func (a *Apu) Init() <-chan []byte {
	al := make(chan []byte, 20)

	a.Square1.Samples = make([]byte, 44100)
	a.Square2.Samples = make([]byte, 44100)
	a.Output = al

	return al
}

func (a *Apu) Step() {
	// Square1
	if a.Square1Enabled {
		a.Square1.Clock()
	}

	if a.Square1.SampleIndex == 0xAC43 {
		a.Output <- a.Square1.Samples
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
	a.Square1Enabled = (v & 0x1) == 0x1
	a.Square2Enabled = ((v >> 1) & 0x1) == 0x1
	a.TriangleEnabled = ((v >> 2) & 0x1) == 0x1
	a.NoiseEnabled = ((v >> 3) & 0x1) == 0x1
	a.DmcEnabled = ((v >> 4) & 0x1) == 0x1

	if !a.Square1Enabled {
		a.Square1.Length = 0
	}

	if !a.Square2Enabled {
		a.Square2.Length = 0
	}

	if !a.TriangleEnabled {
		a.Triangle.Length = 0
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

	if a.Square1.Length > 0 {
		status |= 0x1
	}

	if a.Square2.Length > 0 {
		status |= 0x2
	}

	if a.Triangle.Length > 0 {
		status |= 0x8
	}

	// TODO: Noise -> 0x10

	// Reading this register clears the frame interrupt 
	// flag (but not the DMC interrupt flag).
	// If an interrupt flag was set at the same moment of 
	// the read, it will read back as 1 but it will not be cleared.

	return status
}

// $4017
func (a *Apu) WriteControlFlags2(v Word) {
	// fd-- ----   5-frame cycle, disable frame interrupt
	// fmt.Println("WriteControl2!")
}

// $4000
func (a *Apu) WriteSquare1Control(v Word) {
	// 76543210
	// ||||||||
	// ||||++++- Envelope
	// |||+----- Saw Envelope Disable (0: use internal counter for volume; 1: use Envelope for volume)
	// ||+------ Length Counter Disable (0: use Length Counter; 1: disable Length Counter)
	// ++------- Duty Cycle
	a.Square1.EnvelopeDisabled = (v>>4)&0x1 == 1
	a.Square1.LengthDisabled = (v>>5)&0x1 == 1
	a.Square1.DutyCycle = (v >> 6) & 0x3
	a.Square1.Envelope = v & 0xF
	a.Square1.Envelope = a.Square1.Envelope + ((v >> 1) & 0x10)
}

// $4001
func (a *Apu) WriteSquare1Sweeps(v Word) {
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
	// 76543210
	// ||||||||
	// ||||++++- Envelope
	// |||+----- Saw Envelope Disable (0: use internal counter for volume; 1: use Envelope for volume)
	// ||+------ Length Counter Disable (0: use Length Counter; 1: disable Length Counter)
	// ++------- Duty Cycle
	a.Square2.EnvelopeDisabled = (v>>4)&0x1 == 1
	a.Square2.LengthDisabled = (v>>5)&0x1 == 1
	a.Square2.DutyCycle = (v >> 6) & 0x3
	a.Square2.Envelope = v & 0xF
	a.Square2.Envelope = a.Square2.Envelope + ((v >> 1) & 0x10)
}

// $4005
func (a *Apu) WriteSquare2Sweeps(v Word) {
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
	// |+++++++- Value
	// +-------- Control Flag (0: use internal counters; 1: disable internal counters)
	a.Triangle.Value = v & 0x7F
	a.Triangle.InternalCountersDisabled = (v>>7)&0x1 == 1
}

// $400A
func (a *Apu) WriteTriangleLow(v Word) {
	a.Triangle.Timer = (a.Triangle.Timer & 0x700) | int(v)
}

// $400B
func (a *Apu) WriteTriangleHigh(v Word) {
	a.Triangle.Timer = (a.Triangle.Timer & 0xFF) | int((v&0xF)<<8)
	a.Triangle.Length = v >> 3

	// a.Triangle.Frequency = 1789773 / (32*(a.Triangle.Timer) + 1)
}

func (a *Apu) ClockTriangle() {
	if !a.TriangleEnabled || a.Triangle.Timer <= 0 {
		return
	}

	a.Triangle.Sample = TriangleLookup[a.Triangle.Counter]
	if a.Triangle.Sample > 15 {
		a.Triangle.Sample = 31 - a.Triangle.Sample
	}

	a.Triangle.Counter = (a.Triangle.Counter + 1) % 32
}
