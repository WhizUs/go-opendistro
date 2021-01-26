package index_state_management

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParsePolicyResponse(t *testing.T) {
	tcs := []struct {
		name   string
		input  string
		output PolicyResponse
	}{
		{
			name: "from docs",
			input: `{
  "_id": "policy_1",
  "_version": 1,
  "_primary_term": 1,
  "_seq_no": 7,
  "policy": {
    "policy": {
      "policy_id": "policy_1",
      "description": "ingesting logs",
      "last_updated_time": 1577990761311,
      "schema_version": 1,
      "error_notification": null,
      "default_state": "ingest",
      "states": [
        {
          "name": "ingest",
          "actions": [
            {
              "rollover": {
                "min_doc_count": 5
              }
            }
          ],
          "transitions": [
            {
              "state_name": "search"
            }
          ]
        },
        {
          "name": "search",
          "actions": [],
          "transitions": [
            {
              "state_name": "delete",
              "conditions": {
                "min_index_age": "5m"
              }
            }
          ]
        },
        {
          "name": "delete",
          "actions": [
            {
              "delete": {}
            }
          ],
          "transitions": []
        }
      ]
    }
  }
}
`,
			output: PolicyResponse{
				ID:             "policy_1",
				Version:        1,
				PrimaryTerm:    1,
				SequenceNumber: 7,
				Policy: struct {
					Policy Policy `json:"policy"`
				}{Policy: Policy{
					ID:                "policy_1",
					Description:       "ingesting logs",
					LastUpdatedTime:   1577990761311,
					SchemaVersion:     1,
					ErrorNotification: nil,
					DefaultState:      "ingest",
					States: []State{
						{
							Name: "ingest",
							Actions: []Action{
								{
									Rollover: RolloverAction{MinDocCount: 5},
								},
							},
							Transitions: []Transition{
								{
									StateName: "search",
								},
							},
						},
						{
							Name:    "search",
							Actions: []Action{},
							Transitions: []Transition{
								{
									StateName:  "delete",
									Conditions: TransitionConditions{MinIndexAge: "5m"},
								},
							},
						},
						{
							Name: "delete",
							Actions: []Action{
								{
									Delete: struct{}{},
								},
							},
							Transitions: []Transition{},
						},
					},
				}},
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var out PolicyResponse
			err := json.Unmarshal([]byte(tc.input), &out)
			if err != nil {
				t.Fail()
			}
			if diff := cmp.Diff(tc.output, out); diff != "" {
				t.Errorf("ParseJSON() mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
