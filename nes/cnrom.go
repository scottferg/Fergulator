package nes

type Cnrom struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte

	ActiveBank int
}

func (m *Cnrom) Write(v Word, a int) {
	m.ActiveBank = int(v&0x3) * 2
}

func (m *Cnrom) Read(a int) Word {
	if a >= 0xC000 {
		return m.RomBanks[len(m.RomBanks)-1][a&0x3FFF]
	}

	return m.RomBanks[0][a&0x3FFF]
}

func (m *Cnrom) WriteVram(v Word, a int) {
	if a >= 0x1000 {
		m.VromBanks[m.ActiveBank+1][a&0xFFF] = v
		return
	}

	m.VromBanks[m.ActiveBank][a&0xFFF] = v
}

func (m *Cnrom) ReadVram(a int) Word {
	if a >= 0x1000 {
		return m.VromBanks[m.ActiveBank+1][a&0xFFF]
	}

	return m.VromBanks[m.ActiveBank][a&0xFFF]
}

func (m *Cnrom) ReadTile(a int) []Word {
	if a >= 0x1000 {
		return m.VromBanks[m.ActiveBank+1][a&0xFFF : a+16]
	}

	return m.VromBanks[m.ActiveBank][a&0xFFF : a+16]
}

func (m *Cnrom) BatteryBacked() bool {
	return m.Battery
}
