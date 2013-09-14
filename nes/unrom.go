package nes

type Unrom struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte

	ActiveBank int
}

func (m *Unrom) Write(v Word, a int) {
	m.ActiveBank = int(v & 0x7)
}

func (m *Unrom) Read(a int) Word {
	if a >= 0xC000 {
		return m.RomBanks[len(m.RomBanks)-1][a&0x3FFF]
	}

	return m.RomBanks[m.ActiveBank][a&0x3FFF]
}

func (m *Unrom) WriteVram(v Word, a int) {
	// Nothing to do
}

func (m *Unrom) ReadVram(a int) Word {
	return 0
}

func (m *Unrom) ReadTile(a int) []Word {
	return []Word{}
}

func (m *Unrom) BatteryBacked() bool {
	return m.Battery
}
