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

type ActiongroupService common.Service

type ActiongroupServiceInterface interface {
	Get(ctx context.Context, name string) (*Actiongroup, error)
	List(ctx context.Context) (*[]Actiongroup, error)
	Create(ctx context.Context, name string) error
	common.Modifyable
}

type Actiongroup struct {
	Name           string   `json:"name"`
	Reserved       bool     `json:"reserved"`
	Hidden         bool     `json:"hidden"`
	AllowedActions []string `json:"allowed_actions"`
	Type           string   `json:"type"`
	Description    string   `json:"description"`
	Static         bool     `json:"static"`
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#
func (s *ActiongroupService) Get(ctx context.Context, name string) (*Actiongroup, error) {
	endpoint := common.ActiongroupEndpoint + name

	var actiongroups map[string]*Actiongroup

	err := s.Client.Get(ctx, endpoint, &actiongroups)
	if err != nil {
		return nil, err
	}

	if actiongroups[name] == nil {
		return nil, nil
	}

	actiongroups[name].Name = name

	return actiongroups[name], nil
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#
func (s *ActiongroupService) List(ctx context.Context) (*[]Actiongroup, error) {
	var actiongroups map[string]*Actiongroup

	err := s.Client.Get(ctx, common.ActiongroupEndpoint, &actiongroups)
	if err != nil {
		return nil, err
	}

	var _actiongroups []Actiongroup

	for name, actiongroup := range actiongroups {
		actiongroup.Name = name
		_actiongroups = append(_actiongroups, *actiongroup)
	}

	return &_actiongroups, nil
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#
func (s *ActiongroupService) Delete(ctx context.Context, name string) error {
	endpoint := common.ActiongroupEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodDelete, nil)
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#
func (s *ActiongroupService) Create(ctx context.Context, name string) error {
	endpoint := common.ActiongroupEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodPut, nil)
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#
func (s *ActiongroupService) Update(ctx context.Context, name string, patches *[]common.Patch) error {
	endpoint := common.ActiongroupEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodPatch, patches)
}

//
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#
func (s *ActiongroupService) UpdateBatch(ctx context.Context, patches *[]common.Patch) error {
	return s.Update(ctx, "", patches)
}
