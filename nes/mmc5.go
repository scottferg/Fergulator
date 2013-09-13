package nes

type Mmc5 struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte

	PrgSwitchMode   int
	ChrSwitchMode   int
	ExtendedRamMode int
}

func (m *Mmc5) Write(v Word, a int) {
	/*
		switch a {
		case 0x5100:
			// PRG Switching mode
			m.PrgSwitchMode = v & 0x3
		case 0x5101:
			// CHR Switching mode
			m.ChrSwitchMode = v & 0x3
		case 0x5102:
			// PRG RAM protect 1
			fmt.Println("PRG RAM protect 1")
		case 0x5103:
			// PRG RAM protect 2
			fmt.Println("PRG RAM protect 2")
		case 0x5104:
			// Extended RAM mode
			m.ExtendedRamMode = v & 0x3
		case 0x5105:
			// Nametable mapping
			m.SetNametableMapping(v)
		}
	*/
}

func (m *Mmc5) BatteryBacked() bool {
	return m.Battery
}

func (m *Mmc5) SetNametableMapping(v Word) {
}
