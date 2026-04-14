// Package sampler provides probabilistic and rate-based log line sampling
// for reducing noise in high-volume log streams.
package sampler

import (
	"errors"
	"math/rand"
	"sync"
)

// Mode controls how sampling decisions are made.
type Mode string

const (
	ModeRandom     Mode = "random"     // keep each line with probability Rate
	ModeEveryN     Mode = "every_n"    // keep every Nth line
)

// Config holds sampler configuration.
type Config struct {
	Mode Mode
	// Rate is the keep probability [0.0, 1.0] for ModeRandom.
	Rate float64
	// N is the interval for ModeEveryN (must be >= 1).
	N int
}

// Sampler decides whether a given log line should be forwarded.
type Sampler struct {
	cfg     Config
	mu      sync.Mutex
	counter int
	rng     *rand.Rand
}

// New creates a Sampler from cfg. Returns an error if the configuration is invalid.
func New(cfg Config) (*Sampler, error) {
	switch cfg.Mode {
	case ModeRandom:
		if cfg.Rate < 0.0 || cfg.Rate > 1.0 {
			return nil, errors.New("sampler: rate must be between 0.0 and 1.0")
		}
	case ModeEveryN:
		if cfg.N < 1 {
			return nil, errors.New("sampler: N must be >= 1")
		}
	default:
		return nil, errors.New("sampler: unknown mode " + string(cfg.Mode))
	}
	return &Sampler{
		cfg: cfg,
		//nolint:gosec // non-cryptographic use
		rng: rand.New(rand.NewSource(rand.Int63())),
	}, nil
}

// Allow returns true if the line should be forwarded.
func (s *Sampler) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.cfg.Mode {
	case ModeRandom:
		return s.rng.Float64() < s.cfg.Rate
	case ModeEveryN:
		s.counter++
		if s.counter >= s.cfg.N {
			s.counter = 0
			return true
		}
		return false
	}
	return true
}

// Reset resets internal counters (useful in tests).
func (s *Sampler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter = 0
}
