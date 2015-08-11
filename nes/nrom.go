package nes

// Nrom
type Nrom struct {
	RomBanks  [][]word
	VromBanks [][]word

	PrgBankCount int
	ChrRomCount  int
	Battery      bool
	Data         []byte
}

func (m *Nrom) Load() {
	m.RomBanks = make([][]word, m.PrgBankCount)
	for i := 0; i < m.PrgBankCount; i++ {
		// Move 16kb chunk to 16kb bank
		bank := make([]word, 0x4000)
		for x := 0; x < 0x4000; x++ {
			bank[x] = word(m.Data[(0x4000*i)+x])
		}

		m.RomBanks[i] = bank
	}

	// Everything after PRG-ROM
	chrRom := m.Data[0x4000*len(m.RomBanks):]

	if m.ChrRomCount > 0 {
		m.VromBanks = make([][]word, m.ChrRomCount*2)
	} else {
		m.VromBanks = make([][]word, 2)
	}

	for i := 0; i < cap(m.VromBanks); i++ {
		// Move 16kb chunk to 16kb bank
		m.VromBanks[i] = make([]word, 0x1000, 0x1000)

		// If the game doesn't have CHR banks we
		// just need to allocate VRAM

		for x := 0; x < 0x1000; x++ {
			var val word
			if m.ChrRomCount == 0 {
				val = 0
			} else {
				val = word(chrRom[(0x1000*i)+x])
			}
			m.VromBanks[i][x] = val
		}
	}
}

func (m *Nrom) Write(v word, a int) {
	// Nothing to do
}

func (m *Nrom) Read(a int) word {
	if a >= 0xC000 {
		return m.RomBanks[len(m.RomBanks)-1][a&0x3FFF]
	}

	return m.RomBanks[0][a&0x3FFF]
}

func (m *Nrom) WriteVram(v word, a int) {
	if a >= 0x1000 {
		m.VromBanks[len(m.VromBanks)-1][a&0xFFF] = v
		return
	}

	m.VromBanks[0][a&0xFFF] = v
}

func (m *Nrom) ReadVram(a int) word {
	if a >= 0x1000 {
		return m.VromBanks[len(m.VromBanks)-1][a&0xFFF]
	}

	return m.VromBanks[0][a&0xFFF]
}

func (m *Nrom) ReadTile(a int) []word {
	if a >= 0x1000 {
		return m.VromBanks[len(m.VromBanks)-1][a&0xFFF : a+16]
	}

	return m.VromBanks[0][a&0xFFF : a+16]
}

func (m *Nrom) BatteryBacked() bool {
	return m.Battery
}
