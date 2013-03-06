package main

const (
	ConfigPrgMode = iota
	ConfigChrMode
	ConfigExtendedRamMode
	ConfigMirroring
	ConfigFillModeTile
	ConfigFillModeColor
)

type Mmc5 struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte
}

func (m *Mmc5) Write(v Word, a int) {
	switch a {
	case ConfigPrgMode:
	case ConfigChrMode:
	case ConfigExtendedRamMode:
	case ConfigMirroring:
	case ConfigFillModeTile:
	case ConfigFillModeColor:
	}
}

func (m *Mmc5) BatteryBacked() bool {
	return m.Battery
}
