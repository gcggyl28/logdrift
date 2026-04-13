// Package buffer provides a fixed-size ring buffer for log lines per service.
// It is safe for concurrent use.
package buffer

import (
	"errors"
	"sync"
)

// ErrInvalidCapacity is returned when capacity is less than 1.
var ErrInvalidCapacity = errors.New("buffer: capacity must be >= 1")

// Entry holds a single log line associated with a service name.
type Entry struct {
	Service string
	Line    string
}

// RingBuffer is a thread-safe, fixed-capacity circular buffer of log entries.
type RingBuffer struct {
	mu       sync.Mutex
	entries  []Entry
	cap      int
	head     int
	size     int
}

// New creates a RingBuffer with the given capacity.
func New(capacity int) (*RingBuffer, error) {
	if capacity < 1 {
		return nil, ErrInvalidCapacity
	}
	return &RingBuffer{
		entries: make([]Entry, capacity),
		cap:     capacity,
	}, nil
}

// Push adds an entry to the buffer, overwriting the oldest entry when full.
func (r *RingBuffer) Push(e Entry) {
	r.mu.Lock()
	defer r.mu.Unlock()

	index := (r.head + r.size) % r.cap
	r.entries[index] = e
	if r.size < r.cap {
		r.size++
	} else {
		r.head = (r.head + 1) % r.cap
	}
}

// Entries returns a snapshot of all buffered entries in insertion order.
func (r *RingBuffer) Entries() []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]Entry, r.size)
	for i := 0; i < r.size; i++ {
		out[i] = r.entries[(r.head+i)%r.cap]
	}
	return out
}

// Len returns the current number of entries stored.
func (r *RingBuffer) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.size
}

// Cap returns the maximum capacity of the buffer.
func (r *RingBuffer) Cap() int {
	return r.cap
}

// Reset clears all entries without reallocating.
func (r *RingBuffer) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.head = 0
	r.size = 0
}
