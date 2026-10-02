package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/raifuki/task-management/internal/models"
)

type CommentRepository struct{ db *gorm.DB }

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(c *models.Comment) error { return r.db.Create(c).Error }

func (r *CommentRepository) FindByTask(taskID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.db.Preload("User").Where("task_id = ?", taskID).Order("id ASC").Find(&comments).Error
	return comments, err
}

func (r *CommentRepository) FindByID(id uint) (*models.Comment, error) {
	var c models.Comment
	err := r.db.Preload("User").First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &c, err
}

func (r *CommentRepository) Update(c *models.Comment) error {
	return r.db.Save(c).Error
}

func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Comment{}, id).Error
}
