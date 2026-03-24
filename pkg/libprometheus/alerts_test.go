package libprometheus_test

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	libprom "github.com/portainer/portainer/pkg/libprometheus"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildPostableAlertFiringAlert(t *testing.T) {
	activeAt := time.Now().Add(-5 * time.Minute)
	alert := &rules.Alert{
		State: rules.StateFiring,
		Labels: labels.FromStrings(
			"alertname", "HighCPU",
			"alert_rule_id", "42",
			"severity", "critical",
		),
		Annotations: labels.FromStrings(
			"summary", "CPU usage is high",
			"description", "CPU usage is above 90%",
		),
		ActiveAt: activeAt,
	}

	interval := time.Minute
	before := time.Now()
	pa := libprom.BuildPostableAlert(alert, interval)
	after := time.Now()

	require.NotNil(t, pa)
	assert.Equal(t, "HighCPU", pa.Labels["alertname"])
	assert.Equal(t, "42", pa.Labels["alert_rule_id"])
	assert.Equal(t, "critical", pa.Labels["severity"])
	assert.Equal(t, "CPU usage is high", pa.Annotations["summary"])
	assert.Equal(t, "CPU usage is above 90%", pa.Annotations["description"])
	assert.Equal(t, strfmt.DateTime(activeAt), pa.StartsAt)

	endsAt := time.Time(pa.EndsAt)
	assert.False(t, endsAt.Before(before.Add(4*interval)))
	assert.False(t, endsAt.After(after.Add(4*interval)))
}

func TestBuildPostableAlertResolvedAlert(t *testing.T) {
	activeAt := time.Now().Add(-5 * time.Minute)
	resolvedAt := time.Now().Add(-30 * time.Second)
	alert := &rules.Alert{
		State: rules.StateInactive,
		Labels: labels.FromStrings(
			"alertname", "HighCPU",
			"alert_rule_id", "42",
		),
		Annotations: labels.FromStrings(
			"summary", "CPU usage recovered",
		),
		ActiveAt:   activeAt,
		ResolvedAt: resolvedAt,
	}

	pa := libprom.BuildPostableAlert(alert, time.Minute)

	require.NotNil(t, pa)
	assert.Equal(t, "HighCPU", pa.Labels["alertname"])
	assert.Equal(t, "42", pa.Labels["alert_rule_id"])
	assert.Equal(t, "CPU usage recovered", pa.Annotations["summary"])
	assert.Equal(t, strfmt.DateTime(activeAt), pa.StartsAt)
	assert.Equal(t, strfmt.DateTime(resolvedAt), pa.EndsAt)
}
