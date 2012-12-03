package main

import (
	"log"
)

const (
	InterruptNone = iota
	InterruptIrq
	InterruptReset
	InterruptNmi
)

var (
	ProgramCounter uint16
)

type Cpu struct {
	X            Word
	Y            Word
	A            Word
	P            Word
	CycleCount   int
	StackPointer Word
	Opcode       Word
	Verbose      bool
	Accurate     bool

	InterruptRequested int
	CyclesToWait       int
	Timestamp          int
}

func (c *Cpu) getCarry() bool {
	return c.P&0x01 == 0x01
}

func (c *Cpu) getZero() bool {
	return c.P&0x02 == 0x02
}

func (c *Cpu) getIrqDisable() bool {
	return c.P&0x04 == 0x04
}

func (c *Cpu) getDecimalMode() bool {
	return c.P&0x08 == 0x08
}

func (c *Cpu) getBrkCommand() bool {
	return c.P&0x10 == 0x10
}

func (c *Cpu) getOverflow() bool {
	return c.P&0x40 == 0x40
}

func (c *Cpu) getNegative() bool {
	return c.P&0x80 == 0x80
}

func (c *Cpu) setCarry() {
	c.P = c.P | 0x01
}

func (c *Cpu) setZero() {
	c.P = c.P | 0x02
}

func (c *Cpu) setIrqDisable() {
	c.P = c.P | 0x04
}

func (c *Cpu) setDecimalMode() {
	c.P = c.P | 0x08
}

func (c *Cpu) setBrkCommand() {
	c.P = c.P | 0x10
}

func (c *Cpu) setOverflow() {
	c.P = c.P | 0x40
}

func (c *Cpu) setNegative() {
	c.P = c.P | 0x80
}

func (c *Cpu) clearCarry() {
	c.P = c.P & 0xFE
}

func (c *Cpu) clearZero() {
	c.P = c.P & 0xFD
}

func (c *Cpu) clearIrqDisable() {
	c.P = c.P & 0xFB
}

func (c *Cpu) clearDecimalMode() {
	c.P = c.P & 0xF7
}

func (c *Cpu) clearBrkCommand() {
	c.P = c.P & 0xEF
}

func (c *Cpu) clearOverflow() {
	c.P = c.P & 0xBF
}

func (c *Cpu) clearNegative() {
	c.P = c.P & 0x7F
}

func (c *Cpu) pushToStack(value Word) {
	Ram.Write(0x100+int(c.StackPointer), value)
	c.StackPointer--
}

func (c *Cpu) pullFromStack() Word {
	c.StackPointer++
	val, _ := Ram.Read(0x100 + int(c.StackPointer))

	return val
}

func (c *Cpu) testAndSetNegative(value Word) {
	if value&0x80 == 0x80 {
		c.setNegative()
		return
	}

	c.clearNegative()
}

func (c *Cpu) testAndSetZero(value Word) {
	if value == 0x00 {
		c.setZero()
		return
	}

	c.clearZero()
}

func (c *Cpu) testAndSetCarryAddition(result int) {
	if result > 0xFF {
		c.setCarry()
		return
	}

	c.clearCarry()
}

func (c *Cpu) testAndSetCarrySubtraction(result int) {
	if result < 0x00 {
		c.clearCarry()
		return
	}

	c.setCarry()
}

func (c *Cpu) testAndSetOverflowAddition(a Word, b Word, r Word) {
	if ((a^b)&0x80 == 0x0) && ((a^r)&0x80 == 0x80) {
		c.setOverflow()
	} else {
		c.clearOverflow()
	}
}

func (c *Cpu) testAndSetOverflowSubtraction(a Word, b Word) {
	val := a - b - (1 - c.P&0x01)
	if ((a^val)&0x80) != 0 && ((a^b)&0x80) != 0 {
		c.setOverflow()
	} else {
		c.clearOverflow()
	}
}

func (c *Cpu) immediateAddress() uint16 {
	ProgramCounter++
	return ProgramCounter - 1
}

func (c *Cpu) absoluteAddress() (result uint16) {
	// Switch to an int (or more appropriately uint16) since we 
	// will overflow when shifting the high byte
	high, _ := Ram.Read(ProgramCounter + 1)
	low, _ := Ram.Read(ProgramCounter)

	ProgramCounter += 2
	return (uint16(high) << 8) + uint16(low)
}

func (c *Cpu) zeroPageAddress() uint16 {
	ProgramCounter++
	res, _ := Ram.Read(ProgramCounter - 1)

	return uint16(res)
}

