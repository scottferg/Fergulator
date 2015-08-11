package nes

type Unrom struct {
	RomBanks  [][]word
	VromBanks [][]word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte

	ActiveBank int
}

func (m *Unrom) Write(v word, a int) {
	m.ActiveBank = int(v & 0x7)
}

func (m *Unrom) Read(a int) word {
	if a >= 0xC000 {
		return m.RomBanks[len(m.RomBanks)-1][a&0x3FFF]
	}

	return m.RomBanks[m.ActiveBank][a&0x3FFF]
}

func (m *Unrom) WriteVram(v word, a int) {
	if a >= 0x1000 {
		m.VromBanks[len(m.VromBanks)-1][a&0xFFF] = v
		return
	}

	m.VromBanks[0][a&0xFFF] = v
}

func (m *Unrom) ReadVram(a int) word {
	if a >= 0x1000 {
		return m.VromBanks[len(m.VromBanks)-1][a&0xFFF]
	}

	return m.VromBanks[0][a&0xFFF]
}

func (m *Unrom) ReadTile(a int) []word {
	if a >= 0x1000 {
		a &= 0xFFF
		return m.VromBanks[len(m.VromBanks)-1][a : a+16]
	}

	a &= 0xFFF
	return m.VromBanks[0][a : a+16]
}

func (m *Unrom) BatteryBacked() bool {
	return m.Battery
}
