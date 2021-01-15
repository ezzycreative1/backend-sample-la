package mocks

import (
	"time"

	"backend-sample-la/models"
	"github.com/stretchr/testify/mock"
)

// Usecase represents usecase mock object
type Usecase struct {
	mock.Mock
}

// GetDBTimestamp mock method
func (r *Usecase) GetDBTimestamp() models.HealthCheck {
	r.Called()
	now := time.Now()
	healthCheck := models.HealthCheck{
		CurrentTimestamp: now,
	}
	return healthCheck
}