func (c *Cpu) indirectAbsoluteAddress(addr uint16) (result uint16) {
	high, _ := Ram.Read(addr + 1)
	low, _ := Ram.Read(addr)

	// Indirect jump is bugged on the 6502, it doesn't add 1 to 
	// the full 16-bit value when it reads the second byte, it 
	// adds 1 to the low byte only. So JMP (03FF) reads from 3FF 
	// and 300, not 3FF and 400.
	laddr := (uint16(high) << 8) + uint16(low)
	haddr := (uint16(high) << 8) + ((uint16(low) + 1) & 0xFF)

	ih, _ := Ram.Read(haddr)
	il, _ := Ram.Read(laddr)

	result = (uint16(ih) << 8) + uint16(il)
	return
}

func (c *Cpu) absoluteIndexedAddress(index Word) (result uint16) {
	// Switch to an int (or more appropriately uint16) since we 
	// will overflow when shifting the high byte
	high, _ := Ram.Read(ProgramCounter + 1)
	low, _ := Ram.Read(ProgramCounter)

	address := (uint16(high) << 8) + uint16(low)

	if address&0xFF00 != (address+uint16(index))&0xFF00 {
		c.CycleCount += 1
	}

	address += uint16(index)

	if address > 0xFFFF {
		address = address & 0xFFFF
	}

	ProgramCounter += 2
	return address
}

func (c *Cpu) zeroPageIndexedAddress(index Word) uint16 {
	location, _ := Ram.Read(ProgramCounter)
	ProgramCounter++
	return uint16(location + index)
}

func (c *Cpu) indexedIndirectAddress() uint16 {
	location, _ := Ram.Read(ProgramCounter)
	location = location + c.X

	ProgramCounter++

	// Switch to an int (or more appropriately uint16) since we 
	// will overflow when shifting the high byte
	high, _ := Ram.Read(location + 1)
	low, _ := Ram.Read(location)

	return (uint16(high) << 8) + uint16(low)
}

func (c *Cpu) indirectIndexedAddress() uint16 {
	location, _ := Ram.Read(ProgramCounter)

	// Switch to an int (or more appropriately uint16) since we 
	// will overflow when shifting the high byte
	high, _ := Ram.Read(location + 1)
	low, _ := Ram.Read(location)

	address := (uint16(high) << 8) + uint16(low)

	if address&0xFF00 != (address+uint16(c.Y))&0xFF00 {
		c.CycleCount += 1
	}

	address += uint16(c.Y)

	if address > 0xFFFF {
		address = address & 0xFFFF
	}

	ProgramCounter++
	return address
}

func (c *Cpu) relativeAddress() (a uint16) {
	val, _ := Ram.Read(ProgramCounter)

	a = uint16(val)
	if a < 0x80 {
		a = a + ProgramCounter
	} else {
		a = a + (ProgramCounter - 0x100)
	}

	a++

	return
}

func (c *Cpu) accumulatorAddress() uint16 {
	return 0
}

func (c *Cpu) Adc(location uint16) {
	val, _ := Ram.Read(location)

	cached := c.A

	c.A = cached + val + (c.P & 0x01)

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
	c.testAndSetOverflowAddition(cached, val, cpu.A)
	c.testAndSetCarryAddition(int(cached) + int(val) + int(c.P&0x01))

	c.A = c.A & 0xFF
}

func (c *Cpu) Lda(location uint16) {
	val, _ := Ram.Read(location)
	c.A = val

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
}

func (c *Cpu) Ldx(location uint16) {
	val, _ := Ram.Read(location)
	c.X = val

	c.testAndSetNegative(c.X)
	c.testAndSetZero(c.X)
}

func (c *Cpu) Ldy(location uint16) {
	val, _ := Ram.Read(location)
	c.Y = val

	c.testAndSetNegative(c.Y)
	c.testAndSetZero(c.Y)
}

func (c *Cpu) Sta(location uint16) {
	Ram.Write(location, c.A)
}

func (c *Cpu) Stx(location uint16) {
	Ram.Write(location, c.X)
}

func (c *Cpu) Sty(location uint16) {
	Ram.Write(location, c.Y)
}

func (c *Cpu) Jmp(location uint16) {
	ProgramCounter = location
}

func (c *Cpu) Tax() {
	c.X = c.A

	c.testAndSetNegative(c.X)
	c.testAndSetZero(c.X)
}

func (c *Cpu) Txa() {
	c.A = c.X

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
}

func (c *Cpu) Dex() {
	c.X = c.X - 1

	if c.X == 0 {
		c.setZero()
	}

	c.testAndSetNegative(c.X)
	c.testAndSetZero(c.X)
}

func (c *Cpu) Inx() {
	c.X = c.X + 1

	c.testAndSetNegative(c.X)
	c.testAndSetZero(c.X)
}

