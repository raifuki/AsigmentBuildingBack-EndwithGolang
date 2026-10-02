package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/raifuki/task-management/internal/models"
)

type ProjectRepository struct{ db *gorm.DB }

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(p *models.Project) error {
	return r.db.Create(p).Error
}

func (r *ProjectRepository) FindAll(ownerID uint, page, limit int) ([]models.Project, int64, error) {
	var projects []models.Project
	var total int64

	q := r.db.Model(&models.Project{})
	if ownerID > 0 {
		q = q.Where("owner_id = ?", ownerID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := q.Preload("Owner").Order("id DESC").Offset(offset).Limit(limit).Find(&projects).Error
	return projects, total, err
}

func (r *ProjectRepository) FindByID(id uint) (*models.Project, error) {
	var p models.Project
	err := r.db.Preload("Owner").First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}

func (r *ProjectRepository) Update(p *models.Project) error {
	return r.db.Save(p).Error
}

func (r *ProjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Project{}, id).Error
}
