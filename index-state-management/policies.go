package index_state_management

// The API is derived from this documentation: https://opendistro.github.io/for-elasticsearch-docs/docs/ism/api/#get-policy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/WhizUs/go-opendistro/common"
	"net/http"
)

type PolicyServiceInterface interface {
	Get(ctx context.Context, id string) (*GetPolicyResponse, error)
	Update(ctx context.Context, id string, policy Policy) (*PolicyResponse, error)
	Create(ctx context.Context, id string, policy Policy) (*PolicyResponse, error)
	Delete(ctx context.Context, name string) error
}

type PolicyService common.Service

// Get a policy by id
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/ism/api/#get-policy
func (s *PolicyService) Get(ctx context.Context, id string) (*GetPolicyResponse, error) {
	endpoint := common.ISMPoliciesEndpoint + id

	var response *GetPolicyResponse

	err := s.Client.Get(ctx, endpoint, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Update a policy by id
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/ism/api/#update-policycy
func (s *PolicyService) Update(ctx context.Context, id string, policy Policy) (*PolicyResponse, error) {
	old, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	endpoint := common.ISMPoliciesEndpoint + fmt.Sprintf("%s?if_seq_no=%d&if_primary_term=%d", id, old.SequenceNumber, old.PrimaryTerm)
	data, err := s.Client.Do(ctx, policy, endpoint, http.MethodPut)
	if err != nil {
		return nil, err
	}
	var out PolicyResponse
	err = json.Unmarshal(data, &out)
	return &out, err
}

// Create a policy by id
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/ism/api/#create-policy
func (s *PolicyService) Create(ctx context.Context, id string, policy Policy) (*PolicyResponse, error) {
	endpoint := common.ISMPoliciesEndpoint + id
	data, err := s.Client.Do(ctx, struct {
		Policy Policy
	}{
		Policy: policy,
	}, endpoint, http.MethodPut)
	if err != nil {
		return nil, err
	}
	var out PolicyResponse
	err = json.Unmarshal(data, &out)
	return &out, err
}

// Delete a policy by id
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/ism/api/#delete-policy
func (s *PolicyService) Delete(ctx context.Context, name string) error {
	endpoint := common.ISMPoliciesEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodDelete, nil)
}
