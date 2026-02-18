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
	Email        *string   `gorm:"uniqueIndex"`
	PhoneNumber  *string   `gorm:"uniqueIndex"`
	PasswordHash string
	FirstName    string
	LastName     string
	AvatarURL    string
	GoogleID     *string   `gorm:"uniqueIndex"`
	LineID       *string   `gorm:"uniqueIndex"`
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
		AvatarURL:    m.AvatarURL,
		GoogleID:     m.GoogleID,
		LineID:       m.LineID,
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
		AvatarURL:    e.AvatarURL,
		GoogleID:     e.GoogleID,
		LineID:       e.LineID,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
		DeletedAt:    deletedAt,
	}
}
