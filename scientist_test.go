package scientist

import (
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

type disabledExperiment struct {
	QuickExperiment
}

func (disabledExperiment) IsEnabled(ctx context.Context) bool {
	return false
}

type deepEqualExperiment struct {
	QuickExperiment
}

func (deepEqualExperiment) Compare(ctx context.Context, control, candidate *Observation) bool {
	controlValue := control.Value.([]string)
	candidateValue := candidate.Value.([]string)
	return reflect.DeepEqual(controlValue, candidateValue)
}

func TestRunWithoutControl(t *testing.T) {
	e := NewQuickExperiment()

	_, err := Run(e)
	if !IsControlNotExist(err) {
		t.Fatalf("got %v, expected controlDoesNotExist", err)
	}
}

func TestRunControl(t *testing.T) {
	e := NewQuickExperiment()

	e.Use(func(_ context.Context) (interface{}, error) {
		return "success", nil
	})

	result, err := Run(e)
	if err != nil {
		t.Fatal(err)
	}

	if result != "success" {
		t.Fatalf("run got %v, expected %v", result, "success")
	}
}

func TestRunWithContext(t *testing.T) {
	e := NewQuickExperiment()

	e.Use(func(ctx context.Context) (interface{}, error) {
		return ctx.Value("result"), nil
	})

	ctx := context.Background()
	ctx = context.WithValue(ctx, "result", "success")

	result, err := RunWithContext(ctx, e)
	if err != nil {
		t.Fatal(err)
	}

	if result != "success" {
		t.Fatalf("run got %v, expected %v", result, "success")
	}
}

func TestRunDisableExperiment(t *testing.T) {
	e := disabledExperiment{NewQuickExperiment()}

	ErrorOnMismatch = true
	defer func() { ErrorOnMismatch = false }()

	e.Use(func(_ context.Context) (interface{}, error) {
		return "success", nil
	})

	e.Try("test", func(_ context.Context) (interface{}, error) {
		return "fail", nil
	})

	result, err := Run(e)
	if err != nil {
		t.Fatal(err)
	}

	if result != "success" {
		t.Fatalf("run got %v, expected %v", result, "success")
	}
}

func TestRunErrorOnMisMatch(t *testing.T) {
	e := NewQuickExperiment()

	ErrorOnMismatch = true
	defer func() { ErrorOnMismatch = false }()

	e.Use(func(_ context.Context) (interface{}, error) {
		return "success", nil
	})

	e.Try("test", func(_ context.Context) (interface{}, error) {
		return "fail", nil
	})

	_, err := Run(e)
	if err == nil {
		t.Fatal("expected mismatch error, got nil")
	}

	result := err.(MismatchError).MismatchResult()

	g := len(result.Mistmaches)
	if g != 1 {
		t.Fatalf("mismatches got %v, expected %v, %q", g, 1, result)
	}
}

func TestRunBadBehavior(t *testing.T) {
	e := NewQuickExperiment()

	ErrorOnMismatch = true
	defer func() { ErrorOnMismatch = false }()

	e.Use(func(ctx context.Context) (interface{}, error) {
		return "success", nil
	})

	e.Try("test", func(_ context.Context) (interface{}, error) {
		panic("oh no!")
	})

	_, err := Run(e)
	if err == nil {
		t.Fatal("expected mismatch error, got nil")
	}

	result := err.(MismatchError).MismatchResult()
	m := result.Mistmaches[0]

	if !IsRecoverFromBadBehavior(m.Error) {
		t.Fatalf("got %v, expected recoverFromBadBehavior", m.Error)
	}
}

func TestRunDeepCompare(t *testing.T) {
	e := deepEqualExperiment{NewQuickExperiment()}

	ErrorOnMismatch = true
	defer func() { ErrorOnMismatch = false }()

	e.Use(func(_ context.Context) (interface{}, error) {
		return []string{"1", "2"}, nil
	})

	e.Try("test", func(_ context.Context) (interface{}, error) {
		return []string{"1", "2"}, nil
	})

	result, err := Run(e)
	if err != nil {
		t.Fatal(err)
	}

	w := []string{"1", "2"}
	if !reflect.DeepEqual(result, w) {
		t.Fatalf("run got %v, expected %v", result, w)
	}
}
