# Introduction

[![GoDoc](https://godoc.org/github.com/calavera/go-scientist?status.svg)](https://godoc.org/github.com/calavera/go-scientist)


scientist helps you refactor your Go code with confidence.

Start by creating a new experiment:

```go
experiment := scientist.NewQuickExperiment()
```

Wrap the current behavior into the control function:

```go
// I wonder why this code is so slow :/
control := func(ctx context.Context) (interface{}, error) {
	time.Sleep(10000000 * time.Second)
	return "done", nil
}

experiment.Use(control)
```

Then, create one or more candidate behaviors to compare results:

```go
// This is slightly faster, but I'm getting different results :(
slightlyFasterButWrongResult := func(ctx context.Context) (interface{}, error) {
	time.Sleep(1 * time.Second)
	return "exit", nil
}

experiment.Try("slightly faster call", slightlyFasterButWrongResult)

// I think this is what I want \m/
superFast := func(ctx context.Context) (interface{}, error) {
	return "done", nil
}

experiment.Try("super fast call", superFast)
```

Finally, run the experiment:

```go
value, err := scientist.Run(experiment)
```

This call always returns the result of calling the control function.
It randomizes the call between all three behaviors and measures their duration.
It compares the results and publishes all this information somewhere else to analyze.

## Creating your own experiments

You can create your own experiments by implementing the interface `Experiment`.
The easiest way to do this is by composing your own experiments with `QuickExperiment`
and implementing the methods you want to change, most likely `Name`, `IsEnabled`, `Ignore`,
`Compare` and `Publish`. You can see several examples of this in the `samples` package.

## Failing with mismatches

`scientist.Run` guarantees that the control behavior, your old code, always returns its values.
It might be useful, mostly on testing, to fail the execution when the behaviors don't match,
that way you can test that your experiments are more robust.
To enable this, you can set the global variable `scientist.ErrorOnMismatch` to `true`.

In case of mismatched observations, `scientist.Run` returns `scientist.MismatchResult` as error,
giving you access to all the information about the observations.

## Adding context information

Giving extra information to your experiments is easy using a `context.Context` object.
Use `scientist.RunWithContext` to run your experiment and each behavior will get a copy
of your context object to gather more information.

```go
ctx := context.Background()
ctx = context.WithValue(ctx, "user", models.User{})

control := func(ctx context.Context) (interface{}, error) {
	return ctx.Value("user").(models.User).Login, nil
}

experiment := scientist.NewQuickExperiment()
experiment.Use(control)

login, err := scientist.RunWithContext(ctx, experiment)
```


This package was inspired by GitHub's ruby scientist: https://github.com/github/scientist.
