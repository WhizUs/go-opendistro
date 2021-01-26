package index_state_management

type PolicyResponse struct {
	ID             string `json:"_id"`
	Version        int    `json:"_version"`
	SequenceNumber int    `json:"_seq_no"`
	PrimaryTerm    int    `json:"_primary_term"`
	Policy         struct {
		Policy Policy `json:"policy"`
	} `json:"policy"`
}

type IndexResponse struct {
	UpdatedIndices int      `json:"updated_indices"`
	Failures       bool     `json:"failures"`
	FailedIndices  []string `json:"failed_indices"`
}

type IndexPolicyChange struct {
	PolicyID string `json:"policy_id"`
	State    string `json:"state,omitempty"`
	Include  []struct {
		State string
	} `json:"include,omitempty"`
}

type Policy struct {
	ID                string        `json:"policy_id"`
	Description       string        `json:"description"`
	LastUpdatedTime   int           `json:"last_updated_time"`
	SchemaVersion     int           `json:"schema_version"`
	ErrorNotification *Notification `json:"error_notification"`
	DefaultState      string        `json:"default_state"`
	States            []State       `json:"states"`
}

type Notification struct {
	Destination     string      `json:"destination"`
	MessageTemplate interface{} `json:"message_template"` // TODO: Find concrete fields
}

type State struct {
	Name        string       `json:"name"`
	Actions     []Action     `json:"actions"`
	Transitions []Transition `json:"transitions"`
}

type Action struct {
	Timeout       string              `json:"timeout"`
	Retry         RetryAction         `json:"retry"`
	ForceMerge    ForceAction         `json:"force_merge"`
	ReadOnly      struct{}            `json:"read_only"`
	ReadWrite     struct{}            `json:"read_write"`
	ReplicaCount  ReplicaCountAction  `json:"replica_count"`
	Close         struct{}            `json:"close"`
	Open          struct{}            `json:"open"`
	Delete        struct{}            `json:"delete"`
	Rollover      RolloverAction      `json:"rollover"`
	Notification  Notification        `json:"notification"`
	Snapshot      SnapshotAction      `json:"snapshot"`
	IndexPriority IndexPriorityAction `json:"index_priority"`
	Allocation    AllocationAction    `json:"allocation"`
}

type RetryAction struct {
	Count   int    `json:"count"`
	Backoff string `json:"backoff"`
	Delay   string `json:"delay"`
}

type ForceAction struct {
	MaxNumSegments string `json:"max_num_segments"`
}

type ReplicaCountAction struct {
	NumberOfReplicas int `json:"number_of_replicas"`
}

type RolloverAction struct {
	MinSize     string `json:"min_size"`
	MinDocCount int    `json:"min_doc_count"`
	MinIndexAge string `json:"min_index_age"`
}

type SnapshotAction struct {
	Repository string `json:"repository"`
	Snapshot   string `json:"snapshot"`
}

type IndexPriorityAction struct {
	Priority int `json:"priority"`
}

type AllocationAction struct {
	Require string `json:"require"`
	Include string `json:"include"`
	Exclude string `json:"exclude"`
	WaitFor string `json:"wait_for"`
}

type Transition struct {
	StateName  string               `json:"state_name"`
	Conditions TransitionConditions `json:"conditions"`
}

type TransitionConditions struct {
	MinIndexAge string `json:"min_index_age"`
	MinDocCount int    `json:"min_doc_count"`
	MinSize     string `json:"min_size"`
	Cron        struct {
		Cron struct {
			Expression string `json:"expression"`
			Timezone   string `json:"timezone"`
		} `json:"cron"`
	} `json:"cron"`
}
