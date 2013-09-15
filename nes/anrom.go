package nes

type Anrom struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte

	PrgUpperBank int
	PrgLowerBank int
}

func (m *Anrom) Write(v Word, a int) {
	if v&0x10 == 0x10 {
		ppu.Nametables.SetMirroring(MirroringSingleUpper)
	} else {
		ppu.Nametables.SetMirroring(MirroringSingleLower)
	}

	bank := int((v & 0x7) * 2)

	m.PrgUpperBank = bank + 1
	m.PrgLowerBank = bank
}

func (m *Anrom) Read(a int) Word {
	if a >= 0xC000 {
		return m.RomBanks[m.PrgUpperBank][a&0x3FFF]
	}

	return m.RomBanks[m.PrgLowerBank][a&0x3FFF]
}

func (m *Anrom) BatteryBacked() bool {
	return m.Battery
}

func (m *Anrom) WriteVram(v Word, a int) {
	if a >= 0x1000 {
		m.VromBanks[len(m.VromBanks)-1][a&0xFFF] = v
	}

	m.VromBanks[0][a&0xFFF] = v
}

func (m *Anrom) ReadVram(a int) Word {
	if a >= 0x1000 {
		return m.VromBanks[len(m.VromBanks)-1][a&0xFFF]
	}

	return m.VromBanks[0][a&0xFFF]
}

func (m *Anrom) ReadTile(a int) []Word {
	if a >= 0x1000 {
		a &= 0xFFF
		return m.VromBanks[len(m.VromBanks)-1][a : a+16]
	}

	a &= 0xFFF
	return m.VromBanks[0][a : a+16]
}
