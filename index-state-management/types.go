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

type GetPolicyResponse struct {
	ID             string `json:"_id"`
	Version        int    `json:"_version"`
	SequenceNumber int    `json:"_seq_no"`
	PrimaryTerm    int    `json:"_primary_term"`
	Policy         Policy `json:"policy"`
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
	ID                string        `json:"policy_id,omitempty"`
	Description       string        `json:"description,omitempty"`
	LastUpdatedTime   int           `json:"last_updated_time,omitempty"`
	SchemaVersion     int           `json:"schema_version,omitempty"`
	ErrorNotification *Notification `json:"error_notification,omitempty"`
	DefaultState      string        `json:"default_state,omitempty"`
	States            []State       `json:"states,omitempty"`
}

type Notification struct {
	Destination     string      `json:"destination,omitempty"`
	MessageTemplate interface{} `json:"message_template,omitempty"` // TODO: Find concrete fields
}

type State struct {
	Name        string       `json:"name"`
	Actions     []Action     `json:"actions,omitempty"`
	Transitions []Transition `json:"transitions,omitempty"`
}

type Action struct {
	Timeout       string               `json:"timeout,omitempty"`
	Retry         *RetryAction         `json:"retry,omitempty"`
	ForceMerge    *ForceAction         `json:"force_merge,omitempty"`
	ReadOnly      *struct{}            `json:"read_only,omitempty"`
	ReadWrite     *struct{}            `json:"read_write,omitempty"`
	ReplicaCount  *ReplicaCountAction  `json:"replica_count,omitempty"`
	Close         *struct{}            `json:"close,omitempty"`
	Open          *struct{}            `json:"open,omitempty"`
	Delete        *struct{}            `json:"delete,omitempty"`
	Rollover      *RolloverAction      `json:"rollover,omitempty"`
	Notification  *Notification        `json:"notification,omitempty"`
	Snapshot      *SnapshotAction      `json:"snapshot,omitempty"`
	IndexPriority *IndexPriorityAction `json:"index_priority,omitempty"`
	Allocation    *AllocationAction    `json:"allocation,omitempty"`
}

type RetryAction struct {
	Count   int    `json:"count,omitempty"`
	Backoff string `json:"backoff,omitempty"`
	Delay   string `json:"delay,omitempty"`
}

type ForceAction struct {
	MaxNumSegments int `json:"max_num_segments,omitempty"`
}

type ReplicaCountAction struct {
	NumberOfReplicas int `json:"number_of_replicas,omitempty"`
}

type RolloverAction struct {
	MinSize     string `json:"min_size,omitempty"`
	MinDocCount int    `json:"min_doc_count,omitempty"`
	MinIndexAge string `json:"min_index_age,omitempty"`
}

type SnapshotAction struct {
	Repository string `json:"repository,omitempty"`
	Snapshot   string `json:"snapshot,omitempty"`
}

type IndexPriorityAction struct {
	Priority int `json:"priority,omitempty"`
}

type AllocationAction struct {
	Require string `json:"require,omitempty"`
	Include string `json:"include,omitempty"`
	Exclude string `json:"exclude,omitempty"`
	WaitFor string `json:"wait_for,omitempty"`
}

type Transition struct {
	StateName  string               `json:"state_name"`
	Conditions TransitionConditions `json:"conditions,omitempty"`
}

type TransitionConditions struct {
	MinIndexAge string `json:"min_index_age,omitempty"`
	MinDocCount int    `json:"min_doc_count,omitempty"`
	MinSize     string `json:"min_size,omitempty"`
	Cron        *struct {
		Cron struct {
			Expression string `json:"expression,omitempty"`
			Timezone   string `json:"timezone,omitempty"`
		} `json:"cron,omitempty"`
	} `json:"cron,omitempty"`
}
