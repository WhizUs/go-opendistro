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

type TenantService common.Service

type TenantServiceInterface interface {
	Get(ctx context.Context, name string) (*Tenant, error)
	List(ctx context.Context) (*[]Tenant, error)
	Create(ctx context.Context, name string) error
	common.Modifyable
}

type Tenant struct {
	Name        string `json:"name"`
	Reserved    bool   `json:"reserved"`
	Hidden      bool   `json:"hidden"`
	Description string `json:"description"`
	Static      bool   `json:"static"`
}

// Get a single tenant by name
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#get-tenant
func (s *TenantService) Get(ctx context.Context, name string) (*Tenant, error) {
	endpoint := common.TenantEndpoint + name

	var tenants map[string]*Tenant

	err := s.Client.Get(ctx, endpoint, &tenants)
	if err != nil {
		return nil, err
	}

	if tenants[name] == nil {
		return nil, nil
	}

	tenants[name].Name = name

	return tenants[name], nil
}

// List all tenants
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#get-tenants
func (s *TenantService) List(ctx context.Context) (*[]Tenant, error) {
	var tenants map[string]*Tenant

	err := s.Client.Get(ctx, common.TenantEndpoint, &tenants)
	if err != nil {
		return nil, err
	}

	var _tenants []Tenant

	for name, tenant := range tenants {
		tenant.Name = name
		_tenants = append(_tenants, *tenant)
	}

	return &_tenants, nil
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#delete-tenant
func (s *TenantService) Delete(ctx context.Context, name string) error {
	endpoint := common.TenantEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodDelete, nil)
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#create-tenant
func (s *TenantService) Create(ctx context.Context, name string) error {
	endpoint := common.TenantEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodPut, nil)
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#patch-tenant
func (s *TenantService) Update(ctx context.Context, name string, patches *[]common.Patch) error {
	endpoint := common.TenantEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodPatch, patches)
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#patch-tenants
func (s *TenantService) UpdateBatch(ctx context.Context, patches *[]common.Patch) error {
	return s.Update(ctx, "", patches)
}
