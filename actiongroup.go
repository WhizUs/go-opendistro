package opendistro

import (
	"context"
	"net/http"
)

type ActiongroupService service

type ActiongroupServiceInterface interface {
	Get(ctx context.Context, name string) (*Actiongroup, error)
	List(ctx context.Context) ([]*Actiongroup, error)
	Create(ctx context.Context, name string) error
	Delete(ctx context.Context, name string) error
	Update(ctx context.Context, name string, patches []*Patch) error
	UpdateBatch(ctx context.Context, patches []*Patch) error
}

type Actiongroup struct {
	Name string
}

func (s *ActiongroupService) Get(ctx context.Context, name string) (*Actiongroup, error) {
	endpoint := actiongroupEndpoint + name

	var actiongroups map[string]*Actiongroup

	err := s.client.get(ctx, endpoint, &actiongroups)
	if err != nil {
		return nil, err
	}

	if actiongroups[name] == nil {
		return nil, nil
	}

	actiongroups[name].Name = name

	return actiongroups[name], nil
}

func (s *ActiongroupService) List(ctx context.Context) ([]*Actiongroup, error) {
	var actiongroups map[string]*Actiongroup

	err := s.client.get(ctx, actiongroupEndpoint, &actiongroups)
	if err != nil {
		return nil, err
	}

	var _actiongroups []*Actiongroup

	for name, actiongroup := range actiongroups {
		actiongroup.Name = name
		_actiongroups = append(_actiongroups, actiongroup)
	}

	return _actiongroups, nil
}

func (s *ActiongroupService) Delete(ctx context.Context, name string) error {
	endpoint := actiongroupEndpoint + name

	return s.client.modify(ctx, endpoint, http.MethodDelete, nil)
}

//@todo
func (s *ActiongroupService) Create(ctx context.Context, name string) error {
	endpoint := actiongroupEndpoint + name

	return s.client.modify(ctx, endpoint, http.MethodPut, nil)
}

func (s *ActiongroupService) Update(ctx context.Context, name string, patches []*Patch) error {
	endpoint := actiongroupEndpoint + name

	return s.client.modify(ctx, endpoint, http.MethodPatch, &patches)
}

func (s *ActiongroupService) UpdateBatch(ctx context.Context, patches []*Patch) error {
	return s.Update(ctx, "", patches)
}
