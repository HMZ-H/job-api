package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RoleApplicant UserRole = "applicant"
	RoleCompany   UserRole = "company"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"not null" validate:"required,alpha"`
	Email     string    `json:"email" gorm:"unique;not null" validate:"required,email"`
	Password  string    `json:"-" gorm:"not null" validate:"required,min=8,containsuppercase,containslowercase,containsdigit,containsspecial"`
	Role      UserRole  `json:"role" gorm:"type:varchar(20);not null" validate:"required,oneof=applicant company"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
