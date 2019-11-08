package security

import (
	"context"
	"github.com/WhizUs/go-opendistro/common"
)

type HealthService common.Service

type HealthServiceInterface interface {
	Get(ctx context.Context) (*Health, error)
}

type Health struct {
	Message string `json:"message"`
	Mode    string `json:"mode"`
	Status  string `json:"status"`
}

func (s *HealthService) Get(ctx context.Context) (*Health, error) {
	endpoint := common.HealthEndpoint

	var health *Health

	err := s.Client.Get(ctx, endpoint, &health)
	if err != nil {
		return nil, err
	}

	return health, nil
}
