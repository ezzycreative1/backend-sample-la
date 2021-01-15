package repository

import (
	HealthCheckInterface "backend-sample-la/app/healthcheck"
	"backend-sample-la/models"
	"github.com/jinzhu/gorm"
)

// HealthCheckRepository represents health check repository
type HealthCheckRepository struct {
	Conn *gorm.DB
}

// NewHealthCheckRepository initialize health check repository
func NewHealthCheckRepository(Conn *gorm.DB) HealthCheckInterface.IHealthCheckRepository {
	return &HealthCheckRepository{Conn}
}

// GetDBTimestamp used to get current timestamp on connected DB
func (m *HealthCheckRepository) GetDBTimestamp() models.HealthCheck {
	var healthCheck models.HealthCheck

	tx := m.Conn.Begin()
	tx.Raw("SELECT current_timestamp").Scan(&healthCheck)
	tx.Commit()

	return healthCheck
}
