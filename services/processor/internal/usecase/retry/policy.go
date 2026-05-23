package retry

import (
	"context"
	"math/rand"
	"processor/internal/infra/config"
	"time"
)

type Policy struct {
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Jitter         time.Duration
	MaxAttempts    int
}

func NewPolicy(cfg config.Settings) Policy {
	return Policy{
		InitialBackoff: cfg.ReceiveBackoff.InitialBackoff,
		MaxBackoff:     cfg.ReceiveBackoff.MaxBackoff,
		Jitter:         cfg.ReceiveBackoff.Jitter,
		MaxAttempts:    cfg.ReceiveBackoff.MaxAttempts,
	}
}

func (p Policy) Normalize() Policy {
	if p.MaxAttempts <= 0 {
		p.MaxAttempts = 1
	}
	if p.InitialBackoff < 0 {
		p.InitialBackoff = 0
	}
	if p.MaxBackoff < p.InitialBackoff {
		p.MaxBackoff = p.InitialBackoff
	}
	if p.Jitter < 0 {
		p.Jitter = 0
	}
	return p
}

func (p Policy) Delay(attempt int) time.Duration {
	p = p.Normalize()
	delay := p.InitialBackoff
	for i := 1; i < attempt; i++ {
		delay *= 2
		if delay >= p.MaxBackoff {
			delay = p.MaxBackoff
			break
		}
	}
	if p.Jitter > 0 {
		delay += time.Duration(rand.Int63n(int64(p.Jitter) + 1))
	}
	return delay
}

func (p Policy) Wait(ctx context.Context, attempt int) error {
	delay := p.Delay(attempt)
	if delay <= 0 {
		return nil
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
