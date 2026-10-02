package services

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/raifuki/task-management/internal/cache"
	"github.com/raifuki/task-management/internal/models"
	"github.com/raifuki/task-management/internal/repositories"
)

const (
	projectListTTL   = 5 * time.Minute
	projectDetailTTL = 10 * time.Minute
)

type ProjectService struct {
	repo  *repositories.ProjectRepository
	cache cache.Cache
}

func NewProjectService(r *repositories.ProjectRepository, c cache.Cache) *ProjectService {
	return &ProjectService{repo: r, cache: c}
}

type ProjectInput struct {
	Name        string `json:"name" validate:"required,min=3,max=150"`
	Description string `json:"description" validate:"max=1000"`
}

func (s *ProjectService) Create(ctx context.Context, ownerID uint, in ProjectInput) (*models.Project, error) {
	p := &models.Project{
		Name:        in.Name,
		Description: in.Description,
		OwnerID:     ownerID,
	}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	// Invalidate list cache
	s.invalidateList(ctx)
	return p, nil
}

func (s *ProjectService) List(ctx context.Context, ownerID uint, page, limit int) ([]models.Project, int64, error) {
	key := cache.ListKey("projects", page, limit)
	if ownerID > 0 {
		key += ":owner=" + uintToStr(ownerID)
	}

	// Try cache
	type cached struct {
		Items []models.Project `json:"items"`
		Total int64            `json:"total"`
	}
	var c cached
	if err := s.cache.Get(ctx, key, &c); err == nil {
		log.Printf("🟢 Cache HIT: %s", key)
		return c.Items, c.Total, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		log.Printf("⚠️ Cache error: %v", err)
	} else {
		log.Printf("🔴 Cache MISS: %s", key)
	}

	// Fallback DB
	items, total, err := s.repo.FindAll(ownerID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// Save cache
	_ = s.cache.Set(ctx, key, cached{Items: items, Total: total}, projectListTTL)
	return items, total, nil
}

func (s *ProjectService) Get(ctx context.Context, id uint) (*models.Project, error) {
	key := cache.Key("projects", "detail", uintToStr(id))

	var p models.Project
	if err := s.cache.Get(ctx, key, &p); err == nil {
		log.Printf("🟢 Cache HIT: %s", key)
		return &p, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		log.Printf("⚠️ Cache error: %v", err)
	} else {
		log.Printf("🔴 Cache MISS: %s", key)
	}

	// Fallback DB
	proj, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if proj == nil {
		return nil, errors.New("project not found")
	}

	_ = s.cache.Set(ctx, key, proj, projectDetailTTL)
	return proj, nil
}

func (s *ProjectService) Update(ctx context.Context, id, userID uint, in ProjectInput, isAdmin bool) (*models.Project, error) {
	p, err := s.repo.FindByID(id)
	if err != nil || p == nil {
		return nil, errors.New("project not found")
	}
	if p.OwnerID != userID && !isAdmin {
		return nil, errors.New("forbidden")
	}
	p.Name = in.Name
	p.Description = in.Description
	if err := s.repo.Update(p); err != nil {
		return nil, err
	}

	// Invalidate cache
	s.invalidateDetail(ctx, id)
	s.invalidateList(ctx)
	return p, nil
}

func (s *ProjectService) Delete(ctx context.Context, id, userID uint, isAdmin bool) error {
	p, err := s.repo.FindByID(id)
	if err != nil || p == nil {
		return errors.New("project not found")
	}
	if p.OwnerID != userID && !isAdmin {
		return errors.New("forbidden")
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.invalidateDetail(ctx, id)
	s.invalidateList(ctx)
	return nil
}

// ===== Cache invalidation helpers =====

func (s *ProjectService) invalidateDetail(ctx context.Context, id uint) {
	key := cache.Key("projects", "detail", uintToStr(id))
	_ = s.cache.Delete(ctx, key)
}

func (s *ProjectService) invalidateList(ctx context.Context) {
	_ = s.cache.DeleteByPattern(ctx, "projects:list:*")
}

// Helper nhỏ để tránh import strconv ở nhiều chỗ
func uintToStr(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}
