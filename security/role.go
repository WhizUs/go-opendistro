// Copyright 2019 WhizUs GmbH. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package security

import (
	"context"
	"github.com/WhizUs/go-opendistro/common"
	"net/http"
)

type RoleService common.Service

type RoleServiceInterface interface {
	Get(ctx context.Context, name string) (*Role, error)
	List(ctx context.Context) ([]*Role, error)
	Create(ctx context.Context, name string, rolePermissions *RolePermissions) error
	common.Modifyable
}

type Role struct {
	Name        string
	IsStatic    bool   `json:"static"`
	IsReserved  bool   `json:"reserved"`
	IsHidden    bool   `json:"hidden"`
	Description string `json:"description"`
	RolePermissions
}

type RolePermissions struct {
	ClusterPermissions []string             `json:"cluster_permissions,omitempty"`
	IndexPermissions   *[]IndexPermissions  `json:"index_permissions,omitempty"`
	TenantPermissions  *[]TenantPermissions `json:"tenant_permissions,omitempty"`
}

type IndexPermissions struct {
	IndexPatterns  []string `json:"index_patterns,omitempty"`
	Dls            string   `json:"dls,omitempty"`
	Fls            []string `json:"fls,omitempty"`
	MaskedFields   []string `json:"masked_fields,omitempty"`
	AllowedActions []string `json:"allowed_actions,omitempty"`
}

type TenantPermissions struct {
	TenantPatterns []string `json:"tenant_patterns,omitempty"`
	AllowedActions []string `json:"allowed_actions,omitempty"`
}

// Get a single role by name
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#get-role
func (s *RoleService) Get(ctx context.Context, name string) (*Role, error) {
	endpoint := common.RolesEndpoint + name

	var roles map[string]*Role

	err := s.Client.Get(ctx, endpoint, &roles)
	if err != nil {
		return nil, err
	}

	if roles[name] == nil {
		return nil, nil
	}

	roles[name].Name = name

	return roles[name], nil
}

// List all roles
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#get-roles
func (s *RoleService) List(ctx context.Context) ([]*Role, error) {
	var roles map[string]*Role

	err := s.Client.Get(ctx, common.RolesEndpoint, &roles)
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

// Delete a role by name
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#delete-role
func (s *RoleService) Delete(ctx context.Context, name string) error {
	endpoint := common.RolesEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodDelete, nil)
}

// Create a role with permissions
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#create-role
func (s *RoleService) Create(ctx context.Context, name string, rolePermissions *RolePermissions) error {
	endpoint := common.RolesEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodPut, rolePermissions)
}

// Update a role
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#patch-role
func (s *RoleService) Update(ctx context.Context, name string, patches *[]common.Patch) error {
	endpoint := common.RolesEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodPatch, patches)
}

// Update multiple roles at once
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#patch-roles
func (s *RoleService) UpdateBatch(ctx context.Context, patches *[]common.Patch) error {
	return s.Update(ctx, "", patches)
}
