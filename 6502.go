package main

import (
    "fmt"
)

type Cpu struct {
    X            Word
    Y            Word
    A            Word
    CycleCount   int
    Negative     bool
    Overflow     bool
    BrkCommand   bool
    DecimalMode  bool
    IrqDisable   bool
    Zero         bool
    Carry        bool
    StackPointer Word
    Verbose      bool
}

func (cpu *Cpu) testAndSetNegative(value Word) {
    if value & 0x80 > 0x00 {
        cpu.Negative= true
        return
    }

    cpu.Negative = false
}

func (cpu *Cpu) testAndSetZero(value Word) {
    if value == 0x00 {
        cpu.Zero = true
        return
    }

    cpu.Zero = false
}

func (cpu *Cpu) testAndSetCarryAddition(a Word, b Word) {
    if int(a + b) > 0xff {
        cpu.Carry = true
        return
    }

    cpu.Carry = false
}

func (cpu *Cpu) testAndSetCarrySubtraction(a Word, b Word) {
    if int(a - b) < 0x00 {
        cpu.Carry = false
        return
    }

    cpu.Carry = true
}

func (cpu *Cpu) testAndSetOverflowAddition(a Word, b Word) {
    if (a & 0x80) == (b & 0x80) {
        switch {
        case int(a + b) > 127:
        case int(a + b) < -128:
            cpu.Overflow = true
            return
        }
    }

    cpu.Overflow = false
}

func (cpu *Cpu) testAndSetOverflowSubtraction(a Word, b Word) {
    if (a & 0x80) != (b & 0x80) {
        switch {
        case int(a - b) > 127:
        case int(a - b) < -128:
            cpu.Overflow = true
            return
        }
    }

    cpu.Overflow = false
}

func (cpu *Cpu) immediateAddress() (int) {
    programCounter++
    return programCounter - 1
}

func (cpu *Cpu) absoluteAddress() (result int) {
    // Switch to an int (or more appropriately uint16) since we 
    // will overflow when shifting the high byte
    high := int(Ram[programCounter + 1])
    low := int(Ram[programCounter])

    programCounter += 2
    return (high << 8) + low
}

func (cpu *Cpu) zeroPageAddress() (int) {
    programCounter++
    return int(Ram[programCounter - 1])
}

func (cpu *Cpu) indirectAbsoluteAddress() (result int) {
    result = int((Ram[programCounter + 1] << 8) + Ram[programCounter])
    programCounter++
    return
}

func (cpu *Cpu) absoluteIndexedAddress(index Word) (result int) {
    // Switch to an int (or more appropriately uint16) since we 
    // will overflow when shifting the high byte
    high := int(Ram[programCounter + 1])
    low := int(Ram[programCounter])

    programCounter++
    return (high << 8) + low + int(index)
}

func (cpu *Cpu) zeroPageIndexedAddress(index Word) (int) {
    location := int(Ram[programCounter] + index)
    programCounter++
    return location
}

func (cpu *Cpu) indexedIndirectAddress() (int) {
    location := Ram[programCounter] + cpu.X
    programCounter++

    // Switch to an int (or more appropriately uint16) since we 
    // will overflow when shifting the high byte
    high := int(Ram[location + 1])
    low := int(Ram[location])

    return (high << 8) + low
}

func (cpu *Cpu) indirectIndexedAddress() (int) {
    location := Ram[programCounter]

    // Switch to an int (or more appropriately uint16) since we 
    // will overflow when shifting the high byte
    high := int(Ram[location + 1])
    low := int(Ram[location])

    programCounter++
    return (high << 8) + low + int(cpu.Y)
}

func (cpu *Cpu) relativeAddress() (int) {
    return 0
}

func (cpu *Cpu) accumulatorAddress() (int) {
    return 0
}

func (cpu *Cpu) Adc(location int) {
    cpu.A = cpu.A + Ram[location]

    if cpu.Carry {
        cpu.A++
    }

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
    cpu.testAndSetOverflowAddition(cpu.A, Ram[location])
    cpu.testAndSetCarryAddition(cpu.A, Ram[location])
}

