package nes

type Cnrom struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte
}

func (m *Cnrom) Write(v Word, a int) {
	bank := int(v&0x3) * 2
	WriteVramBank(m.VromBanks, bank, 0x0000, Size4k)
	WriteVramBank(m.VromBanks, bank+1, 0x1000, Size4k)
}

func (m *Cnrom) Read(a int) Word {
	return 0
}

func (m *Cnrom) WriteVram(v Word, a int) {
	// Nothing to do
}

func (m *Cnrom) ReadVram(a int) Word {
	return 0
}

func (m *Cnrom) ReadTile(a int) []Word {
	return []Word{}
}

func (m *Cnrom) BatteryBacked() bool {
	return m.Battery
}
