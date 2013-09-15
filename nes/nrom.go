package nes

// Nrom
type Nrom struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte
}

func (m *Nrom) Write(v Word, a int) {
	// Nothing to do
}

func (m *Nrom) Read(a int) Word {
	if a >= 0xC000 {
		return m.RomBanks[len(m.RomBanks)-1][a&0x3FFF]
	}

	return m.RomBanks[0][a&0x3FFF]
}

func (m *Nrom) WriteVram(v Word, a int) {
	if a >= 0x1000 {
		m.VromBanks[len(m.VromBanks)-1][a&0xFFF] = v
		return
	}

	m.VromBanks[0][a&0xFFF] = v
}

func (m *Nrom) ReadVram(a int) Word {
	if a >= 0x1000 {
		return m.VromBanks[len(m.VromBanks)-1][a&0xFFF]
	}

	return m.VromBanks[0][a&0xFFF]
}

func (m *Nrom) ReadTile(a int) []Word {
	if a >= 0x1000 {
		return m.VromBanks[len(m.VromBanks)-1][a&0xFFF : a+16]
	}

	return m.VromBanks[0][a&0xFFF : a+16]
}

func (m *Nrom) BatteryBacked() bool {
	return m.Battery
}
