package scientist

import "fmt"

// IsBehaviorExist returns true if the error
// was caused because a behavior with a given
// name already exists in the experiment.
func IsBehaviorExist(err error) bool {
	_, ok := err.(behaviorAlreadyExist)
	return ok
}

type behaviorAlreadyExist struct {
	name string
}

func (e behaviorAlreadyExist) Error() string {
	return fmt.Sprintf("behavior already exists: %s", e.name)
}

// IsControlNotExist returns true if the experiment
// doesn't have any control behavior.
func IsControlNotExist(err error) bool {
	_, ok := err.(controlDoesNotExist)
	return ok
}

type controlDoesNotExist struct{}

func (e controlDoesNotExist) Error() string {
	return fmt.Sprintf("control behavior doesn't exist. Call experiment.Use to set the control")
}

// IsRecoverFromBadBehavior returns true if one
// of the behaviors panicked.
func IsRecoverFromBadBehavior(err error) bool {
	_, ok := err.(recoverFromBadBehavior)
	return ok
}

type recoverFromBadBehavior struct {
	name  string
	value interface{}
}

func (e recoverFromBadBehavior) Error() string {
	return fmt.Sprintf("recover from bad behavior %s: %v", e.name, e.value)
}
