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

func (m *Nrom) Load() {
	m.RomBanks = make([][]Word, m.PrgBankCount)
	for i := 0; i < m.PrgBankCount; i++ {
		// Move 16kb chunk to 16kb bank
		bank := make([]Word, 0x4000)
		for x := 0; x < 0x4000; x++ {
			bank[x] = Word(m.Data[(0x4000*i)+x])
		}

		m.RomBanks[i] = bank
	}

	// Everything after PRG-ROM
	chrRom := m.Data[0x4000*len(m.RomBanks):]

	if m.ChrRomCount > 0 {
		m.VromBanks = make([][]Word, m.ChrRomCount*2)
	} else {
		m.VromBanks = make([][]Word, 2)
	}

	for i := 0; i < cap(m.VromBanks); i++ {
		// Move 16kb chunk to 16kb bank
		m.VromBanks[i] = make([]Word, 0x1000, 0x1000)

		// If the game doesn't have CHR banks we
		// just need to allocate VRAM

		for x := 0; x < 0x1000; x++ {
			var val Word
			if m.ChrRomCount == 0 {
				val = 0
			} else {
				val = Word(chrRom[(0x1000*i)+x])
			}
			m.VromBanks[i][x] = val
		}
	}
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
