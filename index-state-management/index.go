package index_state_management

import (
	"context"
	"encoding/json"
	"github.com/WhizUs/go-opendistro/common"
	"net/http"
	"net/url"
)

type IndexService common.Service

// Adds a policy to an index
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/ism/api/#add-policy
func (s *IndexService) AddPolicy(ctx context.Context, index string, policyID string) (*IndexResponse, error) {
	return s.manipulateIndexPolicyReleation(ctx, "add", index, IndexPolicyChange{PolicyID: policyID})
}

// Remove a policy fraom an index
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/ism/api/#remove-policy-from-index
func (s *IndexService) RemovePolicy(ctx context.Context, index string, policyID string) (*IndexResponse, error) {
	return s.manipulateIndexPolicyReleation(ctx, "remove", index, IndexPolicyChange{PolicyID: policyID})
}

// Update the managed index policy of an index
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/ism/api/#update-managed-index-policy
func (s *IndexService) UpdatePolicy(ctx context.Context, index string, policyID string) (*IndexResponse, error) {
	return s.manipulateIndexPolicyReleation(ctx, "change_policy", index, IndexPolicyChange{PolicyID: policyID})
}

func (s *IndexService) manipulateIndexPolicyReleation(ctx context.Context, changeName string, index string, change IndexPolicyChange) (*IndexResponse, error) {
	endpoint := common.IndexStateManagementEndpoint + changeName + "/" + url.PathEscape(index)

	data, err := s.Client.Do(ctx, change, endpoint, http.MethodPost)
	if err != nil {
		return nil, err
	}
	var resp IndexResponse
	err = json.Unmarshal(data, &resp)
	return &resp, err
}

type IndexServiceInterface interface {
	AddPolicy(ctx context.Context, index string, policyID string) (*IndexResponse, error)
	RemovePolicy(ctx context.Context, index string, policyID string) (*IndexResponse, error)
	UpdatePolicy(ctx context.Context, index string, policyID string) (*IndexResponse, error)
}
