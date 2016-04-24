package scientist

import "fmt"

// Result holds information about
// an executed experiment.
type Result struct {
	name string
	// Control is the result of executing the control behavior.
	Control *Observation
	// Candidates are the results of executing all the candidate behaviors.
	Candidates []*Observation
	// Mismatches are the results of behaviors that don't match the control.
	Mistmaches []*Observation
	// Ignored are the results of behaviors that can be ignored.
	Ignored []*Observation
}

// Matches returns true if there are no mismatches and ignored observations.
func (r Result) Matches() bool {
	return len(r.Mistmaches) == 0 && len(r.Ignored) == 0
}

// MismatchError holds the result information
// to inspect when observations don't match.
type MismatchError struct {
	result Result
}

// MismatchResult returns the result of the experiment.
func (m MismatchError) MismatchResult() Result {
	return m.result
}

// Error returns the string representation of the MismatchError.
func (m MismatchError) Error() string {
	return fmt.Sprintf("expriment `%s` has %d mismatched results", m.result.name, len(m.result.Mistmaches))
}
