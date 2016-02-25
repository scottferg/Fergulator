/*
A binding of SDL and SDL_image.

The binding works in pretty much the same way as it does in C, although
some of the functions have been altered to give them an object-oriented
flavor (eg. Rather than sdl.Flip(surface) it's surface.Flip() )
*/
package sdl

// #cgo pkg-config: sdl SDL_image
//
// struct private_hwdata{};
// struct SDL_BlitMap{};
// #define map _map
//
// #include <SDL/SDL.h>
// #include <SDL/SDL_image.h>
import "C"

import (
	"os"
	"reflect"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

type cast unsafe.Pointer

// Mutex for serialization of access to certain SDL functions.
//
// There is no need to use this in application code, the mutex is a public variable
// just because it needs to be accessible from other parts of Go-SDL (such as package "sdl/ttf").
//
// Surface-level functions (such as 'Surface.Blit') are not using this mutex,
// so it is possible to modify multiple surfaces concurrently.
// There is no dependency between 'Surface.Lock' and the global mutex.
var GlobalMutex sync.Mutex

type Surface struct {
	cSurface *C.SDL_Surface
	mutex    sync.RWMutex

	Flags  uint32
	Format *PixelFormat
	W      int32
	H      int32
	Pitch  uint16
	Pixels unsafe.Pointer
	Offset int32

	gcPixels interface{} // Prevents garbage collection of pixels passed to func CreateRGBSurfaceFrom
}

func wrap(cSurface *C.SDL_Surface) *Surface {
	var s *Surface

	if cSurface != nil {
		var surface Surface
		surface.SetCSurface(unsafe.Pointer(cSurface))
		s = &surface
	} else {
		s = nil
	}

	return s
}

// FIXME: Ideally, this should NOT be a public function, but it is needed in the package "ttf" ...
func (s *Surface) SetCSurface(cSurface unsafe.Pointer) {
	s.cSurface = (*C.SDL_Surface)(cSurface)
	s.reload()
}

// Pull data from C.SDL_Surface.
// Make sure to use this when the C surface might have been changed.
func (s *Surface) reload() {
	s.Flags = uint32(s.cSurface.flags)
	s.Format = (*PixelFormat)(cast(s.cSurface.format))
	s.W = int32(s.cSurface.w)
	s.H = int32(s.cSurface.h)
	s.Pitch = uint16(s.cSurface.pitch)
	s.Pixels = s.cSurface.pixels
	s.Offset = int32(s.cSurface.offset)
}

func (s *Surface) destroy() {
	s.cSurface = nil
	s.Format = nil
	s.Pixels = nil
	s.gcPixels = nil
}

// =======
// General
// =======

// The version of Go-SDL bindings.
// The version descriptor changes into a new unique string
// after a semantically incompatible Go-SDL update.
//
// The returned value can be checked by users of this package
// to make sure they are using a version with the expected semantics.
//
// If Go adds some kind of support for package versioning, this function will go away.
func GoSdlVersion() string {
	return "⚛SDL bindings 1.0"
}

// Initializes SDL.
func Init(flags uint32) int {
	GlobalMutex.Lock()
	status := int(C.SDL_Init(C.Uint32(flags)))
	if (status != 0) && (runtime.GOOS == "darwin") && (flags&INIT_VIDEO != 0) {
		if os.Getenv("SDL_VIDEODRIVER") == "" {
			os.Setenv("SDL_VIDEODRIVER", "x11")
			status = int(C.SDL_Init(C.Uint32(flags)))
			if status != 0 {
				os.Setenv("SDL_VIDEODRIVER", "")
			}
		}
	}

	GlobalMutex.Unlock()
	return status
}

// Shuts down SDL
func Quit() {
	GlobalMutex.Lock()

	if currentVideoSurface != nil {
		currentVideoSurface.destroy()
		currentVideoSurface = nil
	}

	C.SDL_Quit()

	GlobalMutex.Unlock()
}

// Initializes subsystems.
func InitSubSystem(flags uint32) int {
	GlobalMutex.Lock()
	status := int(C.SDL_InitSubSystem(C.Uint32(flags)))
	if (status != 0) && (runtime.GOOS == "darwin") && (flags&INIT_VIDEO != 0) {
		if os.Getenv("SDL_VIDEODRIVER") == "" {
			os.Setenv("SDL_VIDEODRIVER", "x11")
			status = int(C.SDL_InitSubSystem(C.Uint32(flags)))
			if status != 0 {
				os.Setenv("SDL_VIDEODRIVER", "")
			}
		}
	}
	GlobalMutex.Unlock()
	return status
}

// Shuts down a subsystem.
func QuitSubSystem(flags uint32) {
	GlobalMutex.Lock()
	C.SDL_QuitSubSystem(C.Uint32(flags))
	GlobalMutex.Unlock()
}

// Checks which subsystems are initialized.
func WasInit(flags uint32) int {
	GlobalMutex.Lock()
	status := int(C.SDL_WasInit(C.Uint32(flags)))
	GlobalMutex.Unlock()
	return status
}

// ==============
// Error Handling
// ==============

// Gets SDL error string
func GetError() string {
	GlobalMutex.Lock()
	s := C.GoString(C.SDL_GetError())
	GlobalMutex.Unlock()
	return s
}

// Clear the current SDL error
func ClearError() {
	GlobalMutex.Lock()
	C.SDL_ClearError()
	GlobalMutex.Unlock()
}

// ======
// Video
// ======

var currentVideoSurface *Surface = nil

// Sets up a video mode with the specified width, height, bits-per-pixel and
// returns a corresponding surface.  You don't need to call the Free method
// of the returned surface, as it will be done automatically by sdl.Quit.
func SetVideoMode(w int, h int, bpp int, flags uint32) *Surface {
	GlobalMutex.Lock()
	var screen = C.SDL_SetVideoMode(C.int(w), C.int(h), C.int(bpp), C.Uint32(flags))
	currentVideoSurface = wrap(screen)
	GlobalMutex.Unlock()
	return currentVideoSurface
}

// Returns a pointer to the current display surface.
func GetVideoSurface() *Surface {
	GlobalMutex.Lock()
	surface := currentVideoSurface
	GlobalMutex.Unlock()
	return surface
}

// Checks to see if a particular video mode is supported.  Returns 0 if not
// supported, or the bits-per-pixel of the closest available mode.
func VideoModeOK(width int, height int, bpp int, flags uint32) int {
	GlobalMutex.Lock()
	status := int(C.SDL_VideoModeOK(C.int(width), C.int(height), C.int(bpp), C.Uint32(flags)))
	GlobalMutex.Unlock()
	return status
}

// Returns the list of available screen dimensions for the given format.
//
// NOTE: The result of this function uses a different encoding than the underlying C function.
// It returns an empty array if no modes are available,
// and nil if any dimension is okay for the given format.
func ListModes(format *PixelFormat, flags uint32) []Rect {
	modes := C.SDL_ListModes((*C.SDL_PixelFormat)(cast(format)), C.Uint32(flags))

	// No modes available
	if modes == nil {
		return make([]Rect, 0)
	}

	// (modes == -1) --> Any dimension is ok
	if uintptr(unsafe.Pointer(modes))+1 == uintptr(0) {
		return nil
	}

	count := 0
	ptr := *modes //first element in the list
	for ptr != nil {
		count++
		ptr = *(**C.SDL_Rect)(unsafe.Pointer(uintptr(unsafe.Pointer(modes)) + uintptr(count*int(unsafe.Sizeof(ptr)))))
	}

	ret := make([]Rect, count)
	for i := 0; i < count; i++ {
		ptr := (**C.SDL_Rect)(unsafe.Pointer(uintptr(unsafe.Pointer(modes)) + uintptr(i*int(unsafe.Sizeof(*modes)))))
		var r *C.SDL_Rect = *ptr
		ret[i].X = int16(r.x)
		ret[i].Y = int16(r.y)
		ret[i].W = uint16(r.w)
		ret[i].H = uint16(r.h)
	}

	return ret
}

type VideoInfo struct {
	HW_available bool         "Flag: Can you create hardware surfaces?"
	WM_available bool         "Flag: Can you talk to a window manager?"
	Blit_hw      bool         "Flag: Accelerated blits HW --> HW"
	Blit_hw_CC   bool         "Flag: Accelerated blits with Colorkey"
	Blit_hw_A    bool         "Flag: Accelerated blits with Alpha"
	Blit_sw      bool         "Flag: Accelerated blits SW --> HW"
	Blit_sw_CC   bool         "Flag: Accelerated blits with Colorkey"
	Blit_sw_A    bool         "Flag: Accelerated blits with Alpha"
	Blit_fill    bool         "Flag: Accelerated color fill"
	Video_mem    uint32       "The total amount of video memory (in K)"
	Vfmt         *PixelFormat "Value: The format of the video surface"
	Current_w    int32        "Value: The current video mode width"
	Current_h    int32        "Value: The current video mode height"
}

func GetVideoInfo() *VideoInfo {
	GlobalMutex.Lock()
	vinfo := (*internalVideoInfo)(cast(C.SDL_GetVideoInfo()))
	GlobalMutex.Unlock()

	flags := vinfo.Flags

	return &VideoInfo{
		HW_available: flags&(1<<0) != 0,
		WM_available: flags&(1<<1) != 0,
		Blit_hw:      flags&(1<<9) != 0,
		Blit_hw_CC:   flags&(1<<10) != 0,
		Blit_hw_A:    flags&(1<<11) != 0,
		Blit_sw:      flags&(1<<12) != 0,
		Blit_sw_CC:   flags&(1<<13) != 0,
		Blit_sw_A:    flags&(1<<14) != 0,
		Blit_fill:    flags&(1<<15) != 0,
		Video_mem:    vinfo.Video_mem,
		Vfmt:         vinfo.Vfmt,
		Current_w:    vinfo.Current_w,
		Current_h:    vinfo.Current_h,
	}
}

// Makes sure the given area is updated on the given screen.  If x, y, w, and
// h are all 0, the whole screen will be updated.
func (screen *Surface) UpdateRect(x int32, y int32, w uint32, h uint32) {
	GlobalMutex.Lock()
	screen.mutex.Lock()

	C.SDL_UpdateRect(screen.cSurface, C.Sint32(x), C.Sint32(y), C.Uint32(w), C.Uint32(h))

	screen.mutex.Unlock()
	GlobalMutex.Unlock()
}

func (screen *Surface) UpdateRects(rects []Rect) {
	if len(rects) > 0 {
		GlobalMutex.Lock()
		screen.mutex.Lock()

		C.SDL_UpdateRects(screen.cSurface, C.int(len(rects)), (*C.SDL_Rect)(cast(&rects[0])))

		screen.mutex.Unlock()
		GlobalMutex.Unlock()
	}
}

// Gets the window title and icon name.
func WM_GetCaption() (title, icon string) {
	GlobalMutex.Lock()

	// SDL seems to free these strings.  TODO: Check to see if that's the case
	var ctitle, cicon *C.char
	C.SDL_WM_GetCaption(&ctitle, &cicon)
	title = C.GoString(ctitle)
	icon = C.GoString(cicon)

	GlobalMutex.Unlock()

	return
}

// Sets the window title and icon name.
func WM_SetCaption(title, icon string) {
	ctitle := C.CString(title)
	cicon := C.CString(icon)

	GlobalMutex.Lock()
	C.SDL_WM_SetCaption(ctitle, cicon)
	GlobalMutex.Unlock()

	C.free(unsafe.Pointer(ctitle))
	C.free(unsafe.Pointer(cicon))
}

// Sets the icon for the display window.
func WM_SetIcon(icon *Surface, mask *uint8) {
	GlobalMutex.Lock()
	C.SDL_WM_SetIcon(icon.cSurface, (*C.Uint8)(mask))
	GlobalMutex.Unlock()
}

// Minimizes the window
func WM_IconifyWindow() int {
	GlobalMutex.Lock()
	status := int(C.SDL_WM_IconifyWindow())
	GlobalMutex.Unlock()
	return status
}

// Toggles fullscreen mode
func WM_ToggleFullScreen(surface *Surface) int {
	GlobalMutex.Lock()
	status := int(C.SDL_WM_ToggleFullScreen(surface.cSurface))
	GlobalMutex.Unlock()
	return status
}

// Swaps OpenGL framebuffers/Update Display.
func GL_SwapBuffers() {
	GlobalMutex.Lock()
	C.SDL_GL_SwapBuffers()
	GlobalMutex.Unlock()
}

func GL_SetAttribute(attr int, value int) int {
	GlobalMutex.Lock()
	status := int(C.SDL_GL_SetAttribute(C.SDL_GLattr(attr), C.int(value)))
	GlobalMutex.Unlock()
	return status
}

// Swaps screen buffers.
func (screen *Surface) Flip() int {
	GlobalMutex.Lock()
	screen.mutex.Lock()

	status := int(C.SDL_Flip(screen.cSurface))

	screen.mutex.Unlock()
	GlobalMutex.Unlock()

	return status
}

// Frees (deletes) a Surface
func (screen *Surface) Free() {
	GlobalMutex.Lock()
	screen.mutex.Lock()

	C.SDL_FreeSurface(screen.cSurface)

	screen.destroy()
	if screen == currentVideoSurface {
		currentVideoSurface = nil
	}

	screen.mutex.Unlock()
	GlobalMutex.Unlock()
}

// Locks a surface for direct access.
func (screen *Surface) Lock() int {
	screen.mutex.Lock()
	status := int(C.SDL_LockSurface(screen.cSurface))
	screen.mutex.Unlock()
	return status
}

// Unlocks a previously locked surface.
func (screen *Surface) Unlock() {
	screen.mutex.Lock()
	C.SDL_UnlockSurface(screen.cSurface)
	screen.mutex.Unlock()
}

// Performs a fast blit from the source surface to the destination surface.
// This is the same as func BlitSurface, but the order of arguments is reversed.
func (dst *Surface) Blit(dstrect *Rect, src *Surface, srcrect *Rect) int {
	GlobalMutex.Lock()
	global := true
	if (src != currentVideoSurface) && (dst != currentVideoSurface) {
		GlobalMutex.Unlock()
		global = false
	}

	// At this point: GlobalMutex is locked only if at least one of 'src' or 'dst'
	//                was identical to 'currentVideoSurface'

	var ret C.int
	{
		src.mutex.RLock()
		dst.mutex.Lock()

		ret = C.SDL_UpperBlit(
			src.cSurface,
			(*C.SDL_Rect)(cast(srcrect)),
			dst.cSurface,
			(*C.SDL_Rect)(cast(dstrect)))

		dst.mutex.Unlock()
		src.mutex.RUnlock()
	}

	if global {
		GlobalMutex.Unlock()
	}

	return int(ret)
}

// Performs a fast blit from the source surface to the destination surface.
func BlitSurface(src *Surface, srcrect *Rect, dst *Surface, dstrect *Rect) int {
	return dst.Blit(dstrect, src, srcrect)
}

// This function performs a fast fill of the given rectangle with some color.
func (dst *Surface) FillRect(dstrect *Rect, color uint32) int {
	dst.mutex.Lock()

	var ret = C.SDL_FillRect(
		dst.cSurface,
		(*C.SDL_Rect)(cast(dstrect)),
		C.Uint32(color))

	dst.mutex.Unlock()

	return int(ret)
}

// Adjusts the alpha properties of a Surface.
func (s *Surface) SetAlpha(flags uint32, alpha uint8) int {
	s.mutex.Lock()
	status := int(C.SDL_SetAlpha(s.cSurface, C.Uint32(flags), C.Uint8(alpha)))
	s.mutex.Unlock()
	return status
}

// Sets the color key (transparent pixel)  in  a  blittable  surface  and
// enables or disables RLE blit acceleration.
func (s *Surface) SetColorKey(flags uint32, ColorKey uint32) int {
	s.mutex.Lock()
	status := int(C.SDL_SetColorKey(s.cSurface, C.Uint32(flags), C.Uint32(ColorKey)))
	s.mutex.Unlock()
	return status
}

// Gets the clipping rectangle for a surface.
func (s *Surface) GetClipRect(r *Rect) {
	s.mutex.RLock()
	C.SDL_GetClipRect(s.cSurface, (*C.SDL_Rect)(cast(r)))
	s.mutex.RUnlock()
}

// Sets the clipping rectangle for a surface.
func (s *Surface) SetClipRect(r *Rect) {
	s.mutex.Lock()
	C.SDL_SetClipRect(s.cSurface, (*C.SDL_Rect)(cast(r)))
	s.mutex.Unlock()
}

// Map a RGBA color value to a pixel format.
func MapRGBA(format *PixelFormat, r, g, b, a uint8) uint32 {
	return (uint32)(C.SDL_MapRGBA((*C.SDL_PixelFormat)(cast(format)), (C.Uint8)(r), (C.Uint8)(g), (C.Uint8)(b), (C.Uint8)(a)))
}

// Gets RGBA values from a pixel in the specified pixel format.
func GetRGBA(color uint32, format *PixelFormat, r, g, b, a *uint8) {
	C.SDL_GetRGBA(C.Uint32(color), (*C.SDL_PixelFormat)(cast(format)), (*C.Uint8)(r), (*C.Uint8)(g), (*C.Uint8)(b), (*C.Uint8)(a))
}

// Loads Surface from file (using IMG_Load).
func Load(file string) *Surface {
	GlobalMutex.Lock()

	cfile := C.CString(file)
	var screen = C.IMG_Load(cfile)
	C.free(unsafe.Pointer(cfile))

	GlobalMutex.Unlock()

	return wrap(screen)
}

// Creates an empty Surface.
func CreateRGBSurface(flags uint32, width int, height int, bpp int, Rmask uint32, Gmask uint32, Bmask uint32, Amask uint32) *Surface {
	GlobalMutex.Lock()

	p := C.SDL_CreateRGBSurface(C.Uint32(flags), C.int(width), C.int(height), C.int(bpp),
		C.Uint32(Rmask), C.Uint32(Gmask), C.Uint32(Bmask), C.Uint32(Amask))

	GlobalMutex.Unlock()

	return wrap(p)
}

// Creates a Surface from existing pixel data. It expects pixels to be a slice, pointer or unsafe.Pointer.
func CreateRGBSurfaceFrom(pixels interface{}, width, height, bpp, pitch int, Rmask, Gmask, Bmask, Amask uint32) *Surface {
	var ptr unsafe.Pointer
	switch v := reflect.ValueOf(pixels); v.Kind() {
	case reflect.Ptr, reflect.UnsafePointer, reflect.Slice:
		ptr = unsafe.Pointer(v.Pointer())
	default:
		panic("Don't know how to handle type: " + v.Kind().String())
	}

	GlobalMutex.Lock()
	p := C.SDL_CreateRGBSurfaceFrom(ptr, C.int(width), C.int(height), C.int(bpp), C.int(pitch),
		C.Uint32(Rmask), C.Uint32(Gmask), C.Uint32(Bmask), C.Uint32(Amask))
	GlobalMutex.Unlock()

	s := wrap(p)
	s.gcPixels = pixels
	return s
}

// Converts a surface to the display format
func (s *Surface) DisplayFormat() *Surface {
	s.mutex.RLock()
	p := C.SDL_DisplayFormat(s.cSurface)
	s.mutex.RUnlock()
	return wrap(p)
}

// Converts a surface to the display format with alpha
func (s *Surface) DisplayFormatAlpha() *Surface {
	s.mutex.RLock()
	p := C.SDL_DisplayFormatAlpha(s.cSurface)
	s.mutex.RUnlock()
	return wrap(p)
}

// ========
// Keyboard
// ========

// Enables UNICODE translation.
func EnableUNICODE(enable int) int {
	GlobalMutex.Lock()
	previous := int(C.SDL_EnableUNICODE(C.int(enable)))
	GlobalMutex.Unlock()
	return previous
}

// Sets keyboard repeat rate.
func EnableKeyRepeat(delay, interval int) int {
	GlobalMutex.Lock()
	status := int(C.SDL_EnableKeyRepeat(C.int(delay), C.int(interval)))
	GlobalMutex.Unlock()
	return status
}

// Gets keyboard repeat rate.
func GetKeyRepeat() (int, int) {
	var delay int
	var interval int

	GlobalMutex.Lock()
	C.SDL_GetKeyRepeat((*C.int)(cast(&delay)), (*C.int)(cast(&interval)))
	GlobalMutex.Unlock()

	return delay, interval
}

// Gets a snapshot of the current keyboard state
func GetKeyState() []uint8 {
	GlobalMutex.Lock()

	var numkeys C.int
	array := C.SDL_GetKeyState(&numkeys)

	var ptr = make([]uint8, numkeys)

	*((**C.Uint8)(unsafe.Pointer(&ptr))) = array // TODO

	GlobalMutex.Unlock()

	return ptr

}

// Modifier
type Mod C.int

// Key
type Key C.int

// Gets the state of modifier keys
func GetModState() Mod {
	GlobalMutex.Lock()
	state := Mod(C.SDL_GetModState())
	GlobalMutex.Unlock()
	return state
}

// Sets the state of modifier keys
func SetModState(modstate Mod) {
	GlobalMutex.Lock()
	C.SDL_SetModState(C.SDLMod(modstate))
	GlobalMutex.Unlock()
}

// Gets the name of an SDL virtual keysym
func GetKeyName(key Key) string {
	GlobalMutex.Lock()
	name := C.GoString(C.SDL_GetKeyName(C.SDLKey(key)))
	GlobalMutex.Unlock()
	return name
}

// ======
// Events
// ======

// Polls for currently pending events
func (event *Event) poll() bool {
	GlobalMutex.Lock()

	var ret = C.SDL_PollEvent((*C.SDL_Event)(cast(event)))

	if ret != 0 {
		if (event.Type == VIDEORESIZE) && (currentVideoSurface != nil) {
			currentVideoSurface.reload()
		}
	}

	GlobalMutex.Unlock()

	return ret != 0
}

// =====
// Mouse
// =====

// Retrieves the current state of the mouse.
func GetMouseState(x, y *int) uint8 {
	GlobalMutex.Lock()
	state := uint8(C.SDL_GetMouseState((*C.int)(cast(x)), (*C.int)(cast(y))))
	GlobalMutex.Unlock()
	return state
}

// Retrieves the current state of the mouse relative to the last time this
// function was called.
func GetRelativeMouseState(x, y *int) uint8 {
	GlobalMutex.Lock()
	state := uint8(C.SDL_GetRelativeMouseState((*C.int)(cast(x)), (*C.int)(cast(y))))
	GlobalMutex.Unlock()
	return state
}

// Toggle whether or not the cursor is shown on the screen.
func ShowCursor(toggle int) int {
	GlobalMutex.Lock()
	state := int(C.SDL_ShowCursor((C.int)(toggle)))
	GlobalMutex.Unlock()
	return state
}

// ========
// Joystick
// ========

type Joystick struct {
	cJoystick *C.SDL_Joystick
}

func wrapJoystick(cJoystick *C.SDL_Joystick) *Joystick {
	var j *Joystick
	if cJoystick != nil {
		var joystick Joystick
		joystick.cJoystick = (*C.SDL_Joystick)(unsafe.Pointer(cJoystick))
		j = &joystick
	} else {
		j = nil
	}
	return j
}

// Count the number of joysticks attached to the system
func NumJoysticks() int {
	GlobalMutex.Lock()
	num := int(C.SDL_NumJoysticks())
	GlobalMutex.Unlock()
	return num
}

// Get the implementation dependent name of a joystick.
// This can be called before any joysticks are opened.
// If no name can be found, this function returns NULL.
func JoystickName(deviceIndex int) string {
	GlobalMutex.Lock()
	name := C.GoString(C.SDL_JoystickName(C.int(deviceIndex)))
	GlobalMutex.Unlock()
	return name
}

// Open a joystick for use The index passed as an argument refers to
// the N'th joystick on the system. This index is the value which will
// identify this joystick in future joystick events.  This function
// returns a joystick identifier, or NULL if an error occurred.
func JoystickOpen(deviceIndex int) *Joystick {
	GlobalMutex.Lock()
	joystick := C.SDL_JoystickOpen(C.int(deviceIndex))
	GlobalMutex.Unlock()
	return wrapJoystick(joystick)
}

// Returns 1 if the joystick has been opened, or 0 if it has not.
func JoystickOpened(deviceIndex int) int {
	GlobalMutex.Lock()
	opened := int(C.SDL_JoystickOpened(C.int(deviceIndex)))
	GlobalMutex.Unlock()
	return opened
}

// Update the current state of the open joysticks. This is called
// automatically by the event loop if any joystick events are enabled.
func JoystickUpdate() {
	GlobalMutex.Lock()
	C.SDL_JoystickUpdate()
	GlobalMutex.Unlock()
}

// Enable/disable joystick event polling. If joystick events are
// disabled, you must call SDL_JoystickUpdate() yourself and check the
// state of the joystick when you want joystick information. The state
// can be one of SDL_QUERY, SDL_ENABLE or SDL_IGNORE.
func JoystickEventState(state int) int {
	GlobalMutex.Lock()
	result := int(C.SDL_JoystickEventState(C.int(state)))
	GlobalMutex.Unlock()
	return result
}

// Close a joystick previously opened with SDL_JoystickOpen()
func (joystick *Joystick) Close() {
	GlobalMutex.Lock()
	C.SDL_JoystickClose(joystick.cJoystick)
	GlobalMutex.Unlock()
}

// Get the number of general axis controls on a joystick
func (joystick *Joystick) NumAxes() int {
	return int(C.SDL_JoystickNumAxes(joystick.cJoystick))
}

// Get the device index of an opened joystick.
func (joystick *Joystick) Index() int {
	return int(C.SDL_JoystickIndex(joystick.cJoystick))
}

// Get the number of buttons on a joystick
func (joystick *Joystick) NumButtons() int {
	return int(C.SDL_JoystickNumButtons(joystick.cJoystick))
}

// Get the number of trackballs on a Joystick trackballs have only
// relative motion events associated with them and their state cannot
// be polled.
func (joystick *Joystick) NumBalls() int {
	return int(C.SDL_JoystickNumBalls(joystick.cJoystick))
}

// Get the number of POV hats on a joystick
func (joystick *Joystick) NumHats() int {
	return int(C.SDL_JoystickNumHats(joystick.cJoystick))
}

// Get the current state of a POV hat on a joystick
// The hat indices start at index 0.
func (joystick *Joystick) GetHat(hat int) uint8 {
	return uint8(C.SDL_JoystickGetHat(joystick.cJoystick, C.int(hat)))
}

// Get the current state of a button on a joystick. The button indices
// start at index 0.
func (joystick *Joystick) GetButton(button int) uint8 {
	return uint8(C.SDL_JoystickGetButton(joystick.cJoystick, C.int(button)))
}

// Get the ball axis change since the last poll. The ball indices
// start at index 0. This returns 0, or -1 if you passed it invalid
// parameters.
func (joystick *Joystick) GetBall(ball int, dx, dy *int) int {
	return int(C.SDL_JoystickGetBall(joystick.cJoystick, C.int(ball), (*C.int)(cast(dx)), (*C.int)(cast(dy))))
}

// Get the current state of an axis control on a joystick. The axis
// indices start at index 0. The state is a value ranging from -32768
// to 32767.
func (joystick *Joystick) GetAxis(axis int) int16 {
	return int16(C.SDL_JoystickGetAxis(joystick.cJoystick, C.int(axis)))
}

// ====
// Time
// ====

// Gets the number of milliseconds since the SDL library initialization.
func GetTicks() uint32 {
	GlobalMutex.Lock()
	t := uint32(C.SDL_GetTicks())
	GlobalMutex.Unlock()
	return t
}

// Waits a specified number of milliseconds before returning.
func Delay(ms uint32) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
