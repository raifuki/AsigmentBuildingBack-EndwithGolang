package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/raifuki/task-management/internal/models"
)

type TaskRepository struct{ db *gorm.DB }

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(t *models.Task) error {
	return r.db.Create(t).Error
}

func (r *TaskRepository) FindAll(projectID uint, status string, page, limit int) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	q := r.db.Model(&models.Task{})
	if projectID > 0 {
		q = q.Where("project_id = ?", projectID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := q.Preload("Assignee").Order("priority DESC, id DESC").
		Offset(offset).Limit(limit).Find(&tasks).Error
	return tasks, total, err
}

func (r *TaskRepository) FindByID(id uint) (*models.Task, error) {
	var t models.Task
	err := r.db.Preload("Assignee").Preload("Comments.User").First(&t, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, err
}

func (r *TaskRepository) Update(t *models.Task) error { return r.db.Save(t).Error }
func (r *TaskRepository) Delete(id uint) error        { return r.db.Delete(&models.Task{}, id).Error }
