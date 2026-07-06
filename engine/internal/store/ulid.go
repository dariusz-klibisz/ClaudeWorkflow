package store

import (
	"crypto/rand"
	"sync"
	"time"
)

// ULID: 48-bit ms timestamp + 80-bit randomness, Crockford base32 (26 chars).
// Implemented in-house to keep the engine dependency-free (07 §3 R2).
// Monotonic within a process: same-millisecond calls increment the random
// component, so ordering by id is stable per writer.

const crockford = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

var (
	ulidMu   sync.Mutex
	lastMs   uint64
	lastRand [10]byte
)

// NewULID returns a 26-char ULID for the current time.
func NewULID() string {
	return newULIDAt(uint64(time.Now().UnixMilli()))
}

func newULIDAt(ms uint64) string {
	ulidMu.Lock()
	defer ulidMu.Unlock()
	if ms == lastMs {
		// increment the 80-bit random component (big-endian)
		for i := 9; i >= 0; i-- {
			lastRand[i]++
			if lastRand[i] != 0 {
				break
			}
		}
	} else {
		lastMs = ms
		if _, err := rand.Read(lastRand[:]); err != nil {
			// crypto/rand never fails on supported platforms; fall back to time bits
			for i := range lastRand {
				lastRand[i] = byte(ms >> (uint(i%8) * 8))
			}
		}
	}
	var b [16]byte
	b[0] = byte(ms >> 40)
	b[1] = byte(ms >> 32)
	b[2] = byte(ms >> 24)
	b[3] = byte(ms >> 16)
	b[4] = byte(ms >> 8)
	b[5] = byte(ms)
	copy(b[6:], lastRand[:])

	// 16 bytes = 128 bits -> 26 base32 chars (first char from 2 bits padding)
	var out [26]byte
	// encode 128 bits big-endian, 5 bits at a time from the top (padded to 130)
	var acc uint32
	bits := 0
	pos := 0
	// prepend 2 zero bits
	acc = 0
	bits = 2
	for _, by := range b {
		acc = acc<<8 | uint32(by)
		bits += 8
		for bits >= 5 {
			bits -= 5
			out[pos] = crockford[(acc>>uint(bits))&31]
			pos++
			acc &= (1 << uint(bits)) - 1
		}
	}
	return string(out[:])
}
