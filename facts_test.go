package scientist

import (
	"testing"

	"golang.org/x/net/context"
)

func emptyBehavior(ctx context.Context) (interface{}, error) {
	return nil, nil
}

func TestControl(t *testing.T) {
	f := NewFacts()

	if c := f.Control(); c != nil {
		t.Fatal("expected no control")
	}

	if err := f.Use(emptyBehavior); err != nil {
		t.Fatal(err)
	}

	if c := f.Control(); c == nil {
		t.Fatal("expected control, got nil")
	}

	if err := f.Use(emptyBehavior); !IsBehaviorExist(err) {
		t.Fatal("expected duplicated behavior error, got nil")
	}
}

func TestBehavior(t *testing.T) {
	f := NewFacts()

	if b := f.Behavior("try something"); b != nil {
		t.Fatal("expected no behavior")
	}

	if err := f.Try("try something", emptyBehavior); err != nil {
		t.Fatal(err)
	}

	if b := f.Behavior("try something"); b == nil {
		t.Fatal("expected behavior, got nil")
	}

	if err := f.Try("try something", emptyBehavior); !IsBehaviorExist(err) {
		t.Fatal("expected duplicated behavior error, got nil")
	}
}
