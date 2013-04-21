package main

func (c *Cpu) InstrInit() {
	c.InstrOpcodes[0x69] = func() {
		c.CycleCount = 2
		c.Adc(c.immediateAddress())
	}
	c.InstrOpcodes[0x65] = func() {
		c.CycleCount = 3
		c.Adc(c.zeroPageAddress())
	}
	c.InstrOpcodes[0x75] = func() {
		c.CycleCount = 4
		c.Adc(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x6D] = func() {
		c.CycleCount = 4
		c.Adc(c.absoluteAddress())
	}
	c.InstrOpcodes[0x7D] = func() {
		c.CycleCount = 4
		c.Adc(c.absoluteIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x79] = func() {
		c.CycleCount = 4
		c.Adc(c.absoluteIndexedAddress(c.Y))
	}
	c.InstrOpcodes[0x61] = func() {
		c.CycleCount = 6
		c.Adc(c.indexedIndirectAddress())
	}
	c.InstrOpcodes[0x71] = func() {
		c.CycleCount = 5
		c.Adc(c.indirectIndexedAddress())
	}
	// LDA
	c.InstrOpcodes[0xA9] = func() {
		c.CycleCount = 2
		c.Lda(c.immediateAddress())
	}
	c.InstrOpcodes[0xA5] = func() {
		c.CycleCount = 3
		c.Lda(c.zeroPageAddress())
	}
	c.InstrOpcodes[0xB5] = func() {
		c.CycleCount = 4
		c.Lda(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0xAD] = func() {
		c.CycleCount = 4
		c.Lda(c.absoluteAddress())
	}
	c.InstrOpcodes[0xBD] = func() {
		c.CycleCount = 4
		c.Lda(c.absoluteIndexedAddress(c.X))
	}
	c.InstrOpcodes[0xB9] = func() {
		c.CycleCount = 4
		c.Lda(c.absoluteIndexedAddress(c.Y))
	}
	c.InstrOpcodes[0xA1] = func() {
		c.CycleCount = 6
		c.Lda(c.indexedIndirectAddress())
	}
	c.InstrOpcodes[0xB1] = func() {
		c.CycleCount = 5
		c.Lda(c.indirectIndexedAddress())
	}
	// LDX
	c.InstrOpcodes[0xA2] = func() {
		c.CycleCount = 2
		c.Ldx(c.immediateAddress())
	}
	c.InstrOpcodes[0xA6] = func() {
		c.CycleCount = 3
		c.Ldx(c.zeroPageAddress())
	}
	c.InstrOpcodes[0xB6] = func() {
		c.CycleCount = 4
		c.Ldx(c.zeroPageIndexedAddress(c.Y))
	}
	c.InstrOpcodes[0xAE] = func() {
		c.CycleCount = 4
		c.Ldx(c.absoluteAddress())
	}
	c.InstrOpcodes[0xBE] = func() {
		c.CycleCount = 4
		c.Ldx(c.absoluteIndexedAddress(c.Y))
	}
	// LDY
	c.InstrOpcodes[0xA0] = func() {
		c.CycleCount = 2
		c.Ldy(c.immediateAddress())
	}
	c.InstrOpcodes[0xA4] = func() {
		c.CycleCount = 3
		c.Ldy(c.zeroPageAddress())
	}
	c.InstrOpcodes[0xB4] = func() {
		c.CycleCount = 4
		c.Ldy(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0xAC] = func() {
		c.CycleCount = 4
		c.Ldy(c.absoluteAddress())
	}
	c.InstrOpcodes[0xBC] = func() {
		c.CycleCount = 4
		c.Ldy(c.absoluteIndexedAddress(c.X))
	}
	// STA
	c.InstrOpcodes[0x85] = func() {
		c.CycleCount = 3
		c.Sta(c.zeroPageAddress())
	}
	c.InstrOpcodes[0x95] = func() {
		c.CycleCount = 4
		c.Sta(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x8D] = func() {
		c.CycleCount = 4
		c.Sta(c.absoluteAddress())
	}
	c.InstrOpcodes[0x9D] = func() {
		c.Sta(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 5
	}
	c.InstrOpcodes[0x99] = func() {
		c.Sta(c.absoluteIndexedAddress(c.Y))
		c.CycleCount = 5
	}
	c.InstrOpcodes[0x81] = func() {
		c.CycleCount = 6
		c.Sta(c.indexedIndirectAddress())
	}
	c.InstrOpcodes[0x91] = func() {
		c.Sta(c.indirectIndexedAddress())
		c.CycleCount = 6
	}
	// STX
	c.InstrOpcodes[0x86] = func() {
		c.CycleCount = 3
		c.Stx(c.zeroPageAddress())
	}
	c.InstrOpcodes[0x96] = func() {
		c.CycleCount = 4
		c.Stx(c.zeroPageIndexedAddress(c.Y))
	}
	c.InstrOpcodes[0x8E] = func() {
		c.CycleCount = 4
		c.Stx(c.absoluteAddress())
	}
	// STY
	c.InstrOpcodes[0x84] = func() {
		c.CycleCount = 3
		c.Sty(c.zeroPageAddress())
	}
	c.InstrOpcodes[0x94] = func() {
		c.CycleCount = 4
		c.Sty(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x8C] = func() {
		c.CycleCount = 4
		c.Sty(c.absoluteAddress())
	}
	// JMP
	c.InstrOpcodes[0x4C] = func() {
		c.CycleCount = 3
		c.Jmp(c.absoluteAddress())
	}
	c.InstrOpcodes[0x6C] = func() {
		c.CycleCount = 5
		c.Jmp(c.indirectAbsoluteAddress(ProgramCounter))
	}
	// JSR
	c.InstrOpcodes[0x20] = func() {
		c.CycleCount = 6
		c.Jsr(c.absoluteAddress())
	}
	// Register Instructions
	c.InstrOpcodes[0xAA] = func() {
		c.CycleCount = 2
		c.Tax()
	}
	c.InstrOpcodes[0x8A] = func() {
		c.CycleCount = 2
		c.Txa()
	}
	c.InstrOpcodes[0xCA] = func() {
		c.CycleCount = 2
		c.Dex()
	}
	c.InstrOpcodes[0xE8] = func() {
		c.CycleCount = 2
		c.Inx()
	}
	c.InstrOpcodes[0xA8] = func() {
		c.CycleCount = 2
		c.Tay()
	}
	c.InstrOpcodes[0x98] = func() {
		c.CycleCount = 2
		c.Tya()
	}
	c.InstrOpcodes[0x88] = func() {
		c.CycleCount = 2
		c.Dey()
	}
	c.InstrOpcodes[0xC8] = func() {
		c.CycleCount = 2
		c.Iny()
	}
	// Branch Instructions
	c.InstrOpcodes[0x10] = func() {
		c.CycleCount = 2
		c.Bpl()
	}
	c.InstrOpcodes[0x30] = func() {
		c.CycleCount = 2
		c.Bmi()
	}
	c.InstrOpcodes[0x50] = func() {
		c.CycleCount = 2
		c.Bvc()
	}
	c.InstrOpcodes[0x70] = func() {
		c.CycleCount = 2
		c.Bvs()
	}
	c.InstrOpcodes[0x90] = func() {
		c.CycleCount = 2
		c.Bcc()
	}
	c.InstrOpcodes[0xB0] = func() {
		c.CycleCount = 2
		c.Bcs()
	}
	c.InstrOpcodes[0xD0] = func() {
		c.CycleCount = 2
		c.Bne()
	}
	c.InstrOpcodes[0xF0] = func() {
		c.CycleCount = 2
		c.Beq()
	}
	// CMP
	c.InstrOpcodes[0xC9] = func() {
		c.CycleCount = 2
		c.Cmp(c.immediateAddress())
	}
	c.InstrOpcodes[0xC5] = func() {
		c.CycleCount = 3
		c.Cmp(c.zeroPageAddress())
	}
	c.InstrOpcodes[0xD5] = func() {
		c.CycleCount = 4
		c.Cmp(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0xCD] = func() {
		c.CycleCount = 4
		c.Cmp(c.absoluteAddress())
	}
	c.InstrOpcodes[0xDD] = func() {
		c.CycleCount = 4
		c.Cmp(c.absoluteIndexedAddress(c.X))
	}
	c.InstrOpcodes[0xD9] = func() {
		c.CycleCount = 4
		c.Cmp(c.absoluteIndexedAddress(c.Y))
	}
	c.InstrOpcodes[0xC1] = func() {
		c.CycleCount = 6
		c.Cmp(c.indexedIndirectAddress())
	}
	c.InstrOpcodes[0xD1] = func() {
		c.CycleCount = 5
		c.Cmp(c.indirectIndexedAddress())
	}
	// CPX
	c.InstrOpcodes[0xE0] = func() {
		c.CycleCount = 2
		c.Cpx(c.immediateAddress())
	}
	c.InstrOpcodes[0xE4] = func() {
		c.CycleCount = 3
		c.Cpx(c.zeroPageAddress())
	}
	c.InstrOpcodes[0xEC] = func() {
		c.CycleCount = 4
		c.Cpx(c.absoluteAddress())
	}
	// CPY
	c.InstrOpcodes[0xC0] = func() {
		c.CycleCount = 2
		c.Cpy(c.immediateAddress())
	}
	c.InstrOpcodes[0xC4] = func() {
		c.CycleCount = 3
		c.Cpy(c.zeroPageAddress())
	}
	c.InstrOpcodes[0xCC] = func() {
		c.CycleCount = 4
		c.Cpy(c.absoluteAddress())
	}
	// SBC
	c.InstrOpcodes[0xE9] = func() {
		c.CycleCount = 2
		c.Sbc(c.immediateAddress())
	}
	c.InstrOpcodes[0xE5] = func() {
		c.CycleCount = 3
		c.Sbc(c.zeroPageAddress())
	}
	c.InstrOpcodes[0xF5] = func() {
		c.CycleCount = 4
		c.Sbc(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0xED] = func() {
		c.CycleCount = 4
		c.Sbc(c.absoluteAddress())
	}
	c.InstrOpcodes[0xFD] = func() {
		c.CycleCount = 4
		c.Sbc(c.absoluteIndexedAddress(c.X))
	}
	c.InstrOpcodes[0xF9] = func() {
		c.CycleCount = 4
		c.Sbc(c.absoluteIndexedAddress(c.Y))
	}
	c.InstrOpcodes[0xE1] = func() {
		c.CycleCount = 6
		c.Sbc(c.indexedIndirectAddress())
	}
	c.InstrOpcodes[0xF1] = func() {
		c.CycleCount = 5
		c.Sbc(c.indirectIndexedAddress())
	}
	// Flag Instructions
	c.InstrOpcodes[0x18] = func() {
		c.CycleCount = 2
		c.Clc()
	}
	c.InstrOpcodes[0x38] = func() {
		c.CycleCount = 2
		c.Sec()
	}
	c.InstrOpcodes[0x58] = func() {
		c.CycleCount = 2
		c.Cli()
	}
	c.InstrOpcodes[0x78] = func() {
		c.CycleCount = 2
		c.Sei()
	}
	c.InstrOpcodes[0xB8] = func() {
		c.CycleCount = 2
		c.Clv()
	}
	c.InstrOpcodes[0xD8] = func() {
		c.CycleCount = 2
		c.Cld()
	}
	c.InstrOpcodes[0xF8] = func() {
		c.CycleCount = 2
		c.Sed()
	}
	// Stack instructions
	c.InstrOpcodes[0x9A] = func() {
		c.CycleCount = 2
		c.Txs()
	}
	c.InstrOpcodes[0xBA] = func() {
		c.CycleCount = 2
		c.Tsx()
	}
	c.InstrOpcodes[0x48] = func() {
		c.CycleCount = 3
		c.Pha()
	}
	c.InstrOpcodes[0x68] = func() {
		c.CycleCount = 4
		c.Pla()
	}
	c.InstrOpcodes[0x08] = func() {
		c.CycleCount = 3
		c.Php()
	}
	c.InstrOpcodes[0x28] = func() {
		c.CycleCount = 4
		c.Plp()
	}
	// AND
	c.InstrOpcodes[0x29] = func() {
		c.CycleCount = 2
		c.And(c.immediateAddress())
	}
	c.InstrOpcodes[0x25] = func() {
		c.CycleCount = 3
		c.And(c.zeroPageAddress())
	}
	c.InstrOpcodes[0x35] = func() {
		c.CycleCount = 4
		c.And(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x2d] = func() {
		c.CycleCount = 4
		c.And(c.absoluteAddress())
	}
	c.InstrOpcodes[0x3d] = func() {
		c.CycleCount = 4
		c.And(c.absoluteIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x39] = func() {
		c.CycleCount = 4
		c.And(c.absoluteIndexedAddress(c.Y))
	}
	c.InstrOpcodes[0x21] = func() {
		c.CycleCount = 6
		c.And(c.indexedIndirectAddress())
	}
	c.InstrOpcodes[0x31] = func() {
		c.CycleCount = 5
		c.And(c.indirectIndexedAddress())
	}
	// ORA
	c.InstrOpcodes[0x09] = func() {
		c.CycleCount = 2
		c.Ora(c.immediateAddress())
	}
	c.InstrOpcodes[0x05] = func() {
		c.CycleCount = 3
		c.Ora(c.zeroPageAddress())
	}
	c.InstrOpcodes[0x15] = func() {
		c.CycleCount = 4
		c.Ora(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x0d] = func() {
		c.CycleCount = 4
		c.Ora(c.absoluteAddress())
	}
	c.InstrOpcodes[0x1d] = func() {
		c.CycleCount = 4
		c.Ora(c.absoluteIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x19] = func() {
		c.CycleCount = 4
		c.Ora(c.absoluteIndexedAddress(c.Y))
	}
	c.InstrOpcodes[0x01] = func() {
		c.CycleCount = 6
		c.Ora(c.indexedIndirectAddress())
	}
	c.InstrOpcodes[0x11] = func() {
		c.CycleCount = 5
		c.Ora(c.indirectIndexedAddress())
	}
	// EOR
	c.InstrOpcodes[0x49] = func() {
		c.CycleCount = 2
		c.Eor(c.immediateAddress())
	}
	c.InstrOpcodes[0x45] = func() {
		c.CycleCount = 3
		c.Eor(c.zeroPageAddress())
	}
	c.InstrOpcodes[0x55] = func() {
		c.CycleCount = 4
		c.Eor(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x4d] = func() {
		c.CycleCount = 4
		c.Eor(c.absoluteAddress())
	}
	c.InstrOpcodes[0x5d] = func() {
		c.CycleCount = 4
		c.Eor(c.absoluteIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x59] = func() {
		c.CycleCount = 4
		c.Eor(c.absoluteIndexedAddress(c.Y))
	}
	c.InstrOpcodes[0x41] = func() {
		c.CycleCount = 6
		c.Eor(c.indexedIndirectAddress())
	}
	c.InstrOpcodes[0x51] = func() {
		c.CycleCount = 5
		c.Eor(c.indirectIndexedAddress())
	}
	// DEC
	c.InstrOpcodes[0xc6] = func() {
		c.CycleCount = 5
		c.Dec(c.zeroPageAddress())
	}
	c.InstrOpcodes[0xd6] = func() {
		c.CycleCount = 6
		c.Dec(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0xce] = func() {
		c.CycleCount = 6
		c.Dec(c.absoluteAddress())
	}
	c.InstrOpcodes[0xde] = func() {
		c.Dec(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 7
	}
	// INC
	c.InstrOpcodes[0xe6] = func() {
		c.CycleCount = 5
		c.Inc(c.zeroPageAddress())
	}
	c.InstrOpcodes[0xf6] = func() {
		c.CycleCount = 6
		c.Inc(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0xee] = func() {
		c.CycleCount = 6
		c.Inc(c.absoluteAddress())
	}
	c.InstrOpcodes[0xfe] = func() {
		c.Inc(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 7
	}
	// BRK
	c.InstrOpcodes[0x00] = func() {
		c.CycleCount = 7
		c.Brk()
	}
	// RTI
	c.InstrOpcodes[0x40] = func() {
		c.CycleCount = 6
		c.Rti()
	}
	// RTS
	c.InstrOpcodes[0x60] = func() {
		c.CycleCount = 6
		c.Rts()
	}
	// NOP
	c.InstrOpcodes[0xea] = func() {
		c.CycleCount = 2
	}
	// LSR
	c.InstrOpcodes[0x4a] = func() {
		c.CycleCount = 2
		c.LsrAcc()
	}
	c.InstrOpcodes[0x46] = func() {
		c.CycleCount = 5
		c.Lsr(c.zeroPageAddress())
	}
	c.InstrOpcodes[0x56] = func() {
		c.CycleCount = 6
		c.Lsr(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x4e] = func() {
		c.CycleCount = 6
		c.Lsr(c.absoluteAddress())
	}
	c.InstrOpcodes[0x5e] = func() {
		c.Lsr(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 7
	}
	// ASL
	c.InstrOpcodes[0x0a] = func() {
		c.CycleCount = 2
		c.AslAcc()
	}
	c.InstrOpcodes[0x06] = func() {
		c.CycleCount = 5
		c.Asl(c.zeroPageAddress())
	}
	c.InstrOpcodes[0x16] = func() {
		c.CycleCount = 6
		c.Asl(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x0e] = func() {
		c.CycleCount = 6
		c.Asl(c.absoluteAddress())
	}
	c.InstrOpcodes[0x1E] = func() {
		c.Asl(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 7
	}
	// ROL
	c.InstrOpcodes[0x2a] = func() {
		c.CycleCount = 2
		c.RolAcc()
	}
	c.InstrOpcodes[0x26] = func() {
		c.CycleCount = 5
		c.Rol(c.zeroPageAddress())
	}
	c.InstrOpcodes[0x36] = func() {
		c.CycleCount = 6
		c.Rol(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x2e] = func() {
		c.CycleCount = 6
		c.Rol(c.absoluteAddress())
	}
	c.InstrOpcodes[0x3e] = func() {
		c.Rol(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 7
	}
	// ROR
	c.InstrOpcodes[0x6a] = func() {
		c.CycleCount = 2
		c.RorAcc()
	}
	c.InstrOpcodes[0x66] = func() {
		c.CycleCount = 5
		c.Ror(c.zeroPageAddress())
	}
	c.InstrOpcodes[0x76] = func() {
		c.CycleCount = 6
		c.Ror(c.zeroPageIndexedAddress(c.X))
	}
	c.InstrOpcodes[0x6e] = func() {
		c.CycleCount = 6
		c.Ror(c.absoluteAddress())
	}
	c.InstrOpcodes[0x7e] = func() {
		c.Ror(c.absoluteIndexedAddress(c.X))
		c.CycleCount = 7
	}
	// BIT
	c.InstrOpcodes[0x24] = func() {
		c.CycleCount = 3
		c.Bit(c.zeroPageAddress())
	}
	c.InstrOpcodes[0x2c] = func() {
		c.CycleCount = 4
		c.Bit(c.absoluteAddress())
	}
}
