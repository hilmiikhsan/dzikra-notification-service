package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-notification-service/internal/module/health_check/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-notification-service/internal/module/health_check/service"
	"github.com/Digitalkeun-Creative/be-dzikra-notification-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-notification-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type healthCheckHandler struct {
	service ports.HeatlhCheckService
}

func NewHealthCheckHandler() *healthCheckHandler {
	var handler = new(healthCheckHandler)

	// service
	healthCheckService := service.NewHealthCheckService()

	// handler
	handler.service = healthCheckService

	return handler
}

func (h *healthCheckHandler) HealthCheckRoute(router fiber.Router) {
	router.Get("/health-check", h.healthCheck)
}

func (h *healthCheckHandler) healthCheck(c *fiber.Ctx) error {
	msg, err := h.service.HealthcheckServices()
	if err != nil {
		log.Error().Err(err).Msg("handler::healthCheck - Failed to get health check")
		_, errs := err_msg.Errors[error](err)
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, msg))
}
