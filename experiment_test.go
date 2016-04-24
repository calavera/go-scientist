package scientist

import (
	"testing"

	"golang.org/x/net/context"
)

func TestComputeResultMismatches(t *testing.T) {
	e := NewQuickExperiment()

	cases := []struct {
		control    *Observation
		candidates []*Observation
		mismatches int
	}{
		{
			control:    &Observation{},
			mismatches: 0,
		},
		{
			control: &Observation{Value: true},
			candidates: []*Observation{
				{Value: true},
				{Value: true},
			},
			mismatches: 0,
		},
		{
			control: &Observation{Value: true},
			candidates: []*Observation{
				{Value: false},
				{Value: true},
			},
			mismatches: 1,
		},
		{
			control: &Observation{Value: true},
			candidates: []*Observation{
				{Value: false},
				{Value: false},
			},
			mismatches: 2,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		r := gatherResult(ctx, e, c.control, c.candidates)
		g := len(r.Mistmaches)
		if g != c.mismatches {
			t.Fatalf("mismatches: got %d, expected %d", g, c.mismatches)
		}
	}
}

type ignoreExperiment struct {
	QuickExperiment
}

func (ignoreExperiment) Ignore(ctx context.Context, control, candidate *Observation) bool {
	return true
}

func TestComputeResultIgnored(t *testing.T) {
	e := ignoreExperiment{NewQuickExperiment()}

	cases := []struct {
		control    *Observation
		candidates []*Observation
		mismatches int
		ignored    int
	}{
		{
			control:    &Observation{},
			mismatches: 0,
			ignored:    0,
		},
		{
			control: &Observation{Value: true},
			candidates: []*Observation{
				{Value: true},
				{Value: true},
			},
			mismatches: 0,
			ignored:    0,
		},
		{
			control: &Observation{Value: true},
			candidates: []*Observation{
				{Value: false},
				{Value: true},
			},
			mismatches: 0,
			ignored:    1,
		},
		{
			control: &Observation{Value: true},
			candidates: []*Observation{
				{Value: false},
				{Value: false},
			},
			mismatches: 0,
			ignored:    2,
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		r := gatherResult(ctx, e, c.control, c.candidates)
		g := len(r.Mistmaches)
		if g != c.mismatches {
			t.Fatalf("mismatches: got %d, expected %d", g, c.mismatches)
		}

		g = len(r.Ignored)
		if g != c.ignored {
			t.Fatalf("ignored: got %d, expected %d", g, c.ignored)
		}
	}
}
