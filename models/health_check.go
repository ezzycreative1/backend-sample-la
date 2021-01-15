package models

import (
	"time"
)

// HealthCheck database model
type HealthCheck struct {
	CurrentTimestamp time.Time `json:"current_timestamp"`
}
