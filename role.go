package opendistro

import (
	"context"
	"net/http"
)

type RoleService service

type RoleServiceInterface interface {
	Get(ctx context.Context, name string) (*Role, error)
	List(ctx context.Context) ([]*Role, error)
	Create(ctx context.Context, name string, rolePermissions *RolePermissions) (*StatusResponse, error)
	Delete(ctx context.Context, name string) (*StatusResponse, error)
	Update(ctx context.Context, name string, patches []*Patch) (*StatusResponse, error)
	UpdateBatch(ctx context.Context, patches []*Patch) (*StatusResponse, error)
}

type Role struct {
	Name        string
	IsStatic    bool   `json:"static"`
	IsReserved  bool   `json:"reserved"`
	IsHidden    bool   `json:"hidden"`
	Description string `json:"description"`
	RolePermissions
}

type _Role struct {
	IsStatic    bool   `json:"static"`
	IsReserved  bool   `json:"reserved"`
	IsHidden    bool   `json:"hidden"`
	Description string `json:"description"`
	RolePermissions
}

type RolePermissions struct {
	ClusterPermissions []string             `json:"cluster_permissions,omitempty"`
	IndexPermissions   []*IndexPermissions  `json:"index_permissions,omitempty"`
	TenantPermissions  []*TenantPermissions `json:"tenant_permissions,omitempty"`
}

type IndexPermissions struct {
	IndexPatterns  []string `json:"index_patterns,omitempty"`
	Dls            []string `json:"dls,omitempty"`
	Fls            []string `json:"fls,omitempty"`
	MaskedFields   []string `json:"masked_fields,omitempty"`
	AllowedActions []string `json:"allowed_actions,omitempty"`
}

type TenantPermissions struct {
	TenantPatterns []string `json:"tenant_patterns,omitempty"`
	AllowedActions []string `json:"allowed_actions,omitempty"`
}

func (s *RoleService) Get(ctx context.Context, name string) (*Role, error) {
	endpoint := rolesEndpoint + name

	var roles map[string]*Role

	err := s.client.get(ctx, endpoint, &roles)
	if err != nil {
		return nil, err
	}

	if roles[name] == nil {
		return nil, nil
	}

	roles[name].Name = name

	return roles[name], nil
}

func (s *RoleService) List(ctx context.Context) ([]*Role, error) {
	var roles map[string]*Role

	err := s.client.get(ctx, rolesEndpoint, &roles)
	if err != nil {
		return nil, err
	}

	var _roles []*Role

	for name, role := range roles {
		role.Name = name
		_roles = append(_roles, role)
	}

	return _roles, nil
}

func (s *RoleService) Delete(ctx context.Context, name string) (*StatusResponse, error) {
	endpoint := rolesEndpoint + name

	return s.client.modify(ctx, endpoint, http.MethodDelete, nil)
}

func (s *RoleService) Create(ctx context.Context, name string, rolePermissions *RolePermissions) (*StatusResponse, error) {
	endpoint := rolesEndpoint + name

	return s.client.modify(ctx, endpoint, http.MethodPut, &rolePermissions)
}

func (s *RoleService) Update(ctx context.Context, name string, patches []*Patch) (*StatusResponse, error) {
	endpoint := rolesEndpoint + name

	return s.client.modify(ctx, endpoint, http.MethodPatch, &patches)
}

func (s *RoleService) UpdateBatch(ctx context.Context, patches []*Patch) (*StatusResponse, error) {
	return s.Update(ctx, "", patches)
}
