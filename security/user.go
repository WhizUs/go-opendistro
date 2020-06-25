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

type UserService common.Service

type UserServiceInterface interface {
	Get(ctx context.Context, name string) (*User, error)
	List(ctx context.Context) (*[]User, error)
	Create(ctx context.Context, name string, userCreate *UserCreate) error
	ChangePassword(ctx context.Context, name string, newPassword string) error
	common.Modifyable
}

type UserCreate struct {
	Password     string            `json:"password,omitempty"`
	BackendRoles []string          `json:"backend_roles,omitempty"`
	Roles        []string          `json:"opendistro_security_roles"`
	Attributes   map[string]string `json:"attributes,omitempty"`
}

type User struct {
	Name         string
	Hash         string            `json:"hash"`
	Reserved     bool              `json:"reserved"`
	Hidden       bool              `json:"hidden"`
	BackendRoles []string          `json:"backend_roles"`
	Attributes   map[string]string `json:"attributes"`
	Description  string            `json:"description"`
	Static       bool              `json:"static"`
}

// Get a single user by name
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#get-user
func (s *UserService) Get(ctx context.Context, name string) (*User, error) {
	endpoint := common.UsersEndpoint + name

	var users map[string]*User

	err := s.Client.Get(ctx, endpoint, &users)
	if err != nil {
		return nil, err
	}

	if users[name] == nil {
		return nil, nil
	}

	users[name].Name = name

	return users[name], nil
}

// List all users
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#get-users
func (s *UserService) List(ctx context.Context) (*[]User, error) {
	var users map[string]*User

	err := s.Client.Get(ctx, common.UsersEndpoint, &users)
	if err != nil {
		return nil, err
	}

	var _users []User

	for name, user := range users {
		user.Name = name
		_users = append(_users, *user)
	}

	return &_users, nil
}

// Delete a user by name
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#delete-user
func (s *UserService) Delete(ctx context.Context, name string) error {
	endpoint := common.UsersEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodDelete, nil)
}

// Create or replace the specified user. Password can be submitted in plain text (password) or hashed (hash). If a plain text password is submitted, the Security Plugin will do the hashing.
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#create-user
func (s *UserService) Create(ctx context.Context, name string, userCreate *UserCreate) error {
	endpoint := common.UsersEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodPut, userCreate)
}

// Update a user by name and providing an update patch
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#patch-user
func (s *UserService) Update(ctx context.Context, name string, patches *[]common.Patch) error {
	endpoint := common.UsersEndpoint + name

	return s.Client.Modify(ctx, endpoint, http.MethodPatch, patches)
}

// Update multiple users at once
//
// see: https://opendistro.github.io/for-elasticsearch-docs/docs/security-access-control/api/#patch-users
func (s *UserService) UpdateBatch(ctx context.Context, patches *[]common.Patch) error {
	return s.Update(ctx, "", patches)
}

// ChangePassword applies the new password to the user provided by name
func (s *UserService) ChangePassword(ctx context.Context, name string, newPassword string) error {
	patch := &[]common.Patch{
		{
			Op:   "add",
			Path: "/" + name,
			Value: map[string]interface{}{
				"password": newPassword,
			},
		},
	}

	return s.UpdateBatch(ctx, patch)
}
