package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:200;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      TaskStatus     `gorm:"type:varchar(20);default:'todo'" json:"status"`
	Priority    int            `gorm:"default:1" json:"priority"` // 1=low, 2=med, 3=high
	DueDate     *time.Time     `json:"due_date,omitempty"`
	ProjectID   uint           `gorm:"not null;index" json:"project_id"`
	AssigneeID  *uint          `gorm:"index" json:"assignee_id,omitempty"`
	Assignee    *User          `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Comments []Comment `gorm:"foreignKey:TaskID" json:"comments,omitempty"`
}