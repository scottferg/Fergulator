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

func PpuRegWrite(v Word, a int) {
    switch a {
    case 0x2000:
        ppu.WriteControl(v)
    case 0x2001:
        ppu.WriteMask(v)
    case 0x2003:
        ppu.WriteOamAddress(v)
    case 0x2004:
        ppu.WriteOamData(v)
    case 0x2005:
        ppu.WriteScroll(v)
    case 0x2006:
        ppu.WriteAddress(v)
    case 0x2007:
        ppu.WriteData(v)
    case 0x4014:
        ppu.WriteDma(v)
    }
}

func PpuRegRead(a int) (Word, error) {
    switch a {
    case 0x2002:
        return ppu.ReadStatus()
    case 0x2004:
        return ppu.ReadOamData()
    case 0x2007:
        return ppu.ReadData()
    }

    return 0, nil
}

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

func (m *Memory) Write(address interface{}, val Word) error {
	if a, err := fitAddressSize(address); err == nil {
		m[a] = val

        if a <= 0x2007 && a >= 0x2000 {
            PpuRegWrite(val, a)
        } else if a == 0x4014 {
            PpuRegWrite(val, a)
        }

		return nil
	}

	return MemoryError{ErrorText: "Invalid address used"}
}

func (m *Memory) Read(address interface{}) (Word, error) {
	a, _ := fitAddressSize(address)

    if a <= 0x2007 && a >= 0x2000 {
        return PpuRegRead(a)
    }

	return m[a], nil
}
