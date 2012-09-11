package main

const (
	InterruptNone = iota
	InterruptIrq
	InterruptReset
	InterruptNmi
)

var (
	ProgramCounter = 0x8000
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
}

func (cpu *Cpu) getCarry() bool {
	return cpu.P&0x01 == 0x01
}

func (cpu *Cpu) getZero() bool {
	return cpu.P&0x02 == 0x02
}

func (cpu *Cpu) getIrqDisable() bool {
	return cpu.P&0x04 == 0x04
}

func (cpu *Cpu) getDecimalMode() bool {
	return cpu.P&0x08 == 0x08
}

func (cpu *Cpu) getBrkCommand() bool {
	return cpu.P&0x10 == 0x10
}

func (cpu *Cpu) getOverflow() bool {
	return cpu.P&0x40 == 0x40
}

func (cpu *Cpu) getNegative() bool {
	return cpu.P&0x80 == 0x80
}

func (cpu *Cpu) setCarry() {
	cpu.P = cpu.P | 0x01
}

func (cpu *Cpu) setZero() {
	cpu.P = cpu.P | 0x02
}

func (cpu *Cpu) setIrqDisable() {
	cpu.P = cpu.P | 0x04
}

func (cpu *Cpu) setDecimalMode() {
	cpu.P = cpu.P | 0x08
}

func (cpu *Cpu) setBrkCommand() {
	cpu.P = cpu.P | 0x10
}

func (cpu *Cpu) setOverflow() {
	cpu.P = cpu.P | 0x40
}

func (cpu *Cpu) setNegative() {
	cpu.P = cpu.P | 0x80
}

func (cpu *Cpu) clearCarry() {
	cpu.P = cpu.P & 0xFE
}

func (cpu *Cpu) clearZero() {
	cpu.P = cpu.P & 0xFD
}

func (cpu *Cpu) clearIrqDisable() {
	cpu.P = cpu.P & 0xFB
}

func (cpu *Cpu) clearDecimalMode() {
	cpu.P = cpu.P & 0xF7
}

func (cpu *Cpu) clearBrkCommand() {
	cpu.P = cpu.P & 0xEF
}

func (cpu *Cpu) clearOverflow() {
	cpu.P = cpu.P & 0xBF
}

func (cpu *Cpu) clearNegative() {
	cpu.P = cpu.P & 0x7F
}

func (cpu *Cpu) pushToStack(value Word) {
	Ram.Write(0x100+int(cpu.StackPointer), value)
	cpu.StackPointer--
}

func (cpu *Cpu) pullFromStack() Word {
	cpu.StackPointer++
	val, _ := Ram.Read(0x100 + int(cpu.StackPointer))

	return val
}

func (cpu *Cpu) testAndSetNegative(value Word) {
	if value&0x80 == 0x80 {
		cpu.setNegative()
		return
	}

	cpu.clearNegative()
}

func (cpu *Cpu) testAndSetZero(value Word) {
	if value == 0x00 {
		cpu.setZero()
		return
	}

	cpu.clearZero()
}

func (cpu *Cpu) testAndSetCarryAddition(result int) {
	if result > 0xff {
		cpu.setCarry()
		return
	}

	cpu.clearCarry()
}

func (cpu *Cpu) testAndSetCarrySubtraction(result int) {
	if result < 0x00 {
		cpu.clearCarry()
		return
	}

	cpu.setCarry()
}

func (cpu *Cpu) testAndSetOverflowAddition(a Word, b Word) {
	if (a & 0x80) == (b & 0x80) {
		switch {
		case int(a+b) > 127:
			fallthrough
		case int(a+b) < -128:
			cpu.setOverflow()
			return
		}
	}

	cpu.clearOverflow()
}

func (cpu *Cpu) testAndSetOverflowSubtraction(a Word, b Word) {
	val := a - b - (1 - cpu.P&0x01)
	if ((a^val)&0x80) != 0 && ((a^b)&0x80) != 0 {
		cpu.setOverflow()
	} else {
		cpu.clearOverflow()
	}
}

