package appdriver

import (
	"time"

	"github.com/google/uuid"
)

type DriverProfile struct {
	ID                   uuid.UUID `json:"id"`
	AccountID            uuid.UUID `json:"account_id"`
	FullName             string    `json:"full_name"`
	DateOfBirth          *time.Time `json:"date_of_birth,omitempty"`
	Gender               string    `json:"gender,omitempty"`
	Address              string    `json:"address,omitempty"`
	GlobalStatus         string    `json:"global_status"`
	Rating               float64   `json:"rating"`
	TotalCompletedOrders int32     `json:"total_completed_orders"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty"`
}
