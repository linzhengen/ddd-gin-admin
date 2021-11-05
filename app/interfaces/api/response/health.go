package response

import "time"

type HealthCheck struct {
	Status    string    `json:"status"`     // Status
	CheckedAt time.Time `json:"checked_at"` // CheckedAt
}
