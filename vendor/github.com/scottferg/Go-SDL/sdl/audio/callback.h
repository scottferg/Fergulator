/*
 * Copyright: ⚛ <0xe2.0x9a.0x9b@gmail.com> 2010
 *
 * The contents of this file can be used freely,
 * except for usages in immoral contexts.
 */

#include <SDL_audio.h>

typedef void (SDLCALL *callback_t)(void *userdata, Uint8 *stream, int len);

extern callback_t callback_getCallback();
extern void       callback_fillBuffer(Uint8 *data, size_t numBytes);
extern void       callback_unblock();

