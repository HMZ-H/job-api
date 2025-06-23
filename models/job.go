package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"not null" validate:"required,min=1,max=100"`
	Description string    `json:"description" gorm:"not null" validate:"required,min=20,max=2000"`
	Location    string    `json:"location"`
	CreatedBy   uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Creator      User          `json:"creator" gorm:"foreignKey:CreatedBy"`
	Applications []Application `json:"applications,omitempty" gorm:"foreignKey:JobID"`
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}
