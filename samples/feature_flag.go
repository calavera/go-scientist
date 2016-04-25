package samples

import (
	"github.com/calavera/go-scientist"
	"golang.org/x/net/context"
)

type featuresTable map[string]map[int]bool

type featuresExperiment struct {
	scientist.QuickExperiment
	features featuresTable
}

func newFeaturesExperiment(features featuresTable) *featuresExperiment {
	return &featuresExperiment{
		QuickExperiment: scientist.NewQuickExperiment(),
		features:        features,
	}
}

func (e *featuresExperiment) IsEnabled(ctx context.Context) bool {
	user := ctx.Value("user_id").(int)

	_, ok := e.features["new_feature"][user]
	return ok
}

func exampleFeaturesExperiment() {
	features := featuresTable{}
	features["new_feature"] = map[int]bool{
		1: true,
		2: false,
	}

	e := newFeaturesExperiment(features)
	e.Use(func(ctx context.Context) (interface{}, error) {
		return nil, nil
	})

	e.Try("new_feature", func(ctx context.Context) (interface{}, error) {
		return nil, nil
	})

	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_id", 1)

	scientist.RunWithContext(ctx, e)
}
