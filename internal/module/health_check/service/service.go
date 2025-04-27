package service

import (
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/health_check/ports"
	"github.com/rs/zerolog/log"
)

var _ ports.HeatlhCheckService = &healthCheckService{}

type healthCheckService struct {
}

func NewHealthCheckService() ports.HeatlhCheckService {
	return &healthCheckService{}
}

func (s *healthCheckService) HealthcheckServices() (string, error) {
	log.Info().Msg("service::healthCheckService - Health check service healthy")
	return "service healthy", nil
}