func (c *Cpu) Tay() {
	c.Y = c.A

	c.testAndSetNegative(c.Y)
	c.testAndSetZero(c.Y)
}

func (c *Cpu) Tya() {
	c.A = c.Y

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
}

func (c *Cpu) Dey() {
	c.Y = c.Y - 1

	if c.X == 0 {
		c.setZero()
	}

	c.testAndSetNegative(c.Y)
	c.testAndSetZero(c.Y)
}

func (c *Cpu) Iny() {
	c.Y = c.Y + 1

	c.testAndSetNegative(c.Y)
	c.testAndSetZero(c.Y)
}

func (c *Cpu) SetBranchCycleCount(a uint16) {
	if ((ProgramCounter - 1) & 0xFF00 >> 8) != ((a & 0xFF00) >> 8) {
		c.CycleCount = 4
	} else {
		c.CycleCount = 3
	}
}

func (c *Cpu) Bpl() {
	if !c.getNegative() {
		a := c.relativeAddress()

		c.SetBranchCycleCount(a)

		ProgramCounter = a
	} else {
		ProgramCounter++
	}
}

func (c *Cpu) Bmi() {
	if c.getNegative() {
		a := c.relativeAddress()

		c.SetBranchCycleCount(a)

		ProgramCounter = a
	} else {
		ProgramCounter++
	}
}

func (c *Cpu) Bvc() {
	if !c.getOverflow() {
		a := c.relativeAddress()

		c.SetBranchCycleCount(a)

		ProgramCounter = a
	} else {
		ProgramCounter++
	}
}

func (c *Cpu) Bvs() {
	if c.getOverflow() {
		a := c.relativeAddress()

		c.SetBranchCycleCount(a)

		ProgramCounter = a
	} else {
		ProgramCounter++
	}
}

func (c *Cpu) Bcc() {
	if !c.getCarry() {
		a := c.relativeAddress()

		c.SetBranchCycleCount(a)

		ProgramCounter = a
	} else {
		ProgramCounter++
	}
}

func (c *Cpu) Bcs() {
	if c.getCarry() {
		a := c.relativeAddress()

		c.SetBranchCycleCount(a)

		ProgramCounter = a
	} else {
		ProgramCounter++
	}
}

func (c *Cpu) Bne() {
	if !c.getZero() {
		a := c.relativeAddress()

		c.SetBranchCycleCount(a)

		ProgramCounter = a
	} else {
		ProgramCounter++
	}
}

func (c *Cpu) Beq() {
	if c.getZero() {
		a := c.relativeAddress()

		c.SetBranchCycleCount(a)

		ProgramCounter = a
	} else {
		ProgramCounter++
	}
}

func (c *Cpu) Txs() {
	c.StackPointer = c.X
}

func (c *Cpu) Tsx() {
	c.X = c.StackPointer

	c.testAndSetZero(c.X)
	c.testAndSetNegative(c.X)
}

func (c *Cpu) Pha() {
	c.pushToStack(c.A)
}

func (c *Cpu) Pla() {
	val := c.pullFromStack()

	c.A = val

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
}

func (c *Cpu) Php() {
	// BRK and PHP push P OR #$10, so that the IRQ handler can tell 
	// whether the entry was from a BRK or from an /IRQ.
	c.pushToStack(c.P | 0x10)
}

func (c *Cpu) Plp() {
	val := c.pullFromStack()

	// Unset bit 5 since it's unused in the NES
	c.P = (val | 0x30) - 0x10
}

func (c *Cpu) Compare(register Word, value Word) {
	r := register - value

	c.testAndSetZero(r)
	c.testAndSetNegative(r)
	c.testAndSetCarrySubtraction(int(register) - int(value))
}

func (c *Cpu) Cmp(location uint16) {
	val, _ := Ram.Read(location)
	c.Compare(c.A, val)
}

func (c *Cpu) Cpx(location uint16) {
	val, _ := Ram.Read(location)
	c.Compare(c.X, val)
}

func (c *Cpu) Cpy(location uint16) {
	val, _ := Ram.Read(location)
	c.Compare(c.Y, val)
}

func (c *Cpu) Sbc(location uint16) {
	val, _ := Ram.Read(location)

	cache := c.A
	c.A = cache - val

	c.A = c.A - (1 - c.P&0x01)

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
	c.testAndSetOverflowSubtraction(cache, val)
	c.testAndSetCarrySubtraction(int(cache) - int(val) - (1 - int(c.P&0x01)))

	c.A = c.A & 0xFF
}

func (c *Cpu) Clc() {
	c.clearCarry()
}

func (c *Cpu) Sec() {
	c.setCarry()
}

