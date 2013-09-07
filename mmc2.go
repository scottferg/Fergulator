package main

const (
	RegisterPrgBankSelect = iota
	RegisterChrBank1Select
	RegisterChrBank2Select
	RegisterChrBank3Select
	RegisterChrBank4Select
	RegisterMirroringSelect
)

type Mmc2 struct {
	RomBanks  [][]Word
	VromBanks [][]Word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte

	LatchLow  int
	LatchHigh int
	LatchFE0  int
	LatchFE1  int
	LatchFD0  int
	LatchFD1  int
}

func NewMmc2(r *Rom) *Mmc2 {
	m := &Mmc2{
		RomBanks:     r.RomBanks,
		VromBanks:    r.VromBanks,
		PrgBankCount: r.PrgBankCount,
		ChrRomCount:  r.ChrRomCount,
		Battery:      r.Battery,
		Data:         r.Data,
	}

	m.LatchLow = 0xFE
	m.LatchHigh = 0xFE
	m.LatchFD0 = 0
	m.LatchFE0 = 4
	m.LatchFD1 = 0
	m.LatchFE1 = 0

	m.LoadRom()

	return m
}

func (m *Mmc2) LoadRom() {
	// The PRG banks are 8192 bytes in size, half the size of an
	// iNES PRG bank. If your emulator or copier handles PRG data
	// in 16384 byte chunks, you can think of the lower bit as
	// selecting the first or second half of the bank
	//
	// http://forums.nesdev.com/viewtopic.php?p=38182#p38182

	// Write swappable PRG banks (0x8000 and 0xA000)
	m.Write8kRamBank(0, 0x8000)

	// Write hardwired PRG banks (0xC000 and 0xE000)
	// Third to last bank
	m.Write8kRamBank(((len(m.RomBanks)-2)*2)+1, 0xA000)
	// Second to last bank
	m.Write8kRamBank((len(m.RomBanks)-1)*2, 0xC000)
	// Last bank
	m.Write8kRamBank(((len(m.RomBanks)-1)*2)+1, 0xE000)

	WriteVramBank(m.VromBanks, 4, 0x0000, Size4k)
	WriteVramBank(m.VromBanks, 0, 0x1000, Size4k)
}

func (m *Mmc2) Write(v Word, a int) {
	switch m.RegisterNumber(a) {
	case RegisterPrgBankSelect:
		m.PrgBankSelect(v)
	case RegisterChrBank1Select:
		m.ChrBankSelect(v, 1)
	case RegisterChrBank2Select:
		m.ChrBankSelect(v, 2)
	case RegisterChrBank3Select:
		m.ChrBankSelect(v, 3)
	case RegisterChrBank4Select:
		m.ChrBankSelect(v, 4)
	case RegisterMirroringSelect:
		m.MirroringSelect(v)
	}
}

func (m *Mmc2) LatchTrigger(a int) {
	a &= 0x3FF0

	switch {
	case a == 0x0FD0 && m.LatchLow != 0xFD:
		m.LatchLow = 0xFD
		WriteVramBank(m.VromBanks, m.LatchFD0%len(m.VromBanks), 0x0000, Size4k)
	case a == 0x0FE0 && m.LatchLow != 0xFE:
		m.LatchLow = 0xFE
		WriteVramBank(m.VromBanks, m.LatchFE0%len(m.VromBanks), 0x0000, Size4k)
	case a == 0x1FD0 && m.LatchHigh != 0xFD:
		m.LatchHigh = 0xFD
		WriteVramBank(m.VromBanks, m.LatchFD1%len(m.VromBanks), 0x1000, Size4k)
	case a == 0x1FE0 && m.LatchHigh != 0xFE:
		m.LatchHigh = 0xFE
		WriteVramBank(m.VromBanks, m.LatchFE1%len(m.VromBanks), 0x1000, Size4k)
	}
}

func (m *Mmc2) BatteryBacked() bool {
	return m.Battery
}

func (m *Mmc2) RegisterNumber(a int) int {
	switch a >> 12 {
	case 0xA:
		return RegisterPrgBankSelect
	case 0xB:
		return RegisterChrBank1Select
	case 0xC:
		return RegisterChrBank2Select
	case 0xD:
		return RegisterChrBank3Select
	case 0xE:
		return RegisterChrBank4Select
	case 0xF:
		return RegisterMirroringSelect
	}

	return -1
}

func (m *Mmc2) PrgBankSelect(v Word) {
	bank := int(v & 0xF)
	m.Write8kRamBank(bank, 0x8000)
}

func (m *Mmc2) ChrBankSelect(v Word, b int) {
	v &= 0x1F

	switch b {
	case 1:
		m.LatchFD0 = int(v)

		if m.LatchLow == 0xFD {
			WriteVramBank(m.VromBanks, m.LatchFD0%len(m.VromBanks), 0x0000, Size4k)
		}
	case 2:
		m.LatchFE0 = int(v)

		if m.LatchLow == 0xFE {
			WriteVramBank(m.VromBanks, m.LatchFE0%len(m.VromBanks), 0x0000, Size4k)
		}
	case 3:
		m.LatchFD1 = int(v)

		if m.LatchHigh == 0xFD {
			WriteVramBank(m.VromBanks, m.LatchFD1%len(m.VromBanks), 0x1000, Size4k)
		}
	case 4:
		m.LatchFE1 = int(v)

		if m.LatchHigh == 0xFE {
			WriteVramBank(m.VromBanks, m.LatchFE1%len(m.VromBanks), 0x1000, Size4k)
		}
	}
}

func (m *Mmc2) MirroringSelect(v Word) {
	if v&0x1 == 0x1 {
		ppu.Nametables.SetMirroring(MirroringHorizontal)
	} else {
		ppu.Nametables.SetMirroring(MirroringVertical)
	}
}

func (m *Mmc2) Write8kRamBank(bank, dest int) {
	b := (bank >> 1) % len(m.RomBanks)
	offset := (bank % 2) * 0x2000

	WriteOffsetRamBank(m.RomBanks, b, dest, Size8k, offset)
}
