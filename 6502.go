package main

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
	InstrOpcodes [0xFF]func()

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

	c.InstrInit()
}

func (c *Cpu) Step() int {
	// Used during a DMA
	if c.CyclesToWait > 0 {
		c.CyclesToWait--
		return 1
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

	opcode, _ := Ram.Read(ProgramCounter)

	c.Opcode = opcode

	ProgramCounter++

	if c.Verbose {
		Disassemble(opcode, c, ProgramCounter)
	}

	c.InstrOpcodes[opcode]()
	c.Timestamp = (c.CycleCount * 15)

	return c.CycleCount
}