func (cpu *Cpu) immediateAddress() int {
	ProgramCounter++
	return ProgramCounter - 1
}

func (cpu *Cpu) absoluteAddress() (result int) {
	// Switch to an int (or more appropriately uint16) since we 
	// will overflow when shifting the high byte
	high, _ := Ram.Read(ProgramCounter + 1)
	low, _ := Ram.Read(ProgramCounter)

	ProgramCounter += 2
	return (int(high) << 8) + int(low)
}

func (cpu *Cpu) zeroPageAddress() int {
	ProgramCounter++
	res, _ := Ram.Read(ProgramCounter - 1)

	return int(res)
}

func (cpu *Cpu) indirectAbsoluteAddress(addr int) (result int) {
	high, _ := Ram.Read(addr + 1)
	low, _ := Ram.Read(addr)

	// Indirect jump is bugged on the 6502, it doesn't add 1 to 
	// the full 16-bit value when it reads the second byte, it 
	// adds 1 to the low byte only. So JMP (03FF) reads from 3FF 
	// and 300, not 3FF and 400.
	laddr := (int(high) << 8) + int(low)
	haddr := (int(high) << 8) + ((int(low) + 1) & 0xFF)

	ih, _ := Ram.Read(haddr)
	il, _ := Ram.Read(laddr)

	result = (int(ih) << 8) + int(il)
	return
}

func (cpu *Cpu) absoluteIndexedAddress(index Word) (result int) {
	// Switch to an int (or more appropriately uint16) since we 
	// will overflow when shifting the high byte
	high, _ := Ram.Read(ProgramCounter + 1)
	low, _ := Ram.Read(ProgramCounter)

	address := (int(high) << 8) + int(low) + int(index)

	if address > 0xFFFF {
		address = address & 0xFFFF
	}

	ProgramCounter += 2
	return address
}

func (cpu *Cpu) zeroPageIndexedAddress(index Word) int {
	location, _ := Ram.Read(ProgramCounter)
	ProgramCounter++
	return int(location + index)
}

func (cpu *Cpu) indexedIndirectAddress() int {
	location, _ := Ram.Read(ProgramCounter)
	location = location + cpu.X

	ProgramCounter++

	// Switch to an int (or more appropriately uint16) since we 
	// will overflow when shifting the high byte
	high, _ := Ram.Read(location + 1)
	low, _ := Ram.Read(location)

	return (int(high) << 8) + int(low)
}

func (cpu *Cpu) indirectIndexedAddress() int {
	location, _ := Ram.Read(ProgramCounter)

	// Switch to an int (or more appropriately uint16) since we 
	// will overflow when shifting the high byte
	high, _ := Ram.Read(location + 1)
	low, _ := Ram.Read(location)

	address := (int(high) << 8) + int(low) + int(cpu.Y)

	if address > 0xFFFF {
		address = address & 0xFFFF
	}

	ProgramCounter++
	return address
}

func (cpu *Cpu) relativeAddress() (a int) {
	val, _ := Ram.Read(ProgramCounter)

	a = int(val)
	if a < 0x80 {
		a = a + ProgramCounter
	} else {
		a = a + (ProgramCounter - 0x100)
	}

	a++

	return
}

func (cpu *Cpu) accumulatorAddress() int {
	return 0
}

func (cpu *Cpu) Adc(location int) {
	val, _ := Ram.Read(location)

	cached := cpu.A

	cpu.A = cpu.A + val

	if cpu.getCarry() {
		cpu.A++
	}

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
	cpu.testAndSetOverflowAddition(cached, val)
	cpu.testAndSetCarryAddition(int(cached) + int(val) + int(cpu.P&0x01))
}

