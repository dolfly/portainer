package libprometheus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const configTemplate = `# Prometheus configuration.
# Generated automatically — do not edit.
global:
  scrape_interval: %s
  evaluation_interval: %s

scrape_configs:
  - job_name: '%s'
    static_configs:
      - targets: ['%s']
`

// WritePrometheusConfig writes a Prometheus config.yaml to dataDir.
// jobName identifies the scrape job (e.g. "portainer" or "edge-agent").
// target is the scrape target address.
func WritePrometheusConfig(dataDir string, scrapeInterval string, jobName string, target string) error {
	configDir := filepath.Join(dataDir, "prometheus")
	if err := os.MkdirAll(configDir, 0o750); err != nil {
		return fmt.Errorf("create prometheus config dir: %w", err)
	}

	// Escape single quotes to prevent YAML injection in interpolated values.
	escapeYAMLSingleQuote := func(s string) string { return strings.ReplaceAll(s, "'", "''") }
	content := fmt.Sprintf(configTemplate, scrapeInterval, scrapeInterval, escapeYAMLSingleQuote(jobName), escapeYAMLSingleQuote(target))
	configPath := filepath.Join(configDir, "config.yaml")

	return os.WriteFile(configPath, []byte(content), 0o600)
}
