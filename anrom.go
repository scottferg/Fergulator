package main

type Anrom Rom

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

func (m *Anrom) Hook() {
	// No hooks
}

func (m *Anrom) LatchTrigger(a int) {}

func (m *Anrom) BatteryBacked() bool {
	return m.Battery
}
