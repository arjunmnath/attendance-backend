package models

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    int        `gorm:"not null;index"`
	DeviceID  int        `gorm:"not null;unique;index"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	Attending Attendance `gorm:"foreignKey:DeviceID;references:DeviceID"`
}

type CurrentEvents struct {
	ID        uuid.UUID    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	StartTime time.Time    `gorm:"not null"`
	EndTime   time.Time    `gorm:"not null"`
	Location  string       `gorm:"not null"`
	Attendees []Attendance `gorm:"foreignKey:EventID"`
}

type Attendance struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	DeviceID       int       `gorm:"index;not-null"`
	EventID        uuid.UUID `gorm:"not-null"`
	ProximityScore int       `gorm:"default:0"`
	PollCount      int       `gorm:"default:0"`
}
