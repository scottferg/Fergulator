package nes

type Anrom struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte
}

func (m *Anrom) Write(v Word, a int) {
	bank := int((v & 0x7) * 2)

	if v&0x10 == 0x10 {
		ppu.Nametables.SetMirroring(MirroringSingleUpper)
	} else {
		ppu.Nametables.SetMirroring(MirroringSingleLower)
	}

	WriteRamBank(m.RomBanks, bank, 0x8000, Size16k)
	WriteRamBank(m.RomBanks, bank+1, 0xC000, Size16k)
}

func (m *Anrom) Read(a int) Word {
	return 0
}

func (m *Anrom) BatteryBacked() bool {
	return m.Battery
}

func (m *Anrom) WriteVram(v Word, a int) {
	// Nothing to do
}

func (m *Anrom) ReadVram(a int) Word {
	return 0
}

func (m *Anrom) ReadTile(a int) []Word {
	return []Word{}
}
