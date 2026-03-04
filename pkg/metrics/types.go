package metrics

// EdgeAlertBatch is the generic envelope pushed by edge agents.
type EdgeAlertBatch struct {
	EnvironmentName string       `json:"environment_name"`
	Timestamp       int64        `json:"timestamp,omitempty"` // unix nanoseconds
	FiredAlerts     []FiredAlert `json:"fired_alerts,omitempty"`
	RawSignals      *RawSignals  `json:"raw_signals,omitempty"`
}

// FiredAlert represents a rule the agent has already evaluated and determined should fire.
type FiredAlert struct {
	RuleID   int    `json:"rule_id"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// RawSignals carries raw K8s cluster health data for future server-side evaluation.
type RawSignals struct {
	Workloads  []WorkloadItem  `json:"workloads,omitempty"`
	Conditions []ConditionItem `json:"conditions,omitempty"`
}

// ResourceRef identifies a Kubernetes resource.
type ResourceRef struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// ConditionItem holds a single Kubernetes condition observation.
type ConditionItem struct {
	Ref      ResourceRef `json:"ref"`
	Type     string      `json:"type"`
	Status   int8        `json:"status"` // -1 unknown, 0 unhealthy, 1 healthy
	Reason   string      `json:"reason"`
	NodeName string      `json:"node_name,omitempty"`
	PodName  string      `json:"pod_name,omitempty"`
}

// WorkloadItem holds a single Kubernetes workload replica observation.
type WorkloadItem struct {
	Ref       ResourceRef `json:"ref"`
	Desired   int         `json:"desired"`
	Available int         `json:"available"`
}
