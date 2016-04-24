package scientist

import "golang.org/x/net/context"

// Experiment is an interface that defines
// how an experiment behaves.
type Experiment interface {
	Name() string
	Control() Behavior
	Shuffle() []string
	Behavior(name string) Behavior
	IsEnabled(ctx context.Context) bool
	Ignore(ctx context.Context, control, candidate *Observation) bool
	Compare(ctx context.Context, control, candidate *Observation) bool
	Publish(ctx context.Context, result Result) error
}

// QuickExperiment is an experiment with a very basic behavior.
// It's always enabled and it does not publishes results anywhere.
type QuickExperiment struct {
	*Facts
}

// NewQuickExperiment creates a new Experiment with a given name.
// It creates an empty context for the experiment.
func NewQuickExperiment() QuickExperiment {
	return QuickExperiment{
		Facts: NewFacts(),
	}
}

// IsEnabled returns true if the experiment is enabled.
// If it's not enabled, the experiment only runs the
// control behavior.
func (e QuickExperiment) IsEnabled(ctx context.Context) bool {
	return true
}

// Ignore returns true if a candidate behavior can be ignored.
// By default there are no behaviors ignored.
func (e QuickExperiment) Ignore(ctx context.Context, control, candidate *Observation) bool {
	return false
}

// Compare returns true if the result of the control behavior is the same
// as the result of a candidate behavior.
func (e QuickExperiment) Compare(ctx context.Context, control, candidate *Observation) bool {
	return control.Error == candidate.Error && control.Value == candidate.Value
}

// Publish allows you to export the result of the experiment somewhere else.
// Use it to compare result information between control and candidates.
func (e QuickExperiment) Publish(ctx context.Context, result Result) error {
	return nil
}
