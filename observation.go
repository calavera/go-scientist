package scientist

import "time"

// Observation holds information about
// an executed behavior.
type Observation struct {
	// Name is the name of the behavior executed.
	Name string
	// Start is the time when the behavior was executed.
	Start time.Time
	// Duration is the time that take the behavior to run.
	Duration time.Duration
	// Value is the value returned by the behavior if any.
	Value interface{}
	// Error is the error returned by the behavior, if any.
	Error error
}
