package sdl

import "time"

var events chan interface{} = make(chan interface{})

// This channel delivers SDL events. Each object received from this channel
// has one of the following types: sdl.QuitEvent, sdl.KeyboardEvent,
// sdl.MouseButtonEvent, sdl.MouseMotionEvent, sdl.ActiveEvent,
// sdl.ResizeEvent, sdl.JoyAxisEvent, sdl.JoyButtonEvent, sdl.JoyHatEvent,
// sdl.JoyBallEvent
var Events <-chan interface{} = events

// Polling interval, in milliseconds
const poll_interval_ms = 10

// Polls SDL events in periodic intervals.
// This function does not return.
func pollEvents() {
	// It is more efficient to create the event-object here once,
	// rather than multiple times within the loop
	event := &Event{}

	for {
		for event.poll() {
			switch event.Type {
			case QUIT:
				events <- *(*QuitEvent)(cast(event))

			case KEYDOWN, KEYUP:
				events <- *(*KeyboardEvent)(cast(event))

			case MOUSEBUTTONDOWN, MOUSEBUTTONUP:
				events <- *(*MouseButtonEvent)(cast(event))

			case MOUSEMOTION:
				events <- *(*MouseMotionEvent)(cast(event))

			case JOYAXISMOTION:
				events <- *(*JoyAxisEvent)(cast(event))

			case JOYBUTTONDOWN, JOYBUTTONUP:
				events <- *(*JoyButtonEvent)(cast(event))

			case JOYHATMOTION:
				events <- *(*JoyHatEvent)(cast(event))

			case JOYBALLMOTION:
				events <- *(*JoyBallEvent)(cast(event))

			case ACTIVEEVENT:
				events <- *(*ActiveEvent)(cast(event))

			case VIDEORESIZE:
				events <- *(*ResizeEvent)(cast(event))
			}
		}

		time.Sleep(poll_interval_ms * 1e6)
	}
}

func init() {
	go pollEvents()
}