func (c *Cpu) Cli() {
	c.clearIrqDisable()
}

func (c *Cpu) Sei() {
	c.setIrqDisable()
}

func (c *Cpu) Clv() {
	c.clearOverflow()
}

func (c *Cpu) Cld() {
	c.clearDecimalMode()
}

func (c *Cpu) Sed() {
	c.setDecimalMode()
}

func (c *Cpu) And(location uint16) {
	val, _ := Ram.Read(location)
	c.A = c.A & val

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
}

func (c *Cpu) Ora(location uint16) {
	val, _ := Ram.Read(location)
	c.A = c.A | val
	c.A &= 0xFF

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
}

func (c *Cpu) Eor(location uint16) {
	val, _ := Ram.Read(location)
	c.A = c.A ^ val

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
}

func (c *Cpu) Dec(location uint16) {
	val, _ := Ram.Read(location)
	val = val - 1

	Ram.Write(location, val)

	c.testAndSetNegative(val)
	c.testAndSetZero(val)
}

func (c *Cpu) Inc(location uint16) {
	val, _ := Ram.Read(location)
	val = val + 1

	Ram.Write(location, val)

	c.testAndSetNegative(val)
	c.testAndSetZero(val)
}

func (c *Cpu) Brk() {
	// perfect example of the confusion the "B flag exists in status register" 
	// causes (pdq, nothing specific to you; this confusion is present in 
	// almost every 6502 book and web page). 
	//
	// As pdq said, BRK does the following: 
	// 
	// 1. Push address of BRK instruction + 2 
	// 2. PHP 
	// 3. SEI 
	// 4. JMP ($FFFE)
	ProgramCounter = ProgramCounter + 1

	c.pushToStack(Word(ProgramCounter >> 8))
	c.pushToStack(Word(ProgramCounter & 0xFF))

	c.Php()
	c.Sei()

	c.setIrqDisable()

	h, _ := Ram.Read(0xFFFF)
	l, _ := Ram.Read(0xFFFE)

	ProgramCounter = uint16(h)<<8 + uint16(l)
}

func (c *Cpu) Jsr(location uint16) {
	high := (ProgramCounter - 1) >> 8
	low := (ProgramCounter - 1) & 0xFF

	c.pushToStack(Word(high))
	c.pushToStack(Word(low))

	ProgramCounter = location
}

func (c *Cpu) Rti() {
	c.Plp()

	low := c.pullFromStack()
	high := c.pullFromStack()

	ProgramCounter = ((uint16(high) << 8) + uint16(low))
}

func (c *Cpu) Rts() {
	low := c.pullFromStack()
	high := c.pullFromStack()

	ProgramCounter = ((uint16(high) << 8) + uint16(low)) + 1
}

func (c *Cpu) Lsr(location uint16) {
	val, _ := Ram.Read(location)

	if val&0x01 > 0x00 {
		c.setCarry()
	} else {
		c.clearCarry()
	}

	Ram.Write(location, val>>1)

	val, _ = Ram.Read(location)

	c.testAndSetNegative(val)
	c.testAndSetZero(val)
}

func (c *Cpu) LsrAcc() {
	if c.A&0x01 > 0 {
		c.setCarry()
	} else {
		c.clearCarry()
	}

	c.A = c.A >> 1

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
}

func (c *Cpu) Asl(location uint16) {
	val, _ := Ram.Read(location)

	if val&0x80 > 0 {
		c.setCarry()
	} else {
		c.clearCarry()
	}

	Ram.Write(location, val<<1)

	val, _ = Ram.Read(location)
	c.testAndSetNegative(val)
	c.testAndSetZero(val)
}

func (c *Cpu) AslAcc() {
	if c.A&0x80 > 0 {
		c.setCarry()
	} else {
		c.clearCarry()
	}

	c.A = c.A << 1

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
}

func (c *Cpu) Rol(location uint16) {
	value, _ := Ram.Read(location)

	carry := value & 0x80

	value = value << 1

	if c.getCarry() {
		value += 1
	}

	if carry > 0x00 {
		c.setCarry()
	} else {
		c.clearCarry()
	}

	Ram.Write(location, value)

	value, _ = Ram.Read(location)
	c.testAndSetNegative(value)
	c.testAndSetZero(value)
}

func (c *Cpu) RolAcc() {
	carry := c.A & 0x80

	c.A = c.A << 1

	if c.getCarry() {
		c.A += 1
	}

	if carry > 0x00 {
		c.setCarry()
	} else {
		c.clearCarry()
	}

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
}