func (cpu *Cpu) Lda(location int) {
    cpu.A = Ram[location]

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Ldx(location int) {
    cpu.X = Ram[location]

    cpu.testAndSetNegative(cpu.X)
    cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu) Ldy(location int) {
    cpu.Y = Ram[location]

    cpu.testAndSetNegative(cpu.Y)
    cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu) Sta(location int) {
    Ram[location] = cpu.A
}

func (cpu *Cpu) Stx(location int) {
    Ram[location] = cpu.X
}

func (cpu *Cpu) Sty(location int) {
    Ram[location] = cpu.Y
}

func (cpu *Cpu) Jmp(location int) {
    programCounter = location
}

func (cpu *Cpu) Tax() {
    cpu.X = cpu.A

    cpu.testAndSetNegative(cpu.X)
    cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu) Txa() {
    cpu.A = cpu.X

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Dex() {
    cpu.X = cpu.X - 1

    if cpu.X == 0 {
        cpu.Zero = true
    }

    cpu.testAndSetNegative(cpu.X)
    cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu) Inx() {
    cpu.X = cpu.X + 1

    cpu.testAndSetNegative(cpu.X)
    cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu) Tay() {
    cpu.Y = cpu.A

    cpu.testAndSetNegative(cpu.Y)
    cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu) Tya() {
    cpu.A = cpu.Y

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Dey() {
    cpu.Y = cpu.Y - 1

    if cpu.X == 0 {
        cpu.Zero = true
    }

    cpu.testAndSetNegative(cpu.Y)
    cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu) Iny() {
    cpu.Y = cpu.Y + 1

    cpu.testAndSetNegative(cpu.Y)
    cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu) Branch(offset Word) {
    switch {
    case offset < 0x80:
        programCounter += int(offset) + 1
    case offset > 0x7f:
        programCounter -= int(((offset ^ 0xff) + 1))
    }
}

func (cpu *Cpu) Bpl() {
    if !cpu.Negative {
        cpu.Branch(Ram[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Bmi() {
    if cpu.Negative {
        cpu.Branch(Ram[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Bvc() {
    if !cpu.Overflow {
        cpu.Branch(Ram[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Bvs() {
    if cpu.Overflow {
        cpu.Branch(Ram[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Bcc() {
    if !cpu.Carry {
        cpu.Branch(Ram[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Bcs() {
    if cpu.Carry {
        cpu.Branch(Ram[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Bne() {
    if !cpu.Zero {
        cpu.Branch(Ram[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Beq() {
    if cpu.Zero {
        cpu.Branch(Ram[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Txs() {
    Ram[cpu.StackPointer] = cpu.X
}

func (cpu *Cpu) Tsx() {
    cpu.X = Ram[cpu.StackPointer]
}

func (cpu *Cpu) Pha() {
    cpu.StackPointer--
    Ram[cpu.StackPointer] = cpu.A
}

func (cpu *Cpu) Pla() {
    cpu.A = Ram[cpu.StackPointer]
    cpu.StackPointer++

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) ProcessorStatus() (status Word) {
    if cpu.Carry {
        status += 0x1
    }

    if cpu.Zero {
        status += 0x2
    }

    if cpu.IrqDisable {
        status += 0x4
    }

    if cpu.DecimalMode {
        status += 0x8
    }

    if cpu.BrkCommand {
        status += 0x10
    }

    if cpu.Overflow {
        status += 0x40
    }

    if cpu.Negative {
        status += 0x80
    }

    return
}

func (cpu *Cpu) SetProcessorStatus(status Word) {
    if status & 0x1 == 0x1 {
        cpu.Carry = true
    }

    if status & 0x2 == 0x2 {
        cpu.Zero = true
    }

    if status & 0x4 == 0x4 {
        cpu.IrqDisable = true
    }

    if status & 0x8 == 0x8 {
        cpu.DecimalMode = true
    }

    if status & 0x10 == 0x10 {
        cpu.BrkCommand = true
    }

    if status & 0x40 == 0x40 {
        cpu.Overflow = true
    }

    if status & 0x80 == 0x80 {
        cpu.Negative = true
    }
}

func (cpu *Cpu) Php() {
    cpu.StackPointer--
    Ram[cpu.StackPointer] = cpu.ProcessorStatus()
}

func (cpu *Cpu) Plp() {
    cpu.SetProcessorStatus(Ram[cpu.StackPointer])
    cpu.StackPointer++
}

func (cpu *Cpu) Compare(register Word, value Word) {
    switch {
    case register < value:
        cpu.Negative = true
        cpu.Zero = false
        cpu.Carry = false
    case register == value:
        cpu.Negative = false
        cpu.Zero = true
        cpu.Carry = true
    case register > value:
        cpu.Negative = false
        cpu.Zero = false
        cpu.Carry = true
    }
}

func (cpu *Cpu) Cmp(location int) {
    cpu.Compare(cpu.A, Ram[location])
}

func (cpu *Cpu) Cpx(location int) {
    cpu.Compare(cpu.X, Ram[location])
}

func (cpu *Cpu) Cpy(location int) {
    cpu.Compare(cpu.Y, Ram[location])
}

func (cpu *Cpu) Sbc(location int) {
    cpu.A = cpu.A - Ram[location]

    if cpu.Carry {
        cpu.A--
    }

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
    cpu.testAndSetOverflowSubtraction(cpu.A, Ram[location])
    cpu.testAndSetCarrySubtraction(cpu.A, Ram[location])

    cpu.A = cpu.A & 0xff
}

func (cpu *Cpu) Clc() {
    cpu.Carry = false
}

func (cpu *Cpu) Sec() {
    cpu.Carry = true
}

func (cpu *Cpu) Cli() {
    cpu.IrqDisable = false
}

func (cpu *Cpu) Sei() {
    cpu.IrqDisable = true
}

func (cpu *Cpu) Clv() {
    cpu.Overflow = false
}

func (cpu *Cpu) Cld() {
    cpu.DecimalMode = false
}

func (cpu *Cpu) Sed() {
    cpu.DecimalMode = true
}

func (cpu *Cpu) And(location int) {
    cpu.A = cpu.A & Ram[location]

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Ora(location int) {
    cpu.A = cpu.A | Ram[location]

    cpu.testAndSetNegative(Ram[location])
    cpu.testAndSetZero(Ram[location])
}

func (cpu *Cpu) Eor(location int) {
    cpu.A = cpu.A ^ Ram[location]

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Dec(location int) {
    Ram[location] = Ram[location] - 1

    cpu.testAndSetNegative(Ram[location])
    cpu.testAndSetZero(Ram[location])
}

func (cpu *Cpu) Inc(location int) {
    Ram[location] = Ram[location] + 1

    cpu.testAndSetNegative(Ram[location])
    cpu.testAndSetZero(Ram[location])
}

func (cpu *Cpu) Brk() {
    cpu.BrkCommand = true
    programCounter++
}

func (cpu *Cpu) Jsr(location int) {
    cpu.StackPointer--
    Ram[cpu.StackPointer] = Word(programCounter - 1)

    programCounter = location
}

func (cpu *Cpu) Rti() {
    cpu.Plp()

    programCounter = int(Ram[cpu.StackPointer])
    cpu.StackPointer++
}

func (cpu *Cpu) Rts() {
    low := Ram[cpu.StackPointer]
    cpu.StackPointer++

    high := Ram[cpu.StackPointer]
    cpu.StackPointer++

    programCounter = int(((high << 8) + low) + 1)
}

func (cpu *Cpu) Lsr(location int) {
    if Ram[location] & 0x01 > 0x00 {
        cpu.Carry = true
    } else {
        cpu.Carry = false
    }

    Ram[location] = Ram[location] >> 1

    cpu.testAndSetNegative(Ram[location])
    cpu.testAndSetZero(Ram[location])
}

func (cpu *Cpu) LsrAcc() {
    if cpu.A & 0x01 > 0 {
        cpu.Carry = true
    } else {
        cpu.Carry = false
    }

    cpu.A = cpu.A >> 1

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Asl(location int) {
    if Ram[location] & 0x80 > 0 {
        cpu.Carry = true
    } else {
        cpu.Carry = false
    }

    Ram[location] = Ram[location] << 1

    cpu.testAndSetNegative(Ram[location])
    cpu.testAndSetZero(Ram[location])
}

func (cpu *Cpu) AslAcc() {
    if cpu.A & 0x80 > 0 {
        cpu.Carry = true
    } else {
        cpu.Carry = false
    }

    cpu.A = cpu.A << 1

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Rol(location int) {
    value := Ram[location]

    carry := value & 0x80

    value = value << 1

    if cpu.Carry {
        value += 1
    }

    if carry > 0x00 {
        cpu.Carry = true
    } else {
        cpu.Carry = false
    }

    Ram[location] = value

    cpu.testAndSetNegative(Ram[location])
    cpu.testAndSetZero(Ram[location])
}

func (cpu *Cpu) RolAcc() {
    carry := cpu.A & 0x80

    cpu.A = cpu.A << 1

    if cpu.Carry {
        cpu.A += 1
    }

    if carry > 0x00 {
        cpu.Carry = true
    } else {
        cpu.Carry = false
    }

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Ror(location int) {
    value := Ram[location]

    carry := value & 0x1

    value = value >> 1

    if cpu.Carry {
        value += 0x80
    }

    if carry > 0x00 {
        cpu.Carry = true
    } else {
        cpu.Carry = false
    }

    Ram[location] = value

    cpu.testAndSetNegative(Ram[location])
    cpu.testAndSetZero(Ram[location])
}

func (cpu *Cpu) RorAcc() {
    carry := cpu.A & 0x1

    cpu.A = cpu.A >> 1

    if cpu.Carry {
        cpu.A += 0x80
    }

    if carry > 0x00 {
        cpu.Carry = true
    } else {
        cpu.Carry = false
    }

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Bit(location int) {
    if Ram[location] & cpu.A == 0 {
        cpu.Zero = true
    } else {
        cpu.Zero = false
    }

    if Ram[location] & 0x80 > 0x00 {
        cpu.Negative = true
    } else {
        cpu.Negative = false
    }

    if Ram[location] & 0x40 > 0x00 {
        cpu.Overflow = true
    } else {
        cpu.Overflow = false
    }
}

func (cpu *Cpu) Init() {
    cpu.Reset()
}

func (cpu *Cpu) Reset() {
    cpu.X = 0
    cpu.Y = 0
    cpu.A = 0
    cpu.CycleCount = 0
    cpu.Negative = false
    cpu.Overflow = false
    cpu.BrkCommand = false
    cpu.DecimalMode = false
    cpu.IrqDisable = false
    cpu.Zero = false
    cpu.Carry = false
    cpu.StackPointer = 0xff
}

func (cpu *Cpu) DumpState() string {
    return fmt.Sprintf("X: 0x%X Y: 0x%X A: 0x%X PC: 0x%X PPU1: 0x%X PPU2: 0x%X Mem: 0x%X", cpu.X, cpu.Y, cpu.A, programCounter, Ram[0x2000], Ram[0x2001], Ram[0x4014])
}

func (cpu *Cpu) Step() {
    if cpu.CycleCount > 1 && false {
        cpu.CycleCount--
        return
    }

    opcode := Ram[programCounter]

    if cpu.Verbose {
        fmt.Printf("Opcode: 0x%X\n", opcode)
    }

    programCounter++


    switch opcode {
    // ADC
    case 0x69:
        cpu.CycleCount = 2
        cpu.Adc(cpu.immediateAddress())
    case 0x65:
        cpu.CycleCount = 3
        cpu.Adc(cpu.zeroPageAddress())
    case 0x75:
        cpu.CycleCount = 4
        cpu.Adc(cpu.zeroPageIndexedAddress(cpu.X))
    case 0x6D:
        cpu.CycleCount = 4
        cpu.Adc(cpu.absoluteAddress())
    case 0x7D:
        cpu.CycleCount = 4
        cpu.Adc(cpu.absoluteIndexedAddress(cpu.X))
    case 0x79:
        cpu.CycleCount = 4
        cpu.Adc(cpu.absoluteIndexedAddress(cpu.Y))
    case 0x61:
        cpu.CycleCount = 6
        cpu.Adc(cpu.indexedIndirectAddress())
    case 0x71:
        cpu.CycleCount = 5
        cpu.Adc(cpu.indirectIndexedAddress())
    // LDA
    case 0xA9:
        cpu.CycleCount = 2
        cpu.Lda(cpu.immediateAddress())
    case 0xA5:
        cpu.CycleCount = 3
        cpu.Lda(cpu.zeroPageAddress())
    case 0xB5:
        cpu.CycleCount = 4
        cpu.Lda(cpu.zeroPageIndexedAddress(cpu.X))
    case 0xAD:
        cpu.CycleCount = 4
        cpu.Lda(cpu.absoluteAddress())
    case 0xBD:
        cpu.CycleCount = 4
        cpu.Lda(cpu.absoluteIndexedAddress(cpu.X))
    case 0xB9:
        cpu.CycleCount = 4
        cpu.Lda(cpu.absoluteIndexedAddress(cpu.Y))
    case 0xA1:
        cpu.CycleCount = 6
        cpu.Lda(cpu.indexedIndirectAddress())
    case 0xB1:
        cpu.CycleCount = 5
        cpu.Lda(cpu.indirectIndexedAddress())
    // LDX
    case 0xA2:
        cpu.CycleCount = 2
        cpu.Ldx(cpu.immediateAddress())
    case 0xA6:
        cpu.CycleCount = 3
        cpu.Ldx(cpu.zeroPageAddress())
    case 0xB6:
        cpu.CycleCount = 4
        cpu.Ldx(cpu.zeroPageIndexedAddress(cpu.Y))
    case 0xAE:
        cpu.CycleCount = 4
        cpu.Ldx(cpu.absoluteAddress())
    case 0xBE:
        cpu.CycleCount = 4
        cpu.Ldx(cpu.absoluteIndexedAddress(cpu.Y))
    // LDY
    case 0xA0:
        cpu.CycleCount = 2
        cpu.Ldy(cpu.immediateAddress())
    case 0xA4:
        cpu.CycleCount = 3
        cpu.Ldy(cpu.zeroPageAddress())
    case 0xB4:
        cpu.CycleCount = 4
        cpu.Ldy(cpu.zeroPageIndexedAddress(cpu.X))
    case 0xAC:
        cpu.CycleCount = 4
        cpu.Ldy(cpu.absoluteAddress())
    case 0xBC:
        cpu.CycleCount = 4
        cpu.Ldy(cpu.absoluteIndexedAddress(cpu.X))
    // STA
    case 0x85:
        cpu.CycleCount = 3
        cpu.Sta(cpu.zeroPageAddress())
    case 0x95:
        cpu.CycleCount = 4
        cpu.Sta(cpu.zeroPageIndexedAddress(cpu.X))
    case 0x8D:
        cpu.CycleCount = 4
        cpu.Sta(cpu.absoluteAddress())
    case 0x9D:
        cpu.CycleCount = 5
        cpu.Sta(cpu.absoluteIndexedAddress(cpu.X))
    case 0x99:
        cpu.CycleCount = 5
        cpu.Sta(cpu.absoluteIndexedAddress(cpu.Y))
    case 0x81:
        cpu.CycleCount = 6
        cpu.Sta(cpu.indexedIndirectAddress())
    case 0x91:
        cpu.CycleCount = 6
        cpu.Sta(cpu.indirectIndexedAddress())
    // STX
    case 0x86:
        cpu.CycleCount = 3
        cpu.Stx(cpu.zeroPageAddress())
    case 0x96:
        cpu.CycleCount = 4
        cpu.Stx(cpu.zeroPageIndexedAddress(cpu.Y))
    case 0x8E:
        cpu.CycleCount = 4
        cpu.Stx(cpu.absoluteAddress())
    // STY
    case 0x84:
        cpu.CycleCount = 3
        cpu.Sty(cpu.zeroPageAddress())
    case 0x94:
        cpu.CycleCount = 4
        cpu.Sty(cpu.zeroPageIndexedAddress(cpu.X))
    case 0x8C:
        cpu.CycleCount = 4
        cpu.Sty(cpu.absoluteAddress())
    // JMP
    case 0x4C:
        cpu.CycleCount = 3
        cpu.Jmp(cpu.absoluteAddress())
    case 0x6C:
        cpu.CycleCount = 5
        cpu.Jmp(cpu.indirectAbsoluteAddress())
    // JSR
    case 0x20:
        cpu.CycleCount = 6
        cpu.Jsr(cpu.absoluteAddress())
    // Register Instructions
    case 0xAA:
        cpu.CycleCount = 2
        cpu.Tax()
    case 0x8A:
        cpu.CycleCount = 2
        cpu.Txa()
    case 0xCA:
        cpu.CycleCount = 2
        cpu.Dex()
    case 0xE8:
        cpu.CycleCount = 2
        cpu.Inx()
    case 0xA8:
        cpu.CycleCount = 2
        cpu.Tay()
    case 0x98:
        cpu.CycleCount = 2
        cpu.Tya()
    case 0x88:
        cpu.CycleCount = 2
        cpu.Dey()
    case 0xC8:
        cpu.CycleCount = 2
        cpu.Iny()
    // Branch Instructions
    case 0x10:
        cpu.CycleCount = 2
        cpu.Bpl()
    case 0x30:
        cpu.CycleCount = 2
        cpu.Bmi()
    case 0x50:
        cpu.CycleCount = 2
        cpu.Bvc()
    case 0x70:
        cpu.CycleCount = 2
        cpu.Bvs()
    case 0x90:
        cpu.CycleCount = 2
        cpu.Bcc()
    case 0xB0:
        cpu.CycleCount = 2
        cpu.Bcs()
    case 0xD0:
        cpu.CycleCount = 2
        cpu.Bne()
    case 0xF0:
        cpu.CycleCount = 2
        cpu.Beq()
    // CMP
    case 0xC9:
        cpu.CycleCount = 2
        cpu.Cmp(cpu.immediateAddress())
    case 0xC5:
        cpu.CycleCount = 3
        cpu.Cmp(cpu.zeroPageAddress())
    case 0xD5:
        cpu.CycleCount = 4
        cpu.Cmp(cpu.zeroPageIndexedAddress(cpu.X))
    case 0xCD:
        cpu.CycleCount = 4
        cpu.Cmp(cpu.absoluteAddress())
    case 0xDD:
        cpu.CycleCount = 4
        cpu.Cmp(cpu.absoluteIndexedAddress(cpu.X))
    case 0xD9:
        cpu.CycleCount = 4
        cpu.Cmp(cpu.absoluteIndexedAddress(cpu.Y))
    case 0xC1:
        cpu.CycleCount = 6
        cpu.Cmp(cpu.indexedIndirectAddress())
    case 0xD1:
        cpu.CycleCount = 5
        cpu.Cmp(cpu.indirectIndexedAddress())
    // CPX
    case 0xE0:
        cpu.CycleCount = 2
        cpu.Cpx(cpu.immediateAddress())
    case 0xE4:
        cpu.CycleCount = 3
        cpu.Cpx(cpu.zeroPageAddress())
    case 0xEC:
        cpu.CycleCount = 4
        cpu.Cpx(cpu.absoluteAddress())
    // CPY
    case 0xC0:
        cpu.CycleCount = 2
        cpu.Cpy(cpu.immediateAddress())
    case 0xC4:
        cpu.CycleCount = 3
        cpu.Cpy(cpu.zeroPageAddress())
    case 0xCC:
        cpu.CycleCount = 4
        cpu.Cpy(cpu.absoluteAddress())
    // SBC
    case 0xE9:
        cpu.CycleCount = 2
        cpu.Sbc(cpu.immediateAddress())
    case 0xE5:
        cpu.CycleCount = 3
        cpu.Sbc(cpu.zeroPageAddress())
    case 0xF5:
        cpu.CycleCount = 4
        cpu.Sbc(cpu.zeroPageIndexedAddress(cpu.X))
    case 0xED:
        cpu.CycleCount = 4
        cpu.Sbc(cpu.absoluteAddress())
    case 0xFD:
        cpu.CycleCount = 4
        cpu.Sbc(cpu.absoluteIndexedAddress(cpu.X))
    case 0xF9:
        cpu.CycleCount = 4
        cpu.Sbc(cpu.absoluteIndexedAddress(cpu.Y))
    case 0xE1:
        cpu.CycleCount = 6
        cpu.Sbc(cpu.indexedIndirectAddress())
    case 0xF1:
        cpu.CycleCount = 5
        cpu.Sbc(cpu.indirectIndexedAddress())
    // Flag Instructions
    case 0x18:
        cpu.CycleCount = 2
        cpu.Clc()
    case 0x38:
        cpu.CycleCount = 2
        cpu.Sec()
    case 0x58:
        cpu.CycleCount = 2
        cpu.Cli()
    case 0x78:
        cpu.CycleCount = 2
        cpu.Sei()
    case 0xb8:
        cpu.CycleCount = 2
        cpu.Clv()
    case 0xd8:
        cpu.CycleCount = 2
        cpu.Cld()
    case 0xf8:
        cpu.CycleCount = 2
        cpu.Sed()
    // Stack instructions
    case 0x9a:
        cpu.CycleCount = 2
        cpu.Txs()
    case 0xba:
        cpu.CycleCount = 2
        cpu.Tsx()
    case 0x48:
        cpu.CycleCount = 3
        cpu.Pha()
    case 0x68:
        cpu.CycleCount = 4
        cpu.Pla()
    case 0x08:
        cpu.CycleCount = 3
        cpu.Php()
    case 0x28:
        cpu.CycleCount = 4
        cpu.Plp()
    // AND
    case 0x29:
        cpu.CycleCount = 2
        cpu.And(cpu.immediateAddress())
    case 0x25:
        cpu.CycleCount = 3
        cpu.And(cpu.zeroPageAddress())
    case 0x35:
        cpu.CycleCount = 4
        cpu.And(cpu.zeroPageIndexedAddress(cpu.X))
    case 0x2d:
        cpu.CycleCount = 4
        cpu.And(cpu.absoluteAddress())
    case 0x3d:
        cpu.CycleCount = 4
        cpu.And(cpu.absoluteIndexedAddress(cpu.X))
    case 0x39:
        cpu.CycleCount = 4
        cpu.And(cpu.absoluteIndexedAddress(cpu.Y))
    case 0x21:
        cpu.CycleCount = 6
        cpu.And(cpu.indexedIndirectAddress())
    case 0x31:
        cpu.CycleCount = 5
        cpu.And(cpu.indirectIndexedAddress())
    // ORA
    case 0x09:
        cpu.CycleCount = 2
        cpu.Ora(cpu.immediateAddress())
    case 0x05:
        cpu.CycleCount = 3
        cpu.Ora(cpu.zeroPageAddress())
    case 0x15:
        cpu.CycleCount = 4
        cpu.Ora(cpu.zeroPageIndexedAddress(cpu.X))
    case 0x0d:
        cpu.CycleCount = 4
        cpu.Ora(cpu.absoluteAddress())
    case 0x1d:
        cpu.CycleCount = 4
        cpu.Ora(cpu.absoluteIndexedAddress(cpu.X))
    case 0x19:
        cpu.CycleCount = 4
        cpu.Ora(cpu.absoluteIndexedAddress(cpu.Y))
    case 0x01:
        cpu.CycleCount = 6
        cpu.Ora(cpu.indexedIndirectAddress())
    case 0x11:
        cpu.CycleCount = 5
        cpu.Ora(cpu.indirectIndexedAddress())
    // EOR
    case 0x49:
        cpu.CycleCount = 2
        cpu.Eor(cpu.immediateAddress())
    case 0x45:
        cpu.CycleCount = 3
        cpu.Eor(cpu.zeroPageAddress())
    case 0x55:
        cpu.CycleCount = 4
        cpu.Eor(cpu.zeroPageIndexedAddress(cpu.X))
    case 0x4d:
        cpu.CycleCount = 4
        cpu.Eor(cpu.absoluteAddress())
    case 0x5d:
        cpu.CycleCount = 4
        cpu.Eor(cpu.absoluteIndexedAddress(cpu.X))
    case 0x59:
        cpu.CycleCount = 4
        cpu.Eor(cpu.absoluteIndexedAddress(cpu.Y))
    case 0x41:
        cpu.CycleCount = 6
        cpu.Eor(cpu.indexedIndirectAddress())
    case 0x51:
        cpu.CycleCount = 5
        cpu.Eor(cpu.indirectIndexedAddress())
    // DEC
    case 0xc6:
        cpu.CycleCount = 5
        cpu.Dec(cpu.zeroPageAddress())
    case 0xd6:
        cpu.CycleCount = 6
        cpu.Dec(cpu.zeroPageIndexedAddress(cpu.X))
    case 0xce:
        cpu.CycleCount = 6
        cpu.Dec(cpu.absoluteAddress())
    case 0xde:
        cpu.CycleCount = 7
        cpu.Dec(cpu.absoluteIndexedAddress(cpu.X))
    // INC
    case 0xe6:
        cpu.CycleCount = 5
        cpu.Inc(cpu.zeroPageAddress())
    case 0xf6:
        cpu.CycleCount = 6
        cpu.Inc(cpu.zeroPageIndexedAddress(cpu.X))
    case 0xee:
        cpu.CycleCount = 6
        cpu.Inc(cpu.absoluteAddress())
    case 0xfe:
        cpu.CycleCount = 7
        cpu.Inc(cpu.absoluteIndexedAddress(cpu.X))
    // BRK
    case 0x00:
        cpu.CycleCount = 7
        cpu.Brk()
    // RTI
    case 0x40:
        cpu.CycleCount = 6
        cpu.Rti()
    // RTS
    case 0x60:
        cpu.CycleCount = 6
        cpu.Rts()
    // NOP
    case 0xea:
        cpu.CycleCount = 2
    // LSR
    case 0x4a:
        cpu.CycleCount = 2
        cpu.LsrAcc()
    case 0x46:
        cpu.CycleCount = 5
        cpu.Lsr(cpu.zeroPageAddress())
    case 0x56:
        cpu.CycleCount = 6
        cpu.Lsr(cpu.zeroPageIndexedAddress(cpu.X))
    case 0x4e:
        cpu.CycleCount = 6
        cpu.Lsr(cpu.absoluteAddress())
    case 0x5e:
        cpu.CycleCount = 7
        cpu.Lsr(cpu.absoluteIndexedAddress(cpu.X))
    // ASL
    case 0x0a:
        cpu.CycleCount = 2
        cpu.AslAcc()
    case 0x06:
        cpu.CycleCount = 5
        cpu.Asl(cpu.zeroPageAddress())
    case 0x16:
        cpu.CycleCount = 6
        cpu.Asl(cpu.zeroPageIndexedAddress(cpu.X))
    case 0x0e:
        cpu.CycleCount = 6
        cpu.Asl(cpu.absoluteAddress())
    case 0x1e:
        cpu.CycleCount = 7
        cpu.Asl(cpu.absoluteIndexedAddress(cpu.X))
    // ROL
    case 0x2a:
        cpu.CycleCount = 2
        cpu.RolAcc()
    case 0x26:
        cpu.CycleCount = 5
        cpu.Rol(cpu.zeroPageAddress())
    case 0x36:
        cpu.CycleCount = 6
        cpu.Rol(cpu.zeroPageIndexedAddress(cpu.X))
    case 0x2e:
        cpu.CycleCount = 6
        cpu.Rol(cpu.absoluteAddress())
    case 0x3e:
        cpu.CycleCount = 7
        cpu.Rol(cpu.absoluteIndexedAddress(cpu.X))
    // ROR
    case 0x6a:
        cpu.CycleCount = 2
        cpu.RorAcc()
    case 0x66:
        cpu.CycleCount = 5
        cpu.Ror(cpu.zeroPageAddress())
    case 0x76:
        cpu.CycleCount = 6
        cpu.Ror(cpu.zeroPageIndexedAddress(cpu.X))
    case 0x6e:
        cpu.CycleCount = 6
        cpu.Ror(cpu.absoluteAddress())
    case 0x7e:
        cpu.CycleCount = 7
        cpu.Ror(cpu.absoluteIndexedAddress(cpu.X))
    // BIT
    case 0x24:
        cpu.CycleCount = 3
        cpu.Bit(cpu.zeroPageAddress())
    case 0x2c:
        cpu.CycleCount = 4
        cpu.Bit(cpu.absoluteAddress())
    }
}
