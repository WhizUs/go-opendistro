package opendistro

import (
    "context"
    "net/http"
)

type TenantService service

type TenantServiceInterface interface {
	Get(ctx context.Context, name string) (*Tenant, error)
	List(ctx context.Context) ([]*Tenant, error)
	Create(ctx context.Context, name string) (*StatusResponse, error)
	Delete(ctx context.Context, name string) (*StatusResponse, error)
	Update(ctx context.Context, name string, patches []*Patch) (*StatusResponse, error)
	UpdateBatch(ctx context.Context, patches []*Patch) (*StatusResponse, error)
}

type Tenant struct{
    Name string
}

func (s *TenantService) Get(ctx context.Context, name string) (*Tenant, error) {
	endpoint := tenantEndpoint + name

    var tenants map[string]*Tenant

    err := s.client.get(ctx, endpoint, &tenants)
    if err != nil {
        return nil, err
    }

    if tenants[name] == nil {
        return nil, nil
    }

    tenants[name].Name = name

    return tenants[name], nil
}

func (s *TenantService) List(ctx context.Context) ([]*Tenant, error) {
    var tenants map[string]*Tenant

    err := s.client.get(ctx, tenantEndpoint, &tenants)
    if err != nil {
        return nil, err
    }

    var _tenants []*Tenant

    for name, tenant := range tenants {
        tenant.Name = name
        _tenants = append(_tenants, tenant)
    }

    return _tenants, nil
}

func (s *TenantService) Delete(ctx context.Context, name string) (*StatusResponse, error) {
    endpoint := usersEndpoint + name

    return s.client.modify(ctx, endpoint, http.MethodDelete, nil)
}

//@todo
func (s *TenantService) Create(ctx context.Context, name string) (*StatusResponse, error) {
    endpoint := usersEndpoint + name

    return s.client.modify(ctx, endpoint, http.MethodPut, nil)
}

func (s *TenantService) Update(ctx context.Context, name string, patches []*Patch) (*StatusResponse, error) {
    endpoint := usersEndpoint + name

    return s.client.modify(ctx, endpoint, http.MethodPatch, &patches)
}

func (s *TenantService) UpdateBatch(ctx context.Context, patches []*Patch) (*StatusResponse, error) {
    return s.Update(ctx, "", patches)
}
