package model

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Device struct {
	ID         uuid.UUID       `json:"id"`
	DeviceUID  string          `json:"device_uid"`
	Platform   string          `json:"platform"`
	DeviceName string          `json:"device_name"`
	OsVersion  string          `json:"os_version"`
	AppVersion string          `json:"app_version"`
	Metadata   json.RawMessage `json:"metadata"`
	BaseModel
}
