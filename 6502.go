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
    high := int(memory[programCounter + 1])
    low := int(memory[programCounter])

    programCounter += 2
    return (high << 8) + low
}

func (cpu *Cpu) zeroPageAddress() (int) {
    programCounter++
    return int(memory[programCounter - 1])
}

func (cpu *Cpu) indirectAbsoluteAddress() (result int) {
    result = int((memory[programCounter + 1] << 8) + memory[programCounter])
    programCounter++
    return
}

func (cpu *Cpu) absoluteIndexedAddress(index Word) (result int) {
    // Switch to an int (or more appropriately uint16) since we 
    // will overflow when shifting the high byte
    high := int(memory[programCounter + 1])
    low := int(memory[programCounter])

    programCounter++
    return (high << 8) + low + int(index)
}

func (cpu *Cpu) zeroPageIndexedAddress(index Word) (int) {
    location := int(memory[programCounter] + index)
    programCounter++
    return location
}

func (cpu *Cpu) indexedIndirectAddress() (int) {
    location := memory[programCounter] + cpu.X
    programCounter++

    // Switch to an int (or more appropriately uint16) since we 
    // will overflow when shifting the high byte
    high := int(memory[location + 1])
    low := int(memory[location])

    return (high << 8) + low
}

func (cpu *Cpu) indirectIndexedAddress() (int) {
    location := memory[programCounter]

    // Switch to an int (or more appropriately uint16) since we 
    // will overflow when shifting the high byte
    high := int(memory[location + 1])
    low := int(memory[location])

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
    cpu.A = cpu.A + memory[location]

    if cpu.Carry {
        cpu.A++
    }

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
    cpu.testAndSetOverflowAddition(cpu.A, memory[location])
    cpu.testAndSetCarryAddition(cpu.A, memory[location])
}

