package samples

import (
	"fmt"
	"strings"

	"github.com/alexcesaro/statsd"
	"github.com/calavera/go-scientist"
	"golang.org/x/net/context"
)

type MetricsExperiment struct {
	scientist.QuickExperiment
	client *statsd.Client
}

func NewMetricsExperiment(client *statsd.Client) *MetricsExperiment {
	return &MetricsExperiment{
		QuickExperiment: scientist.NewQuickExperiment(),
		client:          client,
	}
}

func (m *MetricsExperiment) Publish(ctx context.Context, result scientist.Result) error {
	control := fmt.Sprintf("scientist.metrics.%s.control.duration", m.Name())
	m.client.Timing(control, int64(result.Control.Duration))

	for _, c := range result.Candidates {
		name := strings.Replace(c.Name, " ", "_", -1)

		candidate := fmt.Sprintf("scientist.metrics.%s.%s.duration", m.Name(), name)
		m.client.Timing(candidate, int64(c.Duration))
	}

	switch {
	case result.Matches():
		m.client.Increment(fmt.Sprintf("scientist.metrics.%s.matched", m.Name()))
	case len(result.Ignored) > 0:
		m.client.Increment(fmt.Sprintf("scientist.metrics.%s.ignored", m.Name()))
	default:
		m.client.Increment(fmt.Sprintf("scientist.metrics.%s.mismatched", m.Name()))
	}

	return nil
}
