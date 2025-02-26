package models

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    int       `gorm:"not null;index"`
	DeviceID  int       `gorm:"not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
