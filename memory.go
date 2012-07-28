package main

type Word uint8

type Memory [0x10000]*Word

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
        e = MemoryError{ErrorText: "Invalid type used",}
    }

    return
}

func (m *Memory) Init() {
    for index, _ := range m {
        m[index] = new(Word)
    }
}

func (m *Memory) Write(address interface{}, val Word) error {
    if a, err := fitAddressSize(address); err == nil {
        m[a] = &val
        return nil
    }

    return MemoryError{ErrorText: "Invalid address used"}
}

func (m *Memory) WriteMirrorable(address interface{}, val *Word) error {
    if a, err := fitAddressSize(address); err == nil {
        m[a] = val
        return nil
    }

    return MemoryError{ErrorText: "Invalid address used"}
}

func (m *Memory) Read(address interface{}) (Word, error) {
    a, _ := fitAddressSize(address)
    return *m[a], nil
}

func (m *Memory) ReadMirrorable(address interface{}) (*Word, error) {
    a, _ := fitAddressSize(address)
    return m[a], nil
}
