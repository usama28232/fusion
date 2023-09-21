package fuse

// defines breaker-states
const (
	CLOSED BreakerState = iota
	OPEN
)

// defines default constants like default max-failure-count & max-timeout
const (
	defaultMaxFailureCount = 3
	defaultMaxTimeout      = 5
)
