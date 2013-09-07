package main

import (
	"fmt"
)

var (
	c  *Cpu
	pc uint16
)

func immediateAddress() int {
	val, _ := Ram.Read(pc - 1)
	return int(val)
}

func absoluteAddress() (result int) {
	// Switch to an int (or more appropriately uint16) since we
	// will overflow when shifting the high byte
	high, _ := Ram.Read(pc + 1)
	low, _ := Ram.Read(pc)

	return (int(high) << 8) + int(low)
}

func zeroPageAddress() int {
	res, _ := Ram.Read(pc)

	return int(res)
}

func indirectAbsoluteAddress() (result int) {
	high, _ := Ram.Read(pc + 1)
	low, _ := Ram.Read(pc)

	result = (int(high) << 8) + int(low)
	pc++
	return
}

func absoluteIndexedAddress(index Word) (result int) {
	// Switch to an int (or more appropriately uint16) since we
	// will overflow when shifting the high byte
	high, _ := Ram.Read(pc + 1)
	low, _ := Ram.Read(pc)

	return (int(high) << 8) + int(low) + int(index)
}

func zeroPageIndexedAddress(index Word) int {
	location, _ := Ram.Read(pc)
	return int(location + index)
}

func indexedIndirectAddress() int {
	location, _ := Ram.Read(pc)
	location = location + c.X

	// Switch to an int (or more appropriately uint16) since we
	// will overflow when shifting the high byte
	high, _ := Ram.Read(location + 1)
	low, _ := Ram.Read(location)

	return (int(high) << 8) + int(low)
}

func indirectIndexedAddress() int {
	location, _ := Ram.Read(pc)

	// Switch to an int (or more appropriately uint16) since we
	// will overflow when shifting the high byte
	high, _ := Ram.Read(location + 1)
	low, _ := Ram.Read(location)

	return (int(high) << 8) + int(low) + int(c.Y)
}

func relativeAddress() int {
	return 0
}

func accumulatorAddress() int {
	return 0
}

