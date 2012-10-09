package main

import (
	"encoding/hex"
	"io/ioutil"
	"strings"
	"testing"
)

type CpuState struct {
	A  int
	X  int
	Y  int
	P  int
	S  int
    C  int
	Op uint16
}

func TestGoldLog(test *testing.T) {
	ProgramCounter = 0xC000

	Ram.Init()
	cpu.Reset()
	ppu.Init()

	cpu.P = 0x24

	cpu.Accurate = false

	if contents, err := ioutil.ReadFile("test_roms/nestest.nes"); err == nil {
		if rom, err = LoadRom(contents); err != nil {
			test.Error(err.Error())
			return
		}
	}

	logfile, err := ioutil.ReadFile("test_roms/nestest.log")
	if err != nil {
		test.Error(err.Error())
		return
	}

	log := strings.Split(string(logfile), "\n")

    sentinel := 100
	//sentinel := 5003
	for i := 0; i < sentinel; i++ {
		op, _ := hex.DecodeString(log[i][:4])

		high := op[0]
		low := op[1]

		r := log[i][48:]

		registers := strings.Split(r, " ")

		a, _ := hex.DecodeString(strings.Split(registers[0], ":")[1])
		x, _ := hex.DecodeString(strings.Split(registers[1], ":")[1])
		y, _ := hex.DecodeString(strings.Split(registers[2], ":")[1])
		p, _ := hex.DecodeString(strings.Split(registers[3], ":")[1])
		sp, _ := hex.DecodeString(strings.Split(registers[4], ":")[1])

		expectedState := CpuState{
			A:  int(a[0]),
			X:  int(x[0]),
			Y:  int(y[0]),
			P:  int(p[0]),
			S:  int(sp[0]),
			Op: (uint16(high) << 8) + uint16(low),
		}

		verifyCpuState(ProgramCounter, &cpu, expectedState, test)
		cycles := cpu.Step()

		// 3 PPU cycles for each CPU cycle
		for i := 0; i < 3*cycles; i++ {
			ppu.Step()
		}
	}
}

func verifyCpuState(pc uint16, c *Cpu, e CpuState, test *testing.T) {
	if pc != e.Op {
		test.Errorf("PC was 0x%X, expected 0x%X\n", pc, e.Op)
	}

	if c.A != Word(e.A) {
		test.Errorf("PC: 0x%X Register A was 0x%X, was expecting 0x%X\n", pc, c.A, e.A)
	}

	if c.X != Word(e.X) {
		test.Errorf("PC: 0x%X Register X was 0x%X, was expecting 0x%X\n", pc, c.X, e.X)
	}

	if c.Y != Word(e.Y) {
		test.Errorf("PC: 0x%X Register Y was 0x%X, was expecting 0x%X\n", pc, c.Y, e.Y)
	}

	if c.P != Word(e.P) {
		test.Errorf("PC: 0x%X P register was 0x%X, was expecting 0x%X\n", pc, c.P, e.P)
	}

	if c.StackPointer != Word(e.S) {
		test.Errorf("PC: 0x%X Stack pointer was 0x%X, was expecting 0x%X\n", pc, c.StackPointer, e.S)
	}
}
