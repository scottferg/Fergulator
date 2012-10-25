package main

type Square struct {
	Volume                Word
	SawEnvelopeDisabled   bool
	LengthCounterDisabled bool
	DutyCycle             Word
	LowPeriod             Word
	HighPeriod            Word
	LengthCounter         Word
}

type Triangle struct {
}

type Apu struct {
	Square1Enabled  bool
	Square2Enabled  bool
	TriangleEnabled bool
	NoiseEnabled    bool
	DmcEnabled      bool
	Square1         Square
	Square2         Square
}

func (a *Apu) NewApu() {
}

// $4015
func (a *Apu) WriteFlags(v Word) {
	// 76543210
	//    |||||
	//    ||||+- Square 1 (0: disable; 1: enable)
	//    |||+-- Square 2
	//    ||+--- Triangle
	//    |+---- Noise
	//    +----- DMC
	a.Square1Enabled = (v & 0x1) == 0x1
	a.Square2Enabled = ((v >> 1) & 0x1) == 0x1
	a.TriangeEnabled = ((v >> 2) & 0x1) == 0x1
	a.NoiseEnabled = ((v >> 3) & 0x1) == 0x1
	a.DmcEnabled = ((v >> 4) & 0x1) == 0x1
}

// $4000
func (a *Apu) WriteSquare1Control(v Word) {
	// 76543210
	// ||||||||
	// ||||++++- Volume
	// |||+----- Saw Envelope Disable (0: use internal counter for volume; 1: use Volume for volume)
	// ||+------ Length Counter Disable (0: use Length Counter; 1: disable Length Counter)
	// ++------- Duty Cycle
	a.Square1.Volume = v & 0xF
	a.Square1.SawEnvelopeDisabled = (v>>4)&0x1 == 1
	a.Square1.LengthCounterDisabled = (v>>5)&0x1 == 1
	a.Square1.DutyCycle = (v >> 6) & 0x3
}

// $4001
func (a *Apu) WriteSquare1Sweeps(v Word) {
}

// $4002
func (a *Apu) WriteSquare1Low(v Word) {
	a.Square1.LowPeriod = v
}

// $4003
func (a *Apu) WriteSquare1High(v Word) {
	a.Square1.HighPeriod = v & 0xF
	a.Square1.LengthCounter = v >> 3
}

// $4004
func (a *Apu) WriteSquare2Control(v Word) {
	// 76543210
	// ||||||||
	// ||||++++- Volume
	// |||+----- Saw Envelope Disable (0: use internal counter for volume; 1: use Volume for volume)
	// ||+------ Length Counter Disable (0: use Length Counter; 1: disable Length Counter)
	// ++------- Duty Cycle
	a.Square2.Volume = v & 0xF
	a.Square2.SawEnvelopeDisabled = (v>>4)&0x1 == 1
	a.Square2.LengthCounterDisabled = (v>>5)&0x1 == 1
	a.Square2.DutyCycle = (v >> 6) & 0x3
}

// $4005
func (a *Apu) WriteSquare2Sweeps(v Word) {
}

// $4006
func (a *Apu) WriteSquare2Low(v Word) {
	a.Square2.LowPeriod = v
}

// $4007
func (a *Apu) WriteSquare2High(v Word) {
	a.Square2.HighPeriod = v & 0xF
	a.Square2.LengthCounter = v >> 3
}
