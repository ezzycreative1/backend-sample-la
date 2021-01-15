package usecase

import (
	HealthCheckInterface "backend-sample-la/app/healthcheck"
)

// HealthCheckUsecase represents health check use case
type HealthCheckUsecase struct{}

// NewHealthCheckUsecase initialize health check usecase
func NewHealthCheckUsecase() HealthCheckInterface.IHealthCheckUsecase {
	return &HealthCheckUsecase{}
}

// GetDBTimestamp fetch timestamp from repository
// func (a *HealthCheckUsecase) GetDBTimestamp() models.HealthCheck {
// 	healthCheck := a.HealthCheckRepository.GetDBTimestamp()
// 	return healthCheck
// }
