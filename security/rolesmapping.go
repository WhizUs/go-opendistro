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

type RolesmappingService common.Service

type RolesmappingServiceInterface interface {
	Get(ctx context.Context, name string) (*RoleMapping, error)
	List(ctx context.Context) (*[]RoleMapping, error)
	Create(ctx context.Context, name string, roleMappingRelations *RoleMappingRelations) error
	common.Modifyable
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

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#
func (s *RolesmappingService) Get(ctx context.Context, name string) (*RoleMapping, error) {
	endpoint := common.RolesMappingEndpoint + name

	var rolemappings map[string]*RoleMapping

	err := s.Client.Get(ctx, endpoint, &rolemappings)
	if err != nil {
		return nil, err
	}

	if rolemappings[name] == nil {
		return nil, nil
	}

	rolemappings[name].Name = name

	return rolemappings[name], nil
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#
func (s *RolesmappingService) List(ctx context.Context) (*[]RoleMapping, error) {
	var rolemappings map[string]*RoleMapping

	err := s.Client.Get(ctx, common.RolesMappingEndpoint, &rolemappings)
	if err != nil {
		return nil, err
	}

	var _rolemappings []RoleMapping

	for name, rolemapping := range rolemappings {
		rolemapping.Name = name
		_rolemappings = append(_rolemappings, *rolemapping)
	}

	return &_rolemappings, nil
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#
func (s *RolesmappingService) Delete(ctx context.Context, name string) error {
	endpoint := common.RolesMappingEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodDelete, nil)
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#
func (s *RolesmappingService) Create(ctx context.Context, name string, roleMappingRelations *RoleMappingRelations) error {
	endpoint := common.RolesMappingEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodPut, roleMappingRelations)
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#
func (s *RolesmappingService) Update(ctx context.Context, name string, patches *[]common.Patch) error {
	endpoint := common.RolesMappingEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodPatch, patches)
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#
func (s *RolesmappingService) UpdateBatch(ctx context.Context, patches *[]common.Patch) error {
	return s.Update(ctx, "", patches)
}
