/*
Package scientist helps you refactor your Go code with confidence.

Start by creating a new experiment:

	experiment := scientist.NewQuickExperiment()

Wrap the current behavior into the control function:

	// I wonder why this code is so slow :/
	control := func(ctx context.Context) (interface{}, error) {
		time.Sleep(10000000 * time.Second)
		return "done", nil
	}
	experiment.Use(control)

Then, create one or more candidate behaviors to compare results:

	// This is slightly faster, but I'm getting different results :(
	slightlyFasterButWrongResult := func(ctx context.Context) (interface{}, error) {
		time.Sleep(1 * time.Second)
		return "exit", nil
	}
	experiment.Try("slightly faster call", slightlyFasterWrongResult)

	// I think this is what I want \m/
	superFast := func(ctx context.Context) (interface{}, error) {
		return "done", nil
	}
	experiment.Try("super fast call", superFast)

Finally, run the experiment:

	value, err := scientist.Run(experiment)

This call always returns the result of calling the control function.
It randomizes the call between all three behaviors and measures their duration.
It compares the results and publishes all this information somewhere else to analyze.

Creating your own experiments

You can create your own experiments by implementing the interface `Experiment`.
The easiest way to do this is by composing your own experiments with `QuickExperiment`
and implementing the methods you want to change, most likely `Name`, `IsEnabled`, `Ignore`,
`Compare` and `Publish`. You can see several examples of this in the `samples` package.

Failing with mismatches

`scientist.Run` guarantees that the control behavior, your old code, always returns its values.
It might be useful, mostly on testing, to fail the execution when the behaviors don't match,
that way you can test that your experiments are more robust.
To enable this, you can set the global variable `scientist.ErrorOnMismatch` to `true`.

In case of mismatched observations, `scientist.Run` returns `scientist.MismatchResult` as error,
giving you access to all the information about the observations.

Adding context information

Giving extra information to your experiments is easy using a `context.Context` object.
Use `scientist.RunWithContext` to run your experiment and each behavior will get a copy
of your context object to gather more information.

	ctx := context.Background()
	ctx = context.WithValue(ctx, "user", models.User{})
	control := func(ctx context.Context) (interface{}, error) {
		return ctx.Value("user").(models.User).Login, nil
	}
	experiment := scientist.NewQuickExperiment()
	experiment.Use(control)
	login, err := scientist.RunWithContext(ctx, experiment)

This package was inspired by GitHub's ruby scientist: https://github.com/github/scientist.
*/
package scientist

import (
	"sync"
	"time"

	"golang.org/x/net/context"
)

// ErrorOnMismatch tells scientist to return
// errors when experiments have mismatches.
// Use this to make your tests fail while
// preserving the control candidate behavior
// intact in production.
var ErrorOnMismatch = false

// Behavior is the type of function that defines how
// your experiment behaves. See Experiment.Use and
// Experiment.Try to set those behaviors.
type Behavior func(context.Context) (interface{}, error)

// Run executes the experiment and publishes the results.
// It always returns the result of the control behavior, unless
// ErrorOnMismatch is true and there are mismatches.
// The order of execution between control and candidates
// is always random.
func Run(e Experiment) (interface{}, error) {
	return RunWithContext(context.Background(), e)
}

// RunWithContext executes the experiment and publishes the results.
// It allows to set additional information via the context object.
// It always returns the result of the control behavior, unless
// ErrorOnMismatch is true and there are mismatches.
// The order of execution between control and candidates
// is always random.
func RunWithContext(ctx context.Context, e Experiment) (interface{}, error) {
	c := e.Control()

	if c == nil {
		return "", controlDoesNotExist{}
	}

	behaviors := e.Shuffle()

	// run only the control behavior if the
	// experiment is not enabled or there are
	// no more behaviors.
	if !e.IsEnabled(ctx) || len(behaviors) == 1 {
		return c(ctx)
	}

	control, candidates := runExperiment(ctx, e, behaviors)

	result := gatherResult(ctx, e, control, candidates)

	if err := e.Publish(ctx, result); err != nil {
		return nil, err
	}

	if ErrorOnMismatch && len(result.Mistmaches) > 0 {
		return nil, MismatchError{result}
	}

	return control.Value, control.Error
}

func runExperiment(ctx context.Context, e Experiment, behaviors []string) (*Observation, []*Observation) {
	var control *Observation
	var candidates []*Observation
	var wg sync.WaitGroup

	finished := make(chan *Observation, len(behaviors))

	for _, name := range behaviors {
		wg.Add(1)
		go func(ctx context.Context, name string) {
			defer wg.Done()

			b := e.Behavior(name)
			finished <- observe(ctx, name, b)
		}(ctx, name)
	}
	wg.Wait()
	close(finished)

	for o := range finished {
		if o.Name == "__control__" {
			control = o
		} else {
			candidates = append(candidates, o)
		}
	}

	return control, candidates
}

func observe(ctx context.Context, name string, b Behavior) (obs *Observation) {
	o := &Observation{
		Name: name,
	}
	defer func() {
		if r := recover(); r != nil {
			o.Error = recoverFromBadBehavior{name, r}
			obs = o
		}
	}()
	defer func() {
		o.Duration = time.Since(o.Start)
	}()
	o.Start = time.Now()
	o.Value, o.Error = b(ctx)

	return o
}

func gatherResult(ctx context.Context, e Experiment, control *Observation, candidates []*Observation) Result {
	result := Result{
		name:       e.Name(),
		Control:    control,
		Candidates: candidates,
	}

	for _, o := range candidates {
		match := e.Compare(ctx, control, o)

		if !match {
			if e.Ignore(ctx, control, o) {
				result.Ignored = append(result.Ignored, o)
				continue
			}
			result.Mistmaches = append(result.Mistmaches, o)
		}
	}

	return result
}
