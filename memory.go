package main

type Word uint8

type Memory [0x10000]Word

type MemoryError struct {
	ErrorText string
}

func (e MemoryError) Error() string {
	return e.ErrorText
}

var (
	Ram Memory
)

func fitAddressSize(addr interface{}) (v int, e error) {
	switch a := addr.(type) {
	case Word:
		v = int(a)
	case int:
		v = int(a)
	default:
		e = MemoryError{ErrorText: "Invalid type used"}
	}

	return
}

func (m *Memory) Init() {
	for index, _ := range m {
		m[index] = 0x00
	}
}

func (m *Memory) WriteMirroredRam(v Word, a int) {
	for i := 0; i < 0x2FFB; i += 8 {
		m[0x2002+i] = v
	}
}

func (m *Memory) Write(address interface{}, val Word) error {
	if a, err := fitAddressSize(address); err == nil {
		if a == 0x2002 {
			return nil
		}

		m[a] = val

		if a <= 0x2007 && a >= 0x2000 {
            //ppu.Run(cpu.Timestamp * 3)
			ppu.PpuRegWrite(val, a)
		} else if a == 0x4014 {
            //ppu.Run(cpu.Timestamp * 3)
			ppu.PpuRegWrite(val, a)
		} else if a == 0x4016 {
			controller.Write(val)
		} else if a >= 0x8000 && a <= 0xFFFF {
			// MMC1
			rom.Write(val, a)
			return nil
		}

		return nil
	}

	return MemoryError{ErrorText: "Invalid address used"}
}

func (m *Memory) Read(address interface{}) (Word, error) {
	a, _ := fitAddressSize(address)

    if a == 0x200A {
        return m[0x2002], nil
    }

	if a <= 0x2007 && a >= 0x2000 {
        //ppu.Run(cpu.Timestamp)
		return ppu.PpuRegRead(a)
	} else if a == 0x4016 {
		return controller.Read(), nil
	}

	return m[a], nil
}