func (cpu *Cpu) Lda(location int) {
	val, _ := Ram.Read(location)
	cpu.A = val

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Ldx(location int) {
	val, _ := Ram.Read(location)
	cpu.X = val

	cpu.testAndSetNegative(cpu.X)
	cpu.testAndSetZero(cpu.X)
}

func (cpu *Cpu) Ldy(location int) {
	val, _ := Ram.Read(location)
	cpu.Y = val

	cpu.testAndSetNegative(cpu.Y)
	cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu) Sta(location int) {
	Ram.Write(location, cpu.A)
}

func (cpu *Cpu) Stx(location int) {
	Ram.Write(location, cpu.X)
}

func (cpu *Cpu) Sty(location int) {
	Ram.Write(location, cpu.Y)
}

func (cpu *Cpu) Jmp(location int) {
	ProgramCounter = location
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
		cpu.setZero()
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
		cpu.setZero()
	}

	cpu.testAndSetNegative(cpu.Y)
	cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu) Iny() {
	cpu.Y = cpu.Y + 1

	cpu.testAndSetNegative(cpu.Y)
	cpu.testAndSetZero(cpu.Y)
}

func (cpu *Cpu) Bpl() {
	if !cpu.getNegative() {
		ProgramCounter = cpu.relativeAddress()
	} else {
		ProgramCounter++
	}
}

func (cpu *Cpu) Bmi() {
	if cpu.getNegative() {
		ProgramCounter = cpu.relativeAddress()
	} else {
		ProgramCounter++
	}
}

func (cpu *Cpu) Bvc() {
	if !cpu.getOverflow() {
		ProgramCounter = cpu.relativeAddress()
	} else {
		ProgramCounter++
	}
}

func (cpu *Cpu) Bvs() {
	if cpu.getOverflow() {
		ProgramCounter = cpu.relativeAddress()
	} else {
		ProgramCounter++
	}
}

func (cpu *Cpu) Bcc() {
	if !cpu.getCarry() {
		ProgramCounter = cpu.relativeAddress()
	} else {
		ProgramCounter++
	}
}

func (cpu *Cpu) Bcs() {
	if cpu.getCarry() {
		ProgramCounter = cpu.relativeAddress()
	} else {
		ProgramCounter++
	}
}

func (cpu *Cpu) Bne() {
	if !cpu.getZero() {
		ProgramCounter = cpu.relativeAddress()
	} else {
		ProgramCounter++
	}
}

func (cpu *Cpu) Beq() {
	if cpu.getZero() {
		ProgramCounter = cpu.relativeAddress()
	} else {
		ProgramCounter++
	}
}

func (cpu *Cpu) Txs() {
	cpu.StackPointer = cpu.X
}

func (cpu *Cpu) Tsx() {
	cpu.X = cpu.StackPointer

	cpu.testAndSetZero(cpu.X)
	cpu.testAndSetNegative(cpu.X)
}

func (cpu *Cpu) Pha() {
	cpu.pushToStack(cpu.A)
}

