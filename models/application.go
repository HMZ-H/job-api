package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationStatus string

const (
	StatusApplied   ApplicationStatus = "Applied"
	StatusReviewed  ApplicationStatus = "Reviewed"
	StatusInterview ApplicationStatus = "Interview"
	StatusRejected  ApplicationStatus = "Rejected"
	StatusHired     ApplicationStatus = "Hired"
)

type Application struct {
	ID          uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ApplicantID uuid.UUID         `json:"applicant_id" gorm:"type:uuid;not null"`
	JobID       uuid.UUID         `json:"job_id" gorm:"type:uuid;not null"`
	ResumeLink  string            `json:"resume_link" gorm:"not null" validate:"required,url"`
	CoverLetter string            `json:"cover_letter" validate:"max=200"`
	Status      ApplicationStatus `json:"status" gorm:"type:varchar(20);default:'Applied'"`
	AppliedAt   time.Time         `json:"applied_at"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`

	// Relationships
	Applicant User `json:"applicant" gorm:"foreignKey:ApplicantID"`
	Job       Job  `json:"job" gorm:"foreignKey:JobID"`
}

func (a *Application) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.AppliedAt.IsZero() {
		a.AppliedAt = time.Now()
	}
	return nil
}
