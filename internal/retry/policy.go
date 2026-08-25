package retry

import "time"

type Policy struct {
	Base time.Duration
	Cap  time.Duration
}

func DefaultPolicy() Policy {
	return Policy{Base: time.Second, Cap: time.Minute}
}

func (p Policy) Delay(attempt int) time.Duration {
	if attempt <= 0 {
		return p.Base
	}
	d := p.Base * time.Duration(1<<uint(attempt-1))
	if d > p.Cap {
		return p.Cap
	}
	return d
}
