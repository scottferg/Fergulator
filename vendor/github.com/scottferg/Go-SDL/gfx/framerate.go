/* 
A pure Go version of SDL_framerate
*/

package gfx

import (
	"time"
)

type FPSmanager struct {
	framecount uint32
	rateticks  float64
	lastticks  uint64
	rate       uint32
}

func NewFramerate() *FPSmanager {
	return &FPSmanager{
		framecount: 0,
		rate:       FPS_DEFAULT,
		rateticks:  (1000.0 / float64(FPS_DEFAULT)),
		lastticks:  uint64(time.Now().UnixNano()) / 1e6,
	}
}

func (manager *FPSmanager) SetFramerate(rate uint32) {
	if rate >= FPS_LOWER_LIMIT && rate <= FPS_UPPER_LIMIT {
		manager.framecount = 0
		manager.rate = rate
		manager.rateticks = 1000.0 / float64(rate)
	} else {
	}
}

func (manager *FPSmanager) GetFramerate() uint32 {
	return manager.rate
}

func (manager *FPSmanager) FramerateDelay() {
	var current_ticks, target_ticks, the_delay uint64

	// next frame
	manager.framecount++

	// get/calc ticks
	current_ticks = uint64(time.Now().UnixNano()) / 1e6
	target_ticks = manager.lastticks + uint64(float64(manager.framecount)*manager.rateticks)

	if current_ticks <= target_ticks {
		the_delay = target_ticks - current_ticks
		time.Sleep(time.Duration(the_delay * 1e6))
	} else {
		manager.framecount = 0
		manager.lastticks = uint64(time.Now().UnixNano()) / 1e6
	}
}