func (c *Cpu) Ror(location uint16) {
	value, _ := Ram.Read(location)

	carry := value & 0x1

	value = value >> 1

	if c.getCarry() {
		value += 0x80
	}

	if carry > 0x00 {
		c.setCarry()
	} else {
		c.clearCarry()
	}

	Ram.Write(location, value)

	value, _ = Ram.Read(location)
	c.testAndSetNegative(value)
	c.testAndSetZero(value)
}

func (c *Cpu) RorAcc() {
	carry := c.A & 0x1

	c.A = c.A >> 1

	if c.getCarry() {
		c.A += 0x80
	}

	if carry > 0x00 {
		c.setCarry()
	} else {
		c.clearCarry()
	}

	c.testAndSetNegative(c.A)
	c.testAndSetZero(c.A)
}

func (c *Cpu) Bit(location uint16) {
	val, _ := Ram.Read(location)

	if val&c.A == 0 {
		c.setZero()
	} else {
		c.clearZero()
	}

	if val&0x80 > 0x00 {
		c.setNegative()
	} else {
		c.clearNegative()
	}

	if val&0x40 > 0x00 {
		c.setOverflow()
	} else {
		c.clearOverflow()
	}
}

func (c *Cpu) PerformIrq() {
	high := ProgramCounter >> 8
	low := ProgramCounter & 0xFF

	c.pushToStack(Word(high))
	c.pushToStack(Word(low))

	c.pushToStack(c.P)

	h, _ := Ram.Read(0xFFFF)
	l, _ := Ram.Read(0xFFFE)

	ProgramCounter = uint16(h)<<8 + uint16(l)
}

func (c *Cpu) PerformNmi() {
	high := ProgramCounter >> 8
	low := ProgramCounter & 0xFF

	c.pushToStack(Word(high))
	c.pushToStack(Word(low))

	c.pushToStack(c.P)

	h, _ := Ram.Read(0xFFFB)
	l, _ := Ram.Read(0xFFFA)

	ProgramCounter = uint16(h)<<8 + uint16(l)
}

func (c *Cpu) PerformReset() {
	// $2000.7 enables/disables NMIs
	if ppu.NmiOnVblank != 0x0 {
		high, _ := Ram.Read(0xFFFD)
		low, _ := Ram.Read(0xFFFC)

		ProgramCounter = uint16(high)<<8 + uint16(low)
	}
}

func (c *Cpu) RequestInterrupt(i int) {
	c.InterruptRequested = i
}

func (c *Cpu) Init() {
	c.Reset()
	c.InterruptRequested = InterruptNone
}

func (c *Cpu) Reset() {
	c.X = 0
	c.Y = 0
	c.A = 0
	c.CycleCount = 0
	c.P = 0x34
	c.StackPointer = 0xFD

	c.Accurate = true
	c.InterruptRequested = InterruptNone
}

