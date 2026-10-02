package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:150;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	OwnerID     uint           `gorm:"not null;index" json:"owner_id"`
	Owner       *User          `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // soft delete

	Tasks []Task `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
}
