package mocks

import (
	"time"

	"backend-sample-la/models"
	"github.com/stretchr/testify/mock"
)

// Repository represents repository mock object
type Repository struct {
	mock.Mock
}

// GetDBTimestamp mock method
func (r *Repository) GetDBTimestamp() models.HealthCheck {
	r.Called()
	now := time.Now()
	healthCheck := models.HealthCheck{
		CurrentTimestamp: now,
	}
	return healthCheck
}
