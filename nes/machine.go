package nes

var (
<<<<<<< HEAD
	cpu     Cpu
	ppu     Ppu
	apu     Apu
	rom     Mapper
	Ram     Memory
	Pads    [2]*Controller
=======
	cpu  Cpu
	ppu  Ppu
	apu  Apu
	rom  Mapper
	Ram  Memory
	Pads [2]*Controller

	Running = true
>>>>>>> 4c75604944e216dde4f118842923b7ed5c982af7
)

// Main system runloop. This should be run on it's own goroutine
func RunSystem() {
	var cycles int
	// var lastApuTick int
	// var flip int

	for Running {
		var frame int
		for frame < 81840 {
			for cycles <= 114 {
				cycles += cpu.Step()
				totalCpuCycles += cycles
			}

			for i := 0; i < 341; i++ {
				ppu.Step()
			}

			frame += cycles * 3
			cycles -= 114
		}

		/*
			for i := 0; i < 81840; i++ {
				apu.Step()
			}

			if AudioEnabled {
				if totalCpuCycles-apu.LastFrameTick >= (cpuClockSpeed / 240) {
					apu.FrameSequencerStep()
					apu.LastFrameTick = totalCpuCycles
				}

				if totalCpuCycles-lastApuTick >= ((cpuClockSpeed / 44100) + flip) {
					apu.PushSample()
					lastApuTick = totalCpuCycles

					flip = (flip + 1) & 0x1
				}
			}
		*/
	}
}
