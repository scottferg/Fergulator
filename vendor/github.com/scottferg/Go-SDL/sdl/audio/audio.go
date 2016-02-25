/*
 * Copyright: ⚛ <0xe2.0x9a.0x9b@gmail.com> 2010
 *
 * The contents of this file can be used freely,
 * except for usages in immoral contexts.
 */

/*
An interface to low-level SDL sound functions.
*/
package audio

// #cgo pkg-config: sdl
// #cgo freebsd LDFLAGS: -lrt
// #cgo linux LDFLAGS: -lrt
// #cgo windows LDFLAGS: -lpthread
// #include <SDL_audio.h>
// #include "callback.h"
import "C"
import "unsafe"
import "sync"

// The version of Go-SDL audio bindings.
// The version descriptor changes into a new unique string
// after a semantically incompatible Go-SDL update.
//
// The returned value can be checked by users of this package
// to make sure they are using a version with the expected semantics.
//
// If Go adds some kind of support for package versioning, this function will go away.
func GoSdlAudioVersion() string {
	return "⚛SDL audio bindings 1.0"
}

// Audio format
const (
	AUDIO_U8     = C.AUDIO_U8
	AUDIO_S8     = C.AUDIO_S8
	AUDIO_U16LSB = C.AUDIO_U16LSB
	AUDIO_S16LSB = C.AUDIO_S16LSB
	AUDIO_U16MSB = C.AUDIO_U16MSB
	AUDIO_S16MSB = C.AUDIO_S16MSB
	AUDIO_U16    = C.AUDIO_U16
	AUDIO_S16    = C.AUDIO_S16
)

// Native audio byte ordering
const (
	AUDIO_U16SYS = C.AUDIO_U16SYS
	AUDIO_S16SYS = C.AUDIO_S16SYS
)

type AudioSpec struct {
	Freq        int
	Format      uint16 // If in doubt, use AUDIO_S16SYS
	Channels    uint8  // 1 or 2
	Out_Silence uint8
	Samples     uint16 // A power of 2, preferrably 2^11 (2048) or more
	Out_Size    uint32
}

func OpenAudio(desired, obtained_orNil *AudioSpec) int {
	var C_desired, C_obtained *C.SDL_AudioSpec

	C_desired = new(C.SDL_AudioSpec)
	C_desired.freq = C.int(desired.Freq)
	C_desired.format = C.Uint16(desired.Format)
	C_desired.channels = C.Uint8(desired.Channels)
	C_desired.samples = C.Uint16(desired.Samples)
	C_desired.callback = C.callback_getCallback()

	if obtained_orNil != nil {
		if desired != obtained_orNil {
			C_obtained = new(C.SDL_AudioSpec)
		} else {
			C_obtained = C_desired
		}
	}

	status := C.SDL_OpenAudio(C_desired, C_obtained)

	if status == 0 {
		mutex.Lock()
		opened++
		mutex.Unlock()
	}

	if obtained_orNil != nil {
		obtained := obtained_orNil

		obtained.Freq = int(C_obtained.freq)
		obtained.Format = uint16(C_obtained.format)
		obtained.Channels = uint8(C_obtained.channels)
		obtained.Samples = uint16(C_obtained.samples)
		obtained.Out_Silence = uint8(C_obtained.silence)
		obtained.Out_Size = uint32(C_obtained.size)
	}

	return int(status)
}

func CloseAudio() {
	PauseAudio(true)

	mutex.Lock()
	{
		opened--
		switch {
		case opened == 0:
			userPaused = true
			sdlPaused = true
		case opened < 0:
			panic("SDL audio not opened")
		}
	}
	mutex.Unlock()

	C.callback_unblock()

	C.SDL_CloseAudio()
}

// Audio status
const (
	SDL_AUDIO_STOPPED = C.SDL_AUDIO_STOPPED
	SDL_AUDIO_PLAYING = C.SDL_AUDIO_PLAYING
	SDL_AUDIO_PAUSED  = C.SDL_AUDIO_PAUSED
)

func GetAudioStatus() int {
	return int(C.SDL_GetAudioStatus())
}

var opened int = 0

var userPaused bool = true
var sdlPaused bool = true
var haveData bool = false

var mutex sync.Mutex

// Pause or unpause the audio.
// Unpausing is deferred until a SendAudio function receives some samples.
func PauseAudio(pause_on bool) {
	mutex.Lock()

	if pause_on != sdlPaused {
		if pause_on {
			// Pause SDL audio
			userPaused = true
			sdlPaused = true
			C.SDL_PauseAudio(1)
		} else {
			userPaused = false
			if haveData {
				// Unpause SDL audio
				sdlPaused = false
				C.SDL_PauseAudio(0)
			} else {
				// Defer until SendAudio is called
			}
		}
	}

	mutex.Unlock()
}

func LockAudio() {
	C.SDL_LockAudio()
}

func UnlockAudio() {
	C.SDL_UnlockAudio()
}

// Send samples to the audio device (AUDIO_S16SYS format).
// This function blocks until all the samples are consumed by the SDL audio thread.
func SendAudio_int16(data []int16) {
	if len(data) > 0 {
		sendAudio((*C.Uint8)(unsafe.Pointer(&data[0])), C.size_t(int(unsafe.Sizeof(data[0]))*len(data)))
	}
}

// Send samples to the audio device (AUDIO_U16SYS format).
// This function blocks until all the samples are consumed by the SDL audio thread.
func SendAudio_uint16(data []uint16) {
	if len(data) > 0 {
		sendAudio((*C.Uint8)(unsafe.Pointer(&data[0])), C.size_t(int(unsafe.Sizeof(data[0]))*len(data)))
	}
}

// Send samples to the audio device (AUDIO_S8 format).
// This function blocks until all the samples are consumed by the SDL audio thread.
func SendAudio_int8(data []int8) {
	if len(data) > 0 {
		sendAudio((*C.Uint8)(unsafe.Pointer(&data[0])), C.size_t(int(unsafe.Sizeof(data[0]))*len(data)))
	}
}

// Send samples to the audio device (AUDIO_U8 format).
// This function blocks until all the samples are consumed by the SDL audio thread.
func SendAudio_uint8(data []uint8) {
	if len(data) > 0 {
		sendAudio((*C.Uint8)(unsafe.Pointer(&data[0])), C.size_t(int(unsafe.Sizeof(data[0]))*len(data)))
	}
}

func sendAudio(data *C.Uint8, numBytes C.size_t) {
	if numBytes > 0 {
		mutex.Lock()
		{
			haveData = true

			if (userPaused == false) && (sdlPaused == true) {
				// Unpause SDL audio
				sdlPaused = false
				C.SDL_PauseAudio(0)
			}
		}
		mutex.Unlock()

		C.callback_fillBuffer(data, numBytes)

		mutex.Lock()
		haveData = false
		mutex.Unlock()
	}
}
