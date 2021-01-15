package healthcheck

import "backend-sample-la/models"

// IHealthCheckRepository represents health check repository interface
type IHealthCheckRepository interface {
	GetDBTimestamp() models.HealthCheck
}
