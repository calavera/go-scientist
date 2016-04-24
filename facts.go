package scientist

import (
	"math/rand"
	"time"
)

// Facts holds behavior information for an experiment.
type Facts struct {
	behaviors       map[string]Behavior
	behaviorsAccess []string
}

// NewFacts create a new set of facts for a experiment.
func NewFacts() *Facts {
	return &Facts{
		behaviors: make(map[string]Behavior),
	}
}

// Name returns the facts name as `experiment`.
func (f *Facts) Name() string {
	return "experiment"
}

// Control returns the control behavior.
func (f *Facts) Control() Behavior {
	return f.behaviors["__control__"]
}

// Behavior returns a candidate behavior by its name.
func (f *Facts) Behavior(name string) Behavior {
	return f.behaviors[name]
}

// Use sets the control behavior.
func (f *Facts) Use(behavior Behavior) error {
	return f.tryBehavior("__control__", behavior)
}

// Try adds a new candidate behavior.
// The name of each candidate must be unique.
func (f *Facts) Try(name string, behavior Behavior) error {
	return f.tryBehavior(name, behavior)
}

// Shuffle randomizes the behavior access.
func (f *Facts) Shuffle() []string {
	t := time.Now()
	rand.Seed(int64(t.Nanosecond()))

	arr := f.behaviorsAccess
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i)
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func (f *Facts) tryBehavior(name string, behavior Behavior) error {
	if _, exist := f.behaviors[name]; exist {
		return behaviorAlreadyExist{name}
	}

	f.behaviors[name] = behavior
	f.behaviorsAccess = append(f.behaviorsAccess, name)

	return nil
}