func Disassemble(opcode Word, cpu *Cpu, p uint16) {
	c = cpu
	pc = p

	fmt.Printf("0x%X: 0x%X ", pc-1, opcode)

	switch opcode {
	// ADC
	case 0x69:
		fmt.Printf("ADC $%X\n", immediateAddress())
	case 0x65:
		fmt.Printf("ADC $%X\n", zeroPageAddress())
	case 0x75:
		fmt.Printf("ADC $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0x6D:
		fmt.Printf("ADC $%X\n", absoluteAddress())
	case 0x7D:
		fmt.Printf("ADC $%X,X\n", absoluteIndexedAddress(c.X))
	case 0x79:
		fmt.Printf("ADC $%X,Y\n", absoluteIndexedAddress(c.Y))
	case 0x61:
		fmt.Printf("ADC ($%X,X)\n", indexedIndirectAddress())
	case 0x71:
		fmt.Printf("ADC ($%X),Y\n", indirectIndexedAddress())
	// LDA
	case 0xA9:
		fmt.Printf("LDA $%X\n", immediateAddress())
	case 0xA5:
		fmt.Printf("LDA $%X\n", zeroPageAddress())
	case 0xB5:
		fmt.Printf("LDA $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0xAD:
		fmt.Printf("LDA $%X\n", absoluteAddress())
	case 0xBD:
		fmt.Printf("LDA $%X,X\n", absoluteIndexedAddress(c.X))
	case 0xB9:
		fmt.Printf("LDA $%X,Y\n", absoluteIndexedAddress(c.Y))
	case 0xA1:
		fmt.Printf("LDA ($%X,X)\n", indexedIndirectAddress())
	case 0xB1:
		fmt.Printf("LDA ($%X),Y\n", indirectIndexedAddress())
	// LDX
	case 0xA2:
		fmt.Printf("LDX $%X\n", immediateAddress())
	case 0xA6:
		fmt.Printf("LDX $%X\n", zeroPageAddress())
	case 0xB6:
		fmt.Printf("LDX $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0xAE:
		fmt.Printf("LDX $%X\n", absoluteAddress())
	case 0xBE:
		fmt.Printf("LDX $%X,Y\n", absoluteIndexedAddress(c.Y))
	// LDY
	case 0xA0:
		fmt.Printf("LDY $%X\n", immediateAddress())
	case 0xA4:
		fmt.Printf("LDY $%X\n", zeroPageAddress())
	case 0xB4:
		fmt.Printf("LDY $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0xAC:
		fmt.Printf("LDY $%X\n", absoluteAddress())
	case 0xBC:
		fmt.Printf("LDY $%X,X\n", absoluteIndexedAddress(c.X))
	// STA
	case 0x85:
		fmt.Printf("STA $%X\n", zeroPageAddress())
	case 0x95:
		fmt.Printf("STA $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0x8D:
		fmt.Printf("STA $%X\n", absoluteAddress())
	case 0x9D:
		fmt.Printf("STA $%X,X\n", absoluteIndexedAddress(c.X))
	case 0x99:
		fmt.Printf("STA $%X,Y\n", absoluteIndexedAddress(c.Y))
	case 0x81:
		fmt.Printf("STA ($%X,X)\n", indexedIndirectAddress())
	case 0x91:
		fmt.Printf("STA ($%X),Y\n", indirectIndexedAddress())
	// STX
	case 0x86:
		fmt.Printf("STX $%X\n", zeroPageAddress())
	case 0x96:
		fmt.Printf("STX $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0x8E:
		fmt.Printf("STX $%X\n", absoluteAddress())
	// STY
	case 0x84:
		fmt.Printf("STY $%X\n", zeroPageAddress())
	case 0x94:
		fmt.Printf("STY $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0x8C:
		fmt.Printf("STY $%X\n", absoluteAddress())
	// JMP
	case 0x4C:
		fmt.Printf("JMP $%X\n", absoluteAddress())
	case 0x6C:
		fmt.Printf("JMP $%X\n", indirectAbsoluteAddress())
	// JSR
	case 0x20:
		fmt.Printf("JSR $%X\n", absoluteAddress())
	// Register Instructions
	case 0xAA:
		fmt.Println("TAX")
	case 0x8A:
		fmt.Println("TXA")
	case 0xCA:
		fmt.Println("DEX")
	case 0xE8:
		fmt.Println("INX")
	case 0xA8:
		fmt.Println("TAY")
	case 0x98:
		fmt.Println("TYA")
	case 0x88:
		fmt.Println("DEY")
	case 0xC8:
		fmt.Println("INY")
	// Branch Instructions
	case 0x10:
		fmt.Printf("BPL $%X\n", immediateAddress())
	case 0x30:
		fmt.Printf("BMI $%X\n", immediateAddress())
	case 0x50:
		fmt.Printf("BVC $%X\n", immediateAddress())
	case 0x70:
		fmt.Printf("BVS $%X\n", immediateAddress())
	case 0x90:
		fmt.Printf("BCC $%X\n", immediateAddress())
	case 0xB0:
		fmt.Printf("BCS $%X\n", immediateAddress())
	case 0xD0:
		fmt.Printf("BNE $%X\n", immediateAddress())
	case 0xF0:
		fmt.Printf("BEQ $%X\n", immediateAddress())
	// CMP
	case 0xC9:
		fmt.Printf("CMP $%X\n", immediateAddress())
	case 0xC5:
		fmt.Printf("CMP $%X\n", zeroPageAddress())
	case 0xD5:
		fmt.Printf("CMP $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0xCD:
		fmt.Printf("CMP $%X\n", absoluteAddress())
	case 0xDD:
		fmt.Printf("CMP $%X,X\n", absoluteIndexedAddress(c.X))
	case 0xD9:
		fmt.Printf("CMP $%X,Y\n", absoluteIndexedAddress(c.Y))
	case 0xC1:
		fmt.Printf("CMP ($%X,X)\n", indexedIndirectAddress())
	case 0xD1:
		fmt.Printf("CMP ($%X),Y\n", c.indirectIndexedAddress())
	// CPX
	case 0xE0:
		fmt.Printf("CPX $%X\n", immediateAddress())
	case 0xE4:
		fmt.Printf("CPX $%X\n", zeroPageAddress())
	case 0xEC:
		fmt.Printf("CPX $%X\n", absoluteAddress())
	// CPY
	case 0xC0:
		fmt.Printf("CPY $%X\n", immediateAddress())
	case 0xC4:
		fmt.Printf("CPY $%X\n", zeroPageAddress())
	case 0xCC:
		fmt.Printf("CPY $%X\n", absoluteAddress())
	// SBC
	case 0xE9:
		fmt.Printf("SBC $%X\n", immediateAddress())
	case 0xE5:
		fmt.Printf("SBC $%X\n", zeroPageAddress())
	case 0xF5:
		fmt.Printf("SBC $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0xED:
		fmt.Printf("SBC $%X\n", absoluteAddress())
	case 0xFD:
		fmt.Printf("SBC $%X,X\n", absoluteIndexedAddress(c.X))
	case 0xF9:
		fmt.Printf("SBC $%X,Y\n", absoluteIndexedAddress(c.Y))
	case 0xE1:
		fmt.Printf("SBC ($%X,X)\n", indexedIndirectAddress())
	case 0xF1:
		fmt.Printf("SBC ($%X),Y\n", indirectIndexedAddress())
	// Flag Instructions
	case 0x18:
		fmt.Println("CLC")
	case 0x38:
		fmt.Println("SEC")
	case 0x58:
		fmt.Println("CLI")
	case 0x78:
		fmt.Println("SEI")
	case 0xB8:
		fmt.Println("CLV")
	case 0xD8:
		fmt.Println("CLD")
	case 0xF8:
		fmt.Println("SED")
	// Stack instructions
	case 0x9A:
		fmt.Println("TXS")
	case 0xBA:
		fmt.Println("TSX")
	case 0x48:
		fmt.Println("PHA")
	case 0x68:
		fmt.Println("PLA")
	case 0x08:
		fmt.Println("PHP")
	case 0x28:
		fmt.Println("PLP")
	// AND
	case 0x29:
		fmt.Printf("AND $%X\n", immediateAddress())
	case 0x25:
		fmt.Printf("AND $%X\n", zeroPageAddress())
	case 0x35:
		fmt.Printf("AND $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0x2d:
		fmt.Printf("AND $%X\n", absoluteAddress())
	case 0x3d:
		fmt.Printf("AND $%X,X\n", absoluteIndexedAddress(c.X))
	case 0x39:
		fmt.Printf("AND $%X,Y\n", absoluteIndexedAddress(c.Y))
	case 0x21:
		fmt.Printf("AND ($%X,X)\n", indexedIndirectAddress())
	case 0x31:
		fmt.Printf("AND ($%X),Y\n", indirectIndexedAddress())
	// ORA
	case 0x09:
		fmt.Printf("ORA $%X\n", immediateAddress())
	case 0x05:
		fmt.Printf("ORA $%X\n", zeroPageAddress())
	case 0x15:
		fmt.Printf("ORA $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0x0d:
		fmt.Printf("ORA $%X\n", absoluteAddress())
	case 0x1d:
		fmt.Printf("ORA $%X,X\n", absoluteIndexedAddress(c.X))
	case 0x19:
		fmt.Printf("ORA $%X,Y\n", absoluteIndexedAddress(c.Y))
	case 0x01:
		fmt.Printf("ORA ($%X,X)\n", indexedIndirectAddress())
	case 0x11:
		fmt.Printf("ORA ($%X),Y\n", indirectIndexedAddress())
	// EOR
	case 0x49:
		fmt.Printf("EOR $%X\n", immediateAddress())
	case 0x45:
		fmt.Printf("EOR $%X\n", zeroPageAddress())
	case 0x55:
		fmt.Printf("EOR $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0x4d:
		fmt.Printf("EOR $%X\n", absoluteAddress())
	case 0x5d:
		fmt.Printf("EOR $%X,X\n", absoluteIndexedAddress(c.X))
	case 0x59:
		fmt.Printf("EOR $%X,Y\n", absoluteIndexedAddress(c.Y))
	case 0x41:
		fmt.Printf("EOR ($%X,X)\n", indexedIndirectAddress())
	case 0x51:
		fmt.Printf("EOR ($%X),Y\n", indirectIndexedAddress())
	// DEC
	case 0xc6:
		fmt.Printf("DEC $%X\n", zeroPageAddress())
	case 0xd6:
		fmt.Printf("DEC $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0xce:
		fmt.Printf("DEC $%X\n", absoluteAddress())
	case 0xde:
		fmt.Printf("DEC $%X,X\n", absoluteIndexedAddress(c.X))
	// INC
	case 0xe6:
		fmt.Printf("INC $%X\n", zeroPageAddress())
	case 0xf6:
		fmt.Printf("INC $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0xee:
		fmt.Printf("INC $%X\n", absoluteAddress())
	case 0xfe:
		fmt.Printf("INC $%X,X\n", absoluteIndexedAddress(c.X))
	// BRK
	case 0x00:
		fmt.Println("BRK")
	// RTI
	case 0x40:
		fmt.Println("RTI")
	// RTS
	case 0x60:
		fmt.Println("RTS")
	// NOP
	case 0xea:
		fmt.Println("NOP")
	// LSR
	case 0x4a:
		fmt.Println("LSR A")
	case 0x46:
		fmt.Printf("LSR $%X\n", zeroPageAddress())
	case 0x56:
		fmt.Printf("LSR $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0x4e:
		fmt.Printf("LSR $%X\n", absoluteAddress())
	case 0x5e:
		fmt.Printf("LSR $%X,X\n", absoluteIndexedAddress(c.X))
	// ASL
	case 0x0a:
		fmt.Println("ASL A")
	case 0x06:
		fmt.Printf("ASL $%X\n", zeroPageAddress())
	case 0x16:
		fmt.Printf("ASL $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0x0e:
		fmt.Printf("ASL $%X\n", absoluteAddress())
	case 0x1e:
		fmt.Printf("ASL $%X,X\n", absoluteIndexedAddress(c.X))
	// ROL
	case 0x2a:
		fmt.Println("ROL A")
	case 0x26:
		fmt.Printf("ROL $%X\n", zeroPageAddress())
	case 0x36:
		fmt.Printf("ROL $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0x2e:
		fmt.Printf("ROL $%X\n", absoluteAddress())
	case 0x3e:
		fmt.Printf("ROL $%X,X\n", absoluteIndexedAddress(c.X))
	// ROR
	case 0x6a:
		fmt.Println("ROR A")
	case 0x66:
		fmt.Printf("ROR $%X\n", zeroPageAddress())
	case 0x76:
		fmt.Printf("ROR $%X,X\n", zeroPageIndexedAddress(c.X))
	case 0x6e:
		fmt.Printf("ROR $%X\n", absoluteAddress())
	case 0x7e:
		fmt.Printf("ROR $%X,X\n", absoluteIndexedAddress(c.X))
	// BIT
	case 0x24:
		fmt.Printf("BIT $%X\n", zeroPageAddress())
	case 0x2c:
		fmt.Printf("BIT $%X\n", absoluteAddress())
	}
}
