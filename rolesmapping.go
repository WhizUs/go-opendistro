package opendistro

import (
	"context"
	"net/http"
)

type RolesmappingService service

type RolesmappingServiceInterface interface {
	Get(ctx context.Context, name string) (*RoleMapping, error)
	List(ctx context.Context) ([]*RoleMapping, error)
	Create(ctx context.Context, name string) (*StatusResponse, error)
	Delete(ctx context.Context, name string) (*StatusResponse, error)
	Update(ctx context.Context, name string, patches []*Patch) (*StatusResponse, error)
	UpdateBatch(ctx context.Context, patches []*Patch) (*StatusResponse, error)
}

type RoleMapping struct {
	Name        string
	IsReserved  bool   `json:"reserved"`
	IsHidden    bool   `json:"hidden"`
	Description string `json:"description"`
	RoleMappingRelations
}

type RoleMappingRelations struct {
	BackendRoles []string `json:"backend_roles,omitempty"`
	Hosts        []string `json:"hosts,omitempty"`
	Users        []string `json:"users,omitempty"`
}

func (s *RolesmappingService) Get(ctx context.Context, name string) (*RoleMapping, error) {
	endpoint := rolesMappingEndpoint + name

	var rolemappings map[string]*RoleMapping

	err := s.client.get(ctx, endpoint, &rolemappings)
	if err != nil {
		return nil, err
	}

	if rolemappings[name] == nil {
		return nil, nil
	}

	rolemappings[name].Name = name

	return rolemappings[name], nil
}

func (s *RolesmappingService) List(ctx context.Context) ([]*RoleMapping, error) {
	var rolemappings map[string]*RoleMapping

	err := s.client.get(ctx, rolesMappingEndpoint, &rolemappings)
	if err != nil {
		return nil, err
	}

	var _rolemappings []*RoleMapping

	for name, rolemapping := range rolemappings {
		rolemapping.Name = name
		_rolemappings = append(_rolemappings, rolemapping)
	}

	return _rolemappings, nil
}

func (s *RolesmappingService) Delete(ctx context.Context, name string) (*StatusResponse, error) {
	endpoint := rolesMappingEndpoint + name

	return s.client.modify(ctx, endpoint, http.MethodDelete, nil)
}

func (s *RolesmappingService) Create(ctx context.Context, name string, roleMappingRelations *RoleMappingRelations) (*StatusResponse, error) {
	endpoint := rolesMappingEndpoint + name

	return s.client.modify(ctx, endpoint, http.MethodPut, &roleMappingRelations)
}

func (s *RolesmappingService) Update(ctx context.Context, name string, patches []*Patch) (*StatusResponse, error) {
	endpoint := rolesMappingEndpoint + name

	return s.client.modify(ctx, endpoint, http.MethodPatch, &patches)
}

func (s *RolesmappingService) UpdateBatch(ctx context.Context, patches []*Patch) (*StatusResponse, error) {
	return s.Update(ctx, "", patches)
}
