package model

import (
	"time"

	"tms-core-service/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is the database model for users
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PhoneNumber  string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	FirstName    string
	LastName     string
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// ToEntity converts database model to domain entity
func (m *User) ToEntity() *entity.User {
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return &entity.User{
		ID:           m.ID,
		Email:        m.Email,
		PhoneNumber:  m.PhoneNumber,
		PasswordHash: m.PasswordHash,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}

// FromEntity creates a database model from a domain entity
func FromEntity(e *entity.User) *User {
	var deletedAt gorm.DeletedAt
	if e.DeletedAt != nil {
		deletedAt = gorm.DeletedAt{Time: *e.DeletedAt, Valid: true}
	}

	return &User{
		ID:           e.ID,
		Email:        e.Email,
		PhoneNumber:  e.PhoneNumber,
		PasswordHash: e.PasswordHash,
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}
