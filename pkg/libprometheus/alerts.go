package libprometheus

import (
	"time"

	"github.com/go-openapi/strfmt"
	alertmanagermodels "github.com/prometheus/alertmanager/api/v2/models"
	promLabels "github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/rules"
)

// BuildPostableAlert converts a Prometheus rules.Alert into an Alertmanager
// PostableAlert. Firing alerts get a future EndsAt so Alertmanager keeps them
// active between resends.
func BuildPostableAlert(a *rules.Alert, scrapeInterval time.Duration) *alertmanagermodels.PostableAlert {
	labels := make(alertmanagermodels.LabelSet, a.Labels.Len())
	a.Labels.Range(func(l promLabels.Label) {
		labels[l.Name] = l.Value
	})

	annotations := make(alertmanagermodels.LabelSet, a.Annotations.Len())
	a.Annotations.Range(func(l promLabels.Label) {
		annotations[l.Name] = l.Value
	})

	pa := &alertmanagermodels.PostableAlert{
		Alert: alertmanagermodels.Alert{
			Labels: labels,
		},
		Annotations: annotations,
		StartsAt:    strfmt.DateTime(a.ActiveAt),
	}

	switch a.State {
	case rules.StateFiring:
		pa.EndsAt = strfmt.DateTime(time.Now().Add(4 * scrapeInterval))
	case rules.StateInactive:
		if !a.ResolvedAt.IsZero() {
			pa.EndsAt = strfmt.DateTime(a.ResolvedAt)
		}
	}

	return pa
}
