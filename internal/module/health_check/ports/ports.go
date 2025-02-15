package ports

type HeatlhCheckService interface {
	HealthcheckServices() (string, error)
}
