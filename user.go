package opendistro

import (
    "context"
    "fmt"
    "net/http"
)

type UserService service

type UserServiceInterface interface {
	Get(ctx context.Context, name string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Create(ctx context.Context, name string, userCreate *UserCreate) (*StatusResponse, error)
	Delete(ctx context.Context, name string) (*StatusResponse, error)
	Update(ctx context.Context, name string, patches []*Patch) (*StatusResponse, error)
	UpdateBatch(ctx context.Context, patches []*Patch) (*StatusResponse, error)
}

type UserCreate struct {
	Password     string   `json:"password"` // Passwords must be at least 6 characters long.
	BackendRoles []string `json:"backend_roles,omitempty"`
	Attributes   []string `json:"attributes,omitempty"`
}

type User struct {
	Name         string
	Hash         string      `json:"hash"`
	Reserved     bool        `json:"reserved"`
	Hidden       bool        `json:"hidden"`
	BackendRoles []string    `json:"backend_roles"`
	Attributes   interface{} `json:"attributes"`
	Description  string      `json:"description"`
	Static       bool        `json:"static"`
}

func (s *UserService) Get(ctx context.Context, name string) (*User, error) {
	endpoint := usersEndpoint + name

    var users map[string]*User

    err := s.client.get(ctx, endpoint, &users)
    if err != nil {
        return nil, err
    }

    if users[name] == nil {
        return nil, fmt.Errorf("get user: %s not in response", name)
    }

    users[name].Name = name

    return users[name], nil
}

func (s *UserService) List(ctx context.Context) ([]*User, error) {
    var users map[string]*User

    err := s.client.get(ctx, usersEndpoint, &users)
    if err != nil {
        return nil, err
    }

    var _users []*User

    for name, user := range users {
        user.Name = name
        _users = append(_users, user)
    }

    return _users, nil
}

func (s *UserService) Delete(ctx context.Context, name string) (*StatusResponse, error) {
    endpoint := usersEndpoint + name

    return s.client.modify(ctx, endpoint, http.MethodDelete, nil)
}

func (s *UserService) Create(ctx context.Context, name string, userCreate *UserCreate) (*StatusResponse, error) {
    endpoint := usersEndpoint + name

    return s.client.modify(ctx, endpoint, http.MethodPut, &userCreate)
}

func (s *UserService) Update(ctx context.Context, name string, patches []*Patch) (*StatusResponse, error) {
    endpoint := usersEndpoint + name

    return s.client.modify(ctx, endpoint, http.MethodPatch, &patches)
}

func (s *UserService) UpdateBatch(ctx context.Context, patches []*Patch) (*StatusResponse, error) {
    return s.Update(ctx, "", patches)
}
