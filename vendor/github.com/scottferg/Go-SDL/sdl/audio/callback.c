/*
 * Copyright: ⚛ <0xe2.0x9a.0x9b@gmail.com> 2010
 *
 * The contents of this file can be used freely,
 * except for usages in immoral contexts.
 */

#include "callback.h"
#include <stdio.h>
#include <pthread.h>
#include <assert.h>
#include <string.h>

#define TRUE  1
#define FALSE 0

// Uncomment the next line to enable debugging messages
//#define DEBUG

static pthread_mutex_t m         = PTHREAD_MUTEX_INITIALIZER;
static pthread_cond_t  need      = PTHREAD_COND_INITIALIZER;
static pthread_cond_t  avail     = PTHREAD_COND_INITIALIZER;
static size_t          needed    = 0;	// Number of bytes needed by the consumer
static size_t          available = 0;	// Number of bytes available (from the producer)

static Uint8 *stream;	// Communication buffer between the consumer and the producer

#ifdef DEBUG
#include <time.h>
static int64_t get_time() {
	struct timespec ts;
	if(clock_gettime(CLOCK_MONOTONIC, &ts) == 0)
		return 1000000000*(int64_t)ts.tv_sec + (int64_t)ts.tv_nsec;
	else
		return -1;
}
#endif

#ifdef DEBUG
static uint64_t cummulativeLatency = 0;
static unsigned numCallbacks = 0;
#endif

static void SDLCALL callback(void *userdata, Uint8 *_stream, int _len) {
	assert(_len > 0);

	size_t len = (size_t)_len;

	pthread_mutex_lock(&m);
	{
		assert(available == 0);
		stream = _stream;

		{
			#ifdef DEBUG
				int64_t t1 = get_time();
				printf("consumer: t1=%lld µs\n", (long long)t1/1000);
			#endif

			assert(needed == 0);
			#ifdef DEBUG
				printf("consumer: needed <- %zu\n", len);
			#endif
			needed = len;
			pthread_cond_signal(&need);

			#ifdef DEBUG
				printf("consumer: waiting for data\n");
			#endif
			pthread_cond_wait(&avail, &m);
			assert(needed == 0);
			assert(available == len);

			#ifdef DEBUG
				int64_t t2 = get_time();
				printf("consumer: t2=%lld µs\n", (long long)t2/1000);
				if(t1>0 && t2>0) {
					uint64_t latency = t2-t1;
					cummulativeLatency += latency;
					numCallbacks++;
					printf("consumer: latency=%lld µs, avg=%u µs\n",
					       (long long)(latency/1000),
					       (unsigned)(cummulativeLatency/numCallbacks/1000));
				}
			#endif
		}

		#ifdef DEBUG
			printf("consumer: received %zu bytes of data\n", available);
			printf("consumer: available <- 0\n");
		#endif
		available = 0;
		stream = NULL;
	}
	pthread_mutex_unlock(&m);
}

callback_t callback_getCallback() {
	return &callback;
}

void callback_fillBuffer(Uint8 *data, size_t numBytes) {
	size_t sent = 0;

	pthread_mutex_lock(&m);

	while(sent < numBytes) {
		#ifdef DEBUG
			int64_t t = get_time();
			printf("producer: t=%lld µs\n", (long long)t/1000);
		#endif

		if(needed == 0) {
			#ifdef DEBUG
				printf("producer: waiting until data is needed (1)\n");
			#endif
			pthread_cond_wait(&need, &m);

			// Interrupted from 'callback_unblock' ?
			if(needed == 0) {
				#ifdef DEBUG
					printf("producer: interrupted (1)\n");
				#endif
				break;
			}
		}

		assert(stream != NULL);
		assert(needed > 0);

		// Append a chunk of data to the 'stream'
		size_t n = (needed<(numBytes-sent)) ? needed : (numBytes-sent);
		memcpy(stream+available, data+sent, n);
		available += n;
		sent += n;
		needed -= n;

		#ifdef DEBUG
			printf("producer: added %zu bytes, available=%zu\n", n, available);
		#endif

		if(needed == 0) {
			pthread_cond_signal(&avail);
			if(sent < numBytes) {
				#ifdef DEBUG
					printf("producer: waiting until data is needed (2)\n");
				#endif
				pthread_cond_wait(&need, &m);

				// Interrupted from 'callback_unblock' ?
				if(needed == 0) {
					#ifdef DEBUG
						printf("producer: interrupted (2)\n");
					#endif
					break;
				}
			}
			else {
				break;
			}
		}
	}

	pthread_mutex_unlock(&m);
}

void callback_unblock() {
	pthread_mutex_lock(&m);
	if(needed > 0) {
		// Note: SDL already prefilled the entire 'stream' with silence
		assert(stream != NULL);
		available += needed;
		needed = 0;
		pthread_cond_signal(&avail);
	}
	pthread_cond_signal(&need);
	pthread_mutex_unlock(&m);
}