func (cpu *Cpu) Pla() {
	val := cpu.pullFromStack()

	cpu.A = val

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Php() {
	// BRK and PHP push P OR #$10, so that the IRQ handler can tell 
	// whether the entry was from a BRK or from an /IRQ.
	cpu.pushToStack(cpu.P | 0x10)
}

func (cpu *Cpu) Plp() {
	val := cpu.pullFromStack()

	// Unset bit 5 since it's unused in the NES
	cpu.P = (val | 0x30) - 0x10
}

func (cpu *Cpu) Compare(register Word, value Word) {
	r := register - value

	cpu.testAndSetZero(r)
	cpu.testAndSetNegative(r)
	cpu.testAndSetCarrySubtraction(int(register) - int(value))
}

func (cpu *Cpu) Cmp(location int) {
	val, _ := Ram.Read(location)
	cpu.Compare(cpu.A, val)
}

func (cpu *Cpu) Cpx(location int) {
	val, _ := Ram.Read(location)
	cpu.Compare(cpu.X, val)
}

func (cpu *Cpu) Cpy(location int) {
	val, _ := Ram.Read(location)
	cpu.Compare(cpu.Y, val)
}

func (cpu *Cpu) Sbc(location int) {
	val, _ := Ram.Read(location)

	cache := cpu.A
	cpu.A = cache - val

	cpu.A = cpu.A - (1 - cpu.P&0x01)

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
	cpu.testAndSetOverflowSubtraction(cache, val)
	cpu.testAndSetCarrySubtraction(int(cache) - int(val) - (1 - int(cpu.P&0x01)))

	cpu.A = cpu.A & 0xff
}

func (cpu *Cpu) Clc() {
	cpu.clearCarry()
}

func (cpu *Cpu) Sec() {
	cpu.setCarry()
}

func (cpu *Cpu) Cli() {
	cpu.clearIrqDisable()
}

func (cpu *Cpu) Sei() {
	cpu.setIrqDisable()
}

func (cpu *Cpu) Clv() {
	cpu.clearOverflow()
}

func (cpu *Cpu) Cld() {
	cpu.clearDecimalMode()
}

func (cpu *Cpu) Sed() {
	cpu.setDecimalMode()
}

func (cpu *Cpu) And(location int) {
	val, _ := Ram.Read(location)
	cpu.A = cpu.A & val

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Ora(location int) {
	val, _ := Ram.Read(location)
	cpu.A = cpu.A | val

	cpu.testAndSetNegative(val)
	cpu.testAndSetZero(val)
}

func (cpu *Cpu) Eor(location int) {
	val, _ := Ram.Read(location)
	cpu.A = cpu.A ^ val

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Dec(location int) {
	val, _ := Ram.Read(location)
	val = val - 1

	Ram.Write(location, val)

	cpu.testAndSetNegative(val)
	cpu.testAndSetZero(val)
}

func (cpu *Cpu) Inc(location int) {
	val, _ := Ram.Read(location)
	val = val + 1

	Ram.Write(location, val)

	cpu.testAndSetNegative(val)
	cpu.testAndSetZero(val)
}

func (cpu *Cpu) Brk() {
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

	cpu.pushToStack(Word(ProgramCounter >> 8))
	cpu.pushToStack(Word(ProgramCounter & 0xFF))

	cpu.Php()
	cpu.Sei()

	cpu.setIrqDisable()

	cpu.Jmp(cpu.indirectAbsoluteAddress(0xFFFE))
}

func (cpu *Cpu) Jsr(location int) {
	high := (ProgramCounter - 1) >> 8
	low := (ProgramCounter - 1) & 0xFF

	cpu.pushToStack(Word(high))
	cpu.pushToStack(Word(low))

	ProgramCounter = location
}

func (cpu *Cpu) Rti() {
	cpu.Plp()

	low := cpu.pullFromStack()
	high := cpu.pullFromStack()

	ProgramCounter = ((int(high) << 8) + int(low))
}

func (cpu *Cpu) Rts() {
	low := cpu.pullFromStack()
	high := cpu.pullFromStack()

	ProgramCounter = ((int(high) << 8) + int(low)) + 1
}

func (cpu *Cpu) Lsr(location int) {
	val, _ := Ram.Read(location)

	if val&0x01 > 0x00 {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}

	Ram.Write(location, val>>1)

	val, _ = Ram.Read(location)

	cpu.testAndSetNegative(val)
	cpu.testAndSetZero(val)
}

func (cpu *Cpu) LsrAcc() {
	if cpu.A&0x01 > 0 {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}

	cpu.A = cpu.A >> 1

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Asl(location int) {
	val, _ := Ram.Read(location)

	if val&0x80 > 0 {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}

	Ram.Write(location, val<<1)

	val, _ = Ram.Read(location)
	cpu.testAndSetNegative(val)
	cpu.testAndSetZero(val)
}

func (cpu *Cpu) AslAcc() {
	if cpu.A&0x80 > 0 {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}

	cpu.A = cpu.A << 1

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Rol(location int) {
	value, _ := Ram.Read(location)

	carry := value & 0x80

	value = value << 1

	if cpu.getCarry() {
		value += 1
	}

	if carry > 0x00 {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}

	Ram.Write(location, value)

	value, _ = Ram.Read(location)
	cpu.testAndSetNegative(value)
	cpu.testAndSetZero(value)
}

func (cpu *Cpu) RolAcc() {
	carry := cpu.A & 0x80

	cpu.A = cpu.A << 1

	if cpu.getCarry() {
		cpu.A += 1
	}

	if carry > 0x00 {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Ror(location int) {
	value, _ := Ram.Read(location)

	carry := value & 0x1

	value = value >> 1

	if cpu.getCarry() {
		value += 0x80
	}

	if carry > 0x00 {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}

	Ram.Write(location, value)

	value, _ = Ram.Read(location)
	cpu.testAndSetNegative(value)
	cpu.testAndSetZero(value)
}

func (cpu *Cpu) RorAcc() {
	carry := cpu.A & 0x1

	cpu.A = cpu.A >> 1

	if cpu.getCarry() {
		cpu.A += 0x80
	}

	if carry > 0x00 {
		cpu.setCarry()
	} else {
		cpu.clearCarry()
	}

	cpu.testAndSetNegative(cpu.A)
	cpu.testAndSetZero(cpu.A)
}

func (cpu *Cpu) Bit(location int) {
	val, _ := Ram.Read(location)

	if val&cpu.A == 0 {
		cpu.setZero()
	} else {
		cpu.clearZero()
	}

	if val&0x80 > 0x00 {
		cpu.setNegative()
	} else {
		cpu.clearNegative()
	}

	if val&0x40 > 0x00 {
		cpu.setOverflow()
	} else {
		cpu.clearOverflow()
	}
}

func (cpu *Cpu) PerformNmi() {
	// $2000.7 enables/disables NMIs
	if ppu.NmiOnVblank != 0x0 {
		high := ProgramCounter >> 8
		low := ProgramCounter & 0xFF

		cpu.pushToStack(Word(high))
		cpu.pushToStack(Word(low))

		cpu.pushToStack(cpu.P)

		h, _ := Ram.Read(0xFFFB)
		l, _ := Ram.Read(0xFFFA)

		ProgramCounter = int(h)<<8 + int(l)
	}
}

func (cpu *Cpu) PerformReset() {
	// $2000.7 enables/disables NMIs
	if ppu.NmiOnVblank != 0x0 {
		high, _ := Ram.Read(0xFFFD)
		low, _ := Ram.Read(0xFFFC)

		ProgramCounter = int(high)<<8 + int(low)
	}
}

func (cpu *Cpu) RequestInterrupt(i int) {
	cpu.InterruptRequested = i
}

func (cpu *Cpu) Init() {
	cpu.Reset()
	cpu.InterruptRequested = InterruptNone
}

func (cpu *Cpu) Reset() {
	cpu.X = 0
	cpu.Y = 0
	cpu.A = 0
	cpu.CycleCount = 0
	cpu.P = 0x34
	cpu.StackPointer = 0xFD

	cpu.Accurate = true
	cpu.InterruptRequested = InterruptNone
}

func (cpu *Cpu) Step() int {
	// Used during a DMA
	if cpu.CyclesToWait > 0 {
		cpu.CyclesToWait--
		return 0
	}

	// Check if an interrupt was requested
	switch cpu.InterruptRequested {
	case InterruptNmi:
		cpu.PerformNmi()
		cpu.InterruptRequested = InterruptNone
	case InterruptReset:
		cpu.PerformReset()
		cpu.InterruptRequested = InterruptNone
	}

	opcode, _ := Ram.Read(ProgramCounter)

	cpu.Opcode = opcode

	ProgramCounter++

	if cpu.Verbose {
		Disassemble(opcode, cpu, ProgramCounter)
	}

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
		cpu.Jmp(cpu.indirectAbsoluteAddress(ProgramCounter))
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
	case 0xB8:
		cpu.CycleCount = 2
		cpu.Clv()
	case 0xD8:
		cpu.CycleCount = 2
		cpu.Cld()
	case 0xF8:
		cpu.CycleCount = 2
		cpu.Sed()
	// Stack instructions
	case 0x9A:
		cpu.CycleCount = 2
		cpu.Txs()
	case 0xBA:
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
	default:
		panic("Invalid opcode")
	}

	return cpu.CycleCount
}
