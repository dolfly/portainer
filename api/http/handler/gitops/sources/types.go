package sources

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path"
	"strings"

	gittypes "github.com/portainer/portainer/api/git/types"
	"github.com/portainer/portainer/api/gitops/workflows"
)

// Source represents a unique git repository used as a GitOps source across one or more workflows.
type Source struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Type         string           `json:"type"`
	URL          string           `json:"url"`
	Status       workflows.Status `json:"status"`
	Error        string           `json:"error,omitempty"`
	UsedBy       int              `json:"usedBy"`
	Environments int              `json:"environments"`
	LastSync     int64            `json:"lastSync"`
}

type SourceType string

const (
	SourceTypeGit  SourceType = "git"
	SourceTypeHelm SourceType = "helm"
	SourceTypeOCI  SourceType = "oci"
)

func parseSourceType(s string) (string, error) {
	switch SourceType(s) {
	case SourceTypeGit, SourceTypeHelm, SourceTypeOCI:
		return s, nil
	default:
		return "", fmt.Errorf("invalid source type %q: must be git, helm, or oci", s)
	}
}

type sourceGroupKey struct {
	URL      string
	Username string
	Password string
}

func gitSourceKey(cfg *gittypes.RepoConfig) sourceGroupKey {
	key := sourceGroupKey{URL: cfg.URL}
	if cfg.Authentication != nil {
		key.Username = cfg.Authentication.Username
		key.Password = cfg.Authentication.Password
	}
	return key
}

func sourceID(key sourceGroupKey) string {
	h := sha256.Sum256([]byte(key.URL + "\x00" + key.Username + "\x00" + key.Password))
	return hex.EncodeToString(h[:8])
}

// repoName extracts the repository name from a URL.
// e.g. "https://github.com/org/app-config.git" → "app-config"
func repoName(rawURL string) string {
	base := path.Base(rawURL)
	return strings.TrimSuffix(base, ".git")
}

func worstCaseStatus(statuses []workflows.Status) workflows.Status {
	priority := map[workflows.Status]int{
		workflows.StatusError:   4,
		workflows.StatusSyncing: 3,
		workflows.StatusPaused:  2,
		workflows.StatusHealthy: 1,
		workflows.StatusUnknown: 0,
	}
	worst := workflows.StatusUnknown
	for _, s := range statuses {
		if priority[s] > priority[worst] {
			worst = s
		}
	}
	return worst
}