func (cpu *Cpu) Lda(location int) {
    cpu.A = memory[location]

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Ldx(location int) {
    cpu.X = memory[location]

    cpu.testAndSetNegative(cpu.X)
    cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu) Ldy(location int) {
    cpu.Y = memory[location]

    cpu.testAndSetNegative(cpu.Y)
    cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu) Sta(location int) {
    memory[location] = cpu.A
}

func (cpu *Cpu) Stx(location int) {
    memory[location] = cpu.X
}

func (cpu *Cpu) Sty(location int) {
    memory[location] = cpu.Y
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
        cpu.Branch(memory[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Bmi() {
    if cpu.Negative {
        cpu.Branch(memory[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Bvc() {
    if !cpu.Overflow {
        cpu.Branch(memory[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Bvs() {
    if cpu.Overflow {
        cpu.Branch(memory[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Bcc() {
    if !cpu.Carry {
        cpu.Branch(memory[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Bcs() {
    if cpu.Carry {
        cpu.Branch(memory[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Bne() {
    if !cpu.Zero {
        cpu.Branch(memory[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Beq() {
    if cpu.Zero {
        cpu.Branch(memory[programCounter])
    } else {
        programCounter++
    }
}

func (cpu *Cpu) Txs() {
    memory[cpu.StackPointer] = cpu.X
}

func (cpu *Cpu) Tsx() {
    cpu.X = memory[cpu.StackPointer]
}

func (cpu *Cpu) Pha() {
    cpu.StackPointer--
    memory[cpu.StackPointer] = cpu.A
}

func (cpu *Cpu) Pla() {
    cpu.A = memory[cpu.StackPointer]
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
    memory[cpu.StackPointer] = cpu.ProcessorStatus()
}

func (cpu *Cpu) Plp() {
    cpu.SetProcessorStatus(memory[cpu.StackPointer])
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
    cpu.Compare(cpu.A, memory[location])
}

func (cpu *Cpu) Cpx(location int) {
    cpu.Compare(cpu.X, memory[location])
}

func (cpu *Cpu) Cpy(location int) {
    cpu.Compare(cpu.Y, memory[location])
}

func (cpu *Cpu) Sbc(location int) {
    cpu.A = cpu.A - memory[location]

    if cpu.Carry {
        cpu.A--
    }

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
    cpu.testAndSetOverflowSubtraction(cpu.A, memory[location])
    cpu.testAndSetCarrySubtraction(cpu.A, memory[location])
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
    cpu.A = cpu.A & memory[location]

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Ora(location int) {
    cpu.A = cpu.A | memory[location]

    cpu.testAndSetNegative(memory[location])
    cpu.testAndSetZero(memory[location])
}

func (cpu *Cpu) Eor(location int) {
    cpu.A = cpu.A ^ memory[location]

    cpu.testAndSetNegative(cpu.A)
    cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Dec(location int) {
    memory[location] = memory[location] - 1

    cpu.testAndSetNegative(memory[location])
    cpu.testAndSetZero(memory[location])
}

func (cpu *Cpu) Inc(location int) {
    memory[location] = memory[location] + 1

    cpu.testAndSetNegative(memory[location])
    cpu.testAndSetZero(memory[location])
}

func (cpu *Cpu) Brk() {
    cpu.BrkCommand = true
    programCounter++
}

func (cpu *Cpu) Jsr(location int) {
    cpu.StackPointer--
    memory[cpu.StackPointer] = Word(programCounter - 1)

    programCounter = location
}

func (cpu *Cpu) Rti() {
    cpu.Plp()

    programCounter = int(memory[cpu.StackPointer])
    cpu.StackPointer++
}

func (cpu *Cpu) Rts() {
    low := memory[cpu.StackPointer]
    cpu.StackPointer++

    high := memory[cpu.StackPointer]
    cpu.StackPointer++

    programCounter = int(((high << 8) + low) + 1)
}

func (cpu *Cpu) Lsr(location int) {
    if memory[location] & 0x01 > 0x00 {
        cpu.Carry = true
    } else {
        cpu.Carry = false
    }

    memory[location] = memory[location] >> 1

    cpu.testAndSetNegative(memory[location])
    cpu.testAndSetZero(memory[location])
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
    if memory[location] & 0x80 > 0 {
        cpu.Carry = true
    } else {
        cpu.Carry = false
    }

    memory[location] = memory[location] << 1

    cpu.testAndSetNegative(memory[location])
    cpu.testAndSetZero(memory[location])
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
    value := memory[location]

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

    memory[location] = value

    cpu.testAndSetNegative(memory[location])
    cpu.testAndSetZero(memory[location])
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
    value := memory[location]

    carry := value & 0x80

    value = value >> 1

    if cpu.Carry {
        value += 1
    }

    if carry > 0x00 {
        cpu.Carry = true
    } else {
        cpu.Carry = false
    }

    memory[location] = value

    cpu.testAndSetNegative(memory[location])
    cpu.testAndSetZero(memory[location])
}

func (cpu *Cpu) RorAcc() {
    carry := cpu.A & 0x01

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
    if memory[location] & cpu.A == 0 {
        cpu.Zero = true
    } else {
        cpu.Zero = false
    }

    if memory[location] & 0x80 > 0x00 {
        cpu.Negative = true
    } else {
        cpu.Negative = false
    }

    if memory[location] & 0x40 > 0x00 {
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

func (cpu *Cpu) Step() {
    fmt.Printf("X: 0x%x Y: 0x%x A: 0x%x PC: %d\n", cpu.X, cpu.Y, cpu.A, programCounter)

    if cpu.CycleCount > 1 && false {
        cpu.CycleCount--
        return
    }

    opcode := memory[programCounter]
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
    case 0x6d:
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
    case 0xa9:
        cpu.CycleCount = 2
        cpu.Lda(cpu.immediateAddress())
    case 0xa5:
        cpu.CycleCount = 3
        cpu.Lda(cpu.zeroPageAddress())
    case 0xb5:
        cpu.CycleCount = 4
        cpu.Lda(cpu.zeroPageIndexedAddress(cpu.X))
    case 0xad:
        cpu.CycleCount = 4
        cpu.Lda(cpu.absoluteAddress())
    case 0xbd:
        cpu.CycleCount = 4
        cpu.Lda(cpu.absoluteIndexedAddress(cpu.X))
    case 0xb9:
        cpu.CycleCount = 4
        cpu.Lda(cpu.absoluteIndexedAddress(cpu.Y))
    case 0xa1:
        cpu.CycleCount = 6
        cpu.Lda(cpu.indexedIndirectAddress())
    case 0xb1:
        cpu.CycleCount = 5
        cpu.Lda(cpu.indirectIndexedAddress())
    // LDX
    case 0xa2:
        cpu.CycleCount = 2
        cpu.Ldx(cpu.immediateAddress())
    case 0xa6:
        cpu.CycleCount = 3
        cpu.Ldx(cpu.zeroPageAddress())
    case 0xb6:
        cpu.CycleCount = 4
        cpu.Ldx(cpu.zeroPageIndexedAddress(cpu.Y))
    case 0xae:
        cpu.CycleCount = 4
        cpu.Ldx(cpu.absoluteAddress())
    case 0xbe:
        cpu.CycleCount = 4
        cpu.Ldx(cpu.absoluteIndexedAddress(cpu.Y))
    // LDY
    case 0xa0:
        cpu.CycleCount = 2
        cpu.Ldy(cpu.immediateAddress())
    case 0xa4:
        cpu.CycleCount = 3
        cpu.Ldy(cpu.zeroPageAddress())
    case 0xb4:
        cpu.CycleCount = 4
        cpu.Ldy(cpu.zeroPageIndexedAddress(cpu.X))
    case 0xac:
        cpu.CycleCount = 4
        cpu.Ldy(cpu.absoluteAddress())
    case 0xbc:
        cpu.CycleCount = 4
        cpu.Ldy(cpu.absoluteIndexedAddress(cpu.X))
    // STA
    case 0x85:
        cpu.CycleCount = 3
        cpu.Sta(cpu.zeroPageAddress())
    case 0x95:
        cpu.CycleCount = 4
        cpu.Sta(cpu.zeroPageIndexedAddress(cpu.X))
    case 0x8d:
        cpu.CycleCount = 4
        cpu.Sta(cpu.absoluteAddress())
    case 0x9d:
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
    case 0x8e:
        cpu.CycleCount = 4
        cpu.Stx(cpu.absoluteAddress())
    // STY
    case 0x84:
        cpu.CycleCount = 3
        cpu.Sty(cpu.zeroPageAddress())
    case 0x94:
        cpu.CycleCount = 4
        cpu.Sty(cpu.zeroPageIndexedAddress(cpu.X))
    case 0x8c:
        cpu.CycleCount = 4
        cpu.Sty(cpu.absoluteAddress())
    // JMP
    case 0x4c:
        cpu.CycleCount = 3
        cpu.Jmp(cpu.absoluteAddress())
    case 0x6c:
        cpu.CycleCount = 5
        cpu.Jmp(cpu.indirectAbsoluteAddress())
    // JSR
    case 0x20:
        cpu.CycleCount = 6
        cpu.Jsr(cpu.absoluteAddress())
    // Register Instructions
    case 0xaa:
        cpu.CycleCount = 2
        cpu.Tax()
    case 0x8a:
        cpu.CycleCount = 2
        cpu.Txa()
    case 0xca:
        cpu.CycleCount = 2
        cpu.Dex()
    case 0xe8:
        cpu.CycleCount = 2
        cpu.Inx()
    case 0xa8:
        cpu.CycleCount = 2
        cpu.Tay()
    case 0x98:
        cpu.CycleCount = 2
        cpu.Tya()
    case 0x88:
        cpu.CycleCount = 2
        cpu.Dey()
    case 0xc8:
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
    case 0xb0:
        cpu.CycleCount = 2
        cpu.Bcs()
    case 0xd0:
        cpu.CycleCount = 2
        cpu.Bne()
    case 0xf0:
        cpu.CycleCount = 2
        cpu.Beq()
    // CMP
    case 0xc9:
        cpu.CycleCount = 2
        cpu.Cmp(cpu.immediateAddress())
    case 0xc5:
        cpu.CycleCount = 3
        cpu.Cmp(cpu.zeroPageAddress())
    case 0xd5:
        cpu.CycleCount = 4
        cpu.Cmp(cpu.zeroPageIndexedAddress(cpu.X))
    case 0xcd:
        cpu.CycleCount = 4
        cpu.Cmp(cpu.absoluteAddress())
    case 0xdd:
        cpu.CycleCount = 4
        cpu.Cmp(cpu.absoluteIndexedAddress(cpu.X))
    case 0xd9:
        cpu.CycleCount = 4
        cpu.Cmp(cpu.absoluteIndexedAddress(cpu.Y))
    case 0xc1:
        cpu.CycleCount = 6
        cpu.Cmp(cpu.indexedIndirectAddress())
    case 0xd1:
        cpu.CycleCount = 5
        cpu.Cmp(cpu.indirectIndexedAddress())
    // CPX
    case 0xe0:
        cpu.CycleCount = 2
        cpu.Cpx(cpu.immediateAddress())
    case 0xe4:
        cpu.CycleCount = 3
        cpu.Cpx(cpu.zeroPageAddress())
    case 0xec:
        cpu.CycleCount = 4
        cpu.Cpx(cpu.absoluteAddress())
    // CPY
    case 0xc0:
        cpu.CycleCount = 2
        cpu.Cpy(cpu.immediateAddress())
    case 0xc4:
        cpu.CycleCount = 3
        cpu.Cpy(cpu.zeroPageAddress())
    case 0xcc:
        cpu.CycleCount = 4
        cpu.Cpy(cpu.absoluteAddress())
    // SBC
    case 0xe9:
        cpu.CycleCount = 2
        cpu.Sbc(cpu.immediateAddress())
    case 0xe5:
        cpu.CycleCount = 3
        cpu.Sbc(cpu.zeroPageAddress())
    case 0xf5:
        cpu.CycleCount = 4
        cpu.Sbc(cpu.zeroPageIndexedAddress(cpu.X))
    case 0xed:
        cpu.CycleCount = 4
        cpu.Sbc(cpu.absoluteAddress())
    case 0xfd:
        cpu.CycleCount = 4
        cpu.Sbc(cpu.absoluteIndexedAddress(cpu.X))
    case 0xf9:
        cpu.CycleCount = 4
        cpu.Sbc(cpu.absoluteIndexedAddress(cpu.Y))
    case 0xe1:
        cpu.CycleCount = 6
        cpu.Sbc(cpu.indexedIndirectAddress())
    case 0xf1:
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