func (c *Cpu) Step() int {
	// Used during a DMA
	if c.CyclesToWait > 0 {
		c.CyclesToWait--
		return 0
	}

	// Check if an interrupt was requested
	switch c.InterruptRequested {
	case InterruptIrq:
		if !c.getIrqDisable() {
			c.PerformIrq()
		}
		c.InterruptRequested = InterruptNone
	case InterruptNmi:
		c.PerformNmi()
		c.InterruptRequested = InterruptNone
	case InterruptReset:
		c.PerformReset()
		c.InterruptRequested = InterruptNone
	}

	logpc := ProgramCounter
	opcode, _ := Ram.Read(ProgramCounter)

	c.Opcode = opcode

	ProgramCounter++

	if c.Verbose {
		Disassemble(opcode, c, ProgramCounter)
	}

	switch opcode {
	// ADC
	case 0x69:
		c.CycleCount = 2
		c.Adc(c.immediateAddress())
	case 0x65:
		c.CycleCount = 3
		c.Adc(c.zeroPageAddress())
	case 0x75:
		c.CycleCount = 4
		c.Adc(c.zeroPageIndexedAddress(c.X))
	case 0x6D:
		c.CycleCount = 4
		c.Adc(c.absoluteAddress())
	case 0x7D:
		c.CycleCount = 4
		c.Adc(c.absoluteIndexedAddress(c.X))
	case 0x79:
		c.CycleCount = 4
		c.Adc(c.absoluteIndexedAddress(c.Y))
	case 0x61:
		c.CycleCount = 6
		c.Adc(c.indexedIndirectAddress())
	case 0x71:
		c.CycleCount = 5
		c.Adc(c.indirectIndexedAddress())
	// LDA
	case 0xA9:
		c.CycleCount = 2
		c.Lda(c.immediateAddress())
	case 0xA5:
		c.CycleCount = 3
		c.Lda(c.zeroPageAddress())
	case 0xB5:
		c.CycleCount = 4
		c.Lda(c.zeroPageIndexedAddress(c.X))
	case 0xAD:
		c.CycleCount = 4
		c.Lda(c.absoluteAddress())
	case 0xBD:
		c.CycleCount = 4
		c.Lda(c.absoluteIndexedAddress(c.X))
	case 0xB9:
		c.CycleCount = 4
		c.Lda(c.absoluteIndexedAddress(c.Y))
	case 0xA1:
		c.CycleCount = 6
		c.Lda(c.indexedIndirectAddress())
	case 0xB1:
		c.CycleCount = 5
		c.Lda(c.indirectIndexedAddress())
	// LDX
	case 0xA2:
		c.CycleCount = 2
		c.Ldx(c.immediateAddress())
	case 0xA6:
		c.CycleCount = 3
		c.Ldx(c.zeroPageAddress())
	case 0xB6:
		c.CycleCount = 4
		c.Ldx(c.zeroPageIndexedAddress(c.Y))
	case 0xAE:
		c.CycleCount = 4
		c.Ldx(c.absoluteAddress())
	case 0xBE:
		c.CycleCount = 4
		c.Ldx(c.absoluteIndexedAddress(c.Y))
	// LDY
	case 0xA0:
		c.CycleCount = 2
		c.Ldy(c.immediateAddress())
	case 0xA4:
		c.CycleCount = 3
		c.Ldy(c.zeroPageAddress())
	case 0xB4:
		c.CycleCount = 4
		c.Ldy(c.zeroPageIndexedAddress(c.X))
	case 0xAC:
		c.CycleCount = 4
		c.Ldy(c.absoluteAddress())
	case 0xBC:
		c.CycleCount = 4
		c.Ldy(c.absoluteIndexedAddress(c.X))
	// STA
	case 0x85:
		c.CycleCount = 3
		c.Sta(c.zeroPageAddress())
	case 0x95:
		c.CycleCount = 4
		c.Sta(c.zeroPageIndexedAddress(c.X))
	case 0x8D:
		c.CycleCount = 4
		c.Sta(c.absoluteAddress())
	case 0x9D:
		c.Sta(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 5
	case 0x99:
		c.Sta(c.absoluteIndexedAddress(c.Y))
		c.CycleCount = 5
	case 0x81:
		c.CycleCount = 6
		c.Sta(c.indexedIndirectAddress())
	case 0x91:
		c.Sta(c.indirectIndexedAddress())
		c.CycleCount = 6
	// STX
	case 0x86:
		c.CycleCount = 3
		c.Stx(c.zeroPageAddress())
	case 0x96:
		c.CycleCount = 4
		c.Stx(c.zeroPageIndexedAddress(c.Y))
	case 0x8E:
		c.CycleCount = 4
		c.Stx(c.absoluteAddress())
	// STY
	case 0x84:
		c.CycleCount = 3
		c.Sty(c.zeroPageAddress())
	case 0x94:
		c.CycleCount = 4
		c.Sty(c.zeroPageIndexedAddress(c.X))
	case 0x8C:
		c.CycleCount = 4
		c.Sty(c.absoluteAddress())
	// JMP
	case 0x4C:
		c.CycleCount = 3
		c.Jmp(c.absoluteAddress())
	case 0x6C:
		c.CycleCount = 5
		c.Jmp(c.indirectAbsoluteAddress(ProgramCounter))
	// JSR
	case 0x20:
		c.CycleCount = 6
		c.Jsr(c.absoluteAddress())
	// Register Instructions
	case 0xAA:
		c.CycleCount = 2
		c.Tax()
	case 0x8A:
		c.CycleCount = 2
		c.Txa()
	case 0xCA:
		c.CycleCount = 2
		c.Dex()
	case 0xE8:
		c.CycleCount = 2
		c.Inx()
	case 0xA8:
		c.CycleCount = 2
		c.Tay()
	case 0x98:
		c.CycleCount = 2
		c.Tya()
	case 0x88:
		c.CycleCount = 2
		c.Dey()
	case 0xC8:
		c.CycleCount = 2
		c.Iny()
	// Branch Instructions
	case 0x10:
		c.CycleCount = 2
		c.Bpl()
	case 0x30:
		c.CycleCount = 2
		c.Bmi()
	case 0x50:
		c.CycleCount = 2
		c.Bvc()
	case 0x70:
		c.CycleCount = 2
		c.Bvs()
	case 0x90:
		c.CycleCount = 2
		c.Bcc()
	case 0xB0:
		c.CycleCount = 2
		c.Bcs()
	case 0xD0:
		c.CycleCount = 2
		c.Bne()
	case 0xF0:
		c.CycleCount = 2
		c.Beq()
	// CMP
	case 0xC9:
		c.CycleCount = 2
		c.Cmp(c.immediateAddress())
	case 0xC5:
		c.CycleCount = 3
		c.Cmp(c.zeroPageAddress())
	case 0xD5:
		c.CycleCount = 4
		c.Cmp(c.zeroPageIndexedAddress(c.X))
	case 0xCD:
		c.CycleCount = 4
		c.Cmp(c.absoluteAddress())
	case 0xDD:
		c.CycleCount = 4
		c.Cmp(c.absoluteIndexedAddress(c.X))
	case 0xD9:
		c.CycleCount = 4
		c.Cmp(c.absoluteIndexedAddress(c.Y))
	case 0xC1:
		c.CycleCount = 6
		c.Cmp(c.indexedIndirectAddress())
	case 0xD1:
		c.CycleCount = 5
		c.Cmp(c.indirectIndexedAddress())
	// CPX
	case 0xE0:
		c.CycleCount = 2
		c.Cpx(c.immediateAddress())
	case 0xE4:
		c.CycleCount = 3
		c.Cpx(c.zeroPageAddress())
	case 0xEC:
		c.CycleCount = 4
		c.Cpx(c.absoluteAddress())
	// CPY
	case 0xC0:
		c.CycleCount = 2
		c.Cpy(c.immediateAddress())
	case 0xC4:
		c.CycleCount = 3
		c.Cpy(c.zeroPageAddress())
	case 0xCC:
		c.CycleCount = 4
		c.Cpy(c.absoluteAddress())
	// SBC
	case 0xE9:
		c.CycleCount = 2
		c.Sbc(c.immediateAddress())
	case 0xE5:
		c.CycleCount = 3
		c.Sbc(c.zeroPageAddress())
	case 0xF5:
		c.CycleCount = 4
		c.Sbc(c.zeroPageIndexedAddress(c.X))
	case 0xED:
		c.CycleCount = 4
		c.Sbc(c.absoluteAddress())
	case 0xFD:
		c.CycleCount = 4
		c.Sbc(c.absoluteIndexedAddress(c.X))
	case 0xF9:
		c.CycleCount = 4
		c.Sbc(c.absoluteIndexedAddress(c.Y))
	case 0xE1:
		c.CycleCount = 6
		c.Sbc(c.indexedIndirectAddress())
	case 0xF1:
		c.CycleCount = 5
		c.Sbc(c.indirectIndexedAddress())
	// Flag Instructions
	case 0x18:
		c.CycleCount = 2
		c.Clc()
	case 0x38:
		c.CycleCount = 2
		c.Sec()
	case 0x58:
		c.CycleCount = 2
		c.Cli()
	case 0x78:
		c.CycleCount = 2
		c.Sei()
	case 0xB8:
		c.CycleCount = 2
		c.Clv()
	case 0xD8:
		c.CycleCount = 2
		c.Cld()
	case 0xF8:
		c.CycleCount = 2
		c.Sed()
	// Stack instructions
	case 0x9A:
		c.CycleCount = 2
		c.Txs()
	case 0xBA:
		c.CycleCount = 2
		c.Tsx()
	case 0x48:
		c.CycleCount = 3
		c.Pha()
	case 0x68:
		c.CycleCount = 4
		c.Pla()
	case 0x08:
		c.CycleCount = 3
		c.Php()
	case 0x28:
		c.CycleCount = 4
		c.Plp()
	// AND
	case 0x29:
		c.CycleCount = 2
		c.And(c.immediateAddress())
	case 0x25:
		c.CycleCount = 3
		c.And(c.zeroPageAddress())
	case 0x35:
		c.CycleCount = 4
		c.And(c.zeroPageIndexedAddress(c.X))
	case 0x2d:
		c.CycleCount = 4
		c.And(c.absoluteAddress())
	case 0x3d:
		c.CycleCount = 4
		c.And(c.absoluteIndexedAddress(c.X))
	case 0x39:
		c.CycleCount = 4
		c.And(c.absoluteIndexedAddress(c.Y))
	case 0x21:
		c.CycleCount = 6
		c.And(c.indexedIndirectAddress())
	case 0x31:
		c.CycleCount = 5
		c.And(c.indirectIndexedAddress())
	// ORA
	case 0x09:
		c.CycleCount = 2
		c.Ora(c.immediateAddress())
	case 0x05:
		c.CycleCount = 3
		c.Ora(c.zeroPageAddress())
	case 0x15:
		c.CycleCount = 4
		c.Ora(c.zeroPageIndexedAddress(c.X))
	case 0x0d:
		c.CycleCount = 4
		c.Ora(c.absoluteAddress())
	case 0x1d:
		c.CycleCount = 4
		c.Ora(c.absoluteIndexedAddress(c.X))
	case 0x19:
		c.CycleCount = 4
		c.Ora(c.absoluteIndexedAddress(c.Y))
	case 0x01:
		c.CycleCount = 6
		c.Ora(c.indexedIndirectAddress())
	case 0x11:
		c.CycleCount = 5
		c.Ora(c.indirectIndexedAddress())
	// EOR
	case 0x49:
		c.CycleCount = 2
		c.Eor(c.immediateAddress())
	case 0x45:
		c.CycleCount = 3
		c.Eor(c.zeroPageAddress())
	case 0x55:
		c.CycleCount = 4
		c.Eor(c.zeroPageIndexedAddress(c.X))
	case 0x4d:
		c.CycleCount = 4
		c.Eor(c.absoluteAddress())
	case 0x5d:
		c.CycleCount = 4
		c.Eor(c.absoluteIndexedAddress(c.X))
	case 0x59:
		c.CycleCount = 4
		c.Eor(c.absoluteIndexedAddress(c.Y))
	case 0x41:
		c.CycleCount = 6
		c.Eor(c.indexedIndirectAddress())
	case 0x51:
		c.CycleCount = 5
		c.Eor(c.indirectIndexedAddress())
	// DEC
	case 0xc6:
		c.CycleCount = 5
		c.Dec(c.zeroPageAddress())
	case 0xd6:
		c.CycleCount = 6
		c.Dec(c.zeroPageIndexedAddress(c.X))
	case 0xce:
		c.CycleCount = 6
		c.Dec(c.absoluteAddress())
	case 0xde:
		c.Dec(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 7
	// INC
	case 0xe6:
		c.CycleCount = 5
		c.Inc(c.zeroPageAddress())
	case 0xf6:
		c.CycleCount = 6
		c.Inc(c.zeroPageIndexedAddress(c.X))
	case 0xee:
		c.CycleCount = 6
		c.Inc(c.absoluteAddress())
	case 0xfe:
		c.Inc(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 7
	// BRK
	case 0x00:
		c.CycleCount = 7
		c.Brk()
	// RTI
	case 0x40:
		c.CycleCount = 6
		c.Rti()
	// RTS
	case 0x60:
		c.CycleCount = 6
		c.Rts()
	// NOP
	case 0xea:
		c.CycleCount = 2
	// LSR
	case 0x4a:
		c.CycleCount = 2
		c.LsrAcc()
	case 0x46:
		c.CycleCount = 5
		c.Lsr(c.zeroPageAddress())
	case 0x56:
		c.CycleCount = 6
		c.Lsr(c.zeroPageIndexedAddress(c.X))
	case 0x4e:
		c.CycleCount = 6
		c.Lsr(c.absoluteAddress())
	case 0x5e:
		c.Lsr(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 7
	// ASL
	case 0x0a:
		c.CycleCount = 2
		c.AslAcc()
	case 0x06:
		c.CycleCount = 5
		c.Asl(c.zeroPageAddress())
	case 0x16:
		c.CycleCount = 6
		c.Asl(c.zeroPageIndexedAddress(c.X))
	case 0x0e:
		c.CycleCount = 6
		c.Asl(c.absoluteAddress())
	case 0x1E:
		c.Asl(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 7
	// ROL
	case 0x2a:
		c.CycleCount = 2
		c.RolAcc()
	case 0x26:
		c.CycleCount = 5
		c.Rol(c.zeroPageAddress())
	case 0x36:
		c.CycleCount = 6
		c.Rol(c.zeroPageIndexedAddress(c.X))
	case 0x2e:
		c.CycleCount = 6
		c.Rol(c.absoluteAddress())
	case 0x3e:
		c.Rol(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 7
	// ROR
	case 0x6a:
		c.CycleCount = 2
		c.RorAcc()
	case 0x66:
		c.CycleCount = 5
		c.Ror(c.zeroPageAddress())
	case 0x76:
		c.CycleCount = 6
		c.Ror(c.zeroPageIndexedAddress(c.X))
	case 0x6e:
		c.CycleCount = 6
		c.Ror(c.absoluteAddress())
	case 0x7e:
		c.Ror(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 7
	// BIT
	case 0x24:
		c.CycleCount = 3
		c.Bit(c.zeroPageAddress())
	case 0x2c:
		c.CycleCount = 4
		c.Bit(c.absoluteAddress())
	default:
		log.Fatalf("Invalid opcode at 0x%X: 0x%X", logpc, opcode)
	}

	c.Timestamp = (c.CycleCount * 15)

	return c.CycleCount
}
