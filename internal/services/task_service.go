package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/raifuki/task-management/internal/cache"
	"github.com/raifuki/task-management/internal/models"
	"github.com/raifuki/task-management/internal/repositories"
)

const (
	taskListTTL   = 2 * time.Minute
	taskDetailTTL = 5 * time.Minute
)

type TaskService struct {
	repo        *repositories.TaskRepository
	projectRepo *repositories.ProjectRepository
	cache       cache.Cache
}

func NewTaskService(
	r *repositories.TaskRepository,
	pr *repositories.ProjectRepository,
	c cache.Cache,
) *TaskService {
	return &TaskService{repo: r, projectRepo: pr, cache: c}
}

type TaskInput struct {
	Title       string     `json:"title" validate:"required,min=3,max=200"`
	Description string     `json:"description" validate:"max=2000"`
	Status      string     `json:"status" validate:"omitempty,oneof=todo in_progress done"`
	Priority    int        `json:"priority" validate:"omitempty,min=1,max=3"`
	DueDate     *time.Time `json:"due_date"`
	AssigneeID  *uint      `json:"assignee_id"`
}

func (s *TaskService) Create(ctx context.Context, projectID uint, in TaskInput) (*models.Task, error) {
	p, err := s.projectRepo.FindByID(projectID)
	if err != nil || p == nil {
		return nil, errors.New("project not found")
	}

	status := models.TaskStatus(in.Status)
	if status == "" {
		status = models.StatusTodo
	}
	priority := in.Priority
	if priority == 0 {
		priority = 1
	}

	t := &models.Task{
		Title:       in.Title,
		Description: in.Description,
		Status:      status,
		Priority:    priority,
		DueDate:     in.DueDate,
		ProjectID:   projectID,
		AssigneeID:  in.AssigneeID,
	}
	if err := s.repo.Create(t); err != nil {
		return nil, err
	}

	s.invalidateTaskList(ctx, projectID)
	return t, nil
}

func (s *TaskService) List(ctx context.Context, projectID uint, status string, page, limit int) ([]models.Task, int64, error) {
	key := cache.ListKey("tasks", page, limit)
	if projectID > 0 {
		key += ":project=" + uintToStr(projectID)
	}
	if status != "" {
		key += ":status=" + status
	}

	type cached struct {
		Items []models.Task `json:"items"`
		Total int64         `json:"total"`
	}
	var c cached
	if err := s.cache.Get(ctx, key, &c); err == nil {
		log.Printf("🟢 Cache HIT: %s", key)
		return c.Items, c.Total, nil
	} else if errors.Is(err, cache.ErrCacheMiss) {
		log.Printf("🔴 Cache MISS: %s", key)
	} else {
		log.Printf("⚠️ Cache error: %v", err)
	}

	items, total, err := s.repo.FindAll(projectID, status, page, limit)
	if err != nil {
		return nil, 0, err
	}

	_ = s.cache.Set(ctx, key, cached{Items: items, Total: total}, taskListTTL)
	return items, total, nil
}

func (s *TaskService) Get(ctx context.Context, id uint) (*models.Task, error) {
	key := cache.Key("tasks", "detail", uintToStr(id))

	var t models.Task
	if err := s.cache.Get(ctx, key, &t); err == nil {
		log.Printf("🟢 Cache HIT: %s", key)
		return &t, nil
	} else if errors.Is(err, cache.ErrCacheMiss) {
		log.Printf("🔴 Cache MISS: %s", key)
	}

	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	_ = s.cache.Set(ctx, key, task, taskDetailTTL)
	return task, nil
}

func (s *TaskService) Update(ctx context.Context, id uint, in TaskInput) (*models.Task, error) {
	t, err := s.repo.FindByID(id)
	if err != nil || t == nil {
		return nil, errors.New("task not found")
	}
	t.Title = in.Title
	t.Description = in.Description
	if in.Status != "" {
		t.Status = models.TaskStatus(in.Status)
	}
	if in.Priority > 0 {
		t.Priority = in.Priority
	}
	t.DueDate = in.DueDate
	t.AssigneeID = in.AssigneeID

	if err := s.repo.Update(t); err != nil {
		return nil, err
	}

	s.invalidateTaskDetail(ctx, id)
	s.invalidateTaskList(ctx, t.ProjectID)
	return t, nil
}

func (s *TaskService) Delete(ctx context.Context, id uint) error {
	t, err := s.repo.FindByID(id)
	if err != nil || t == nil {
		return errors.New("task not found")
	}
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	s.invalidateTaskDetail(ctx, id)
	s.invalidateTaskList(ctx, t.ProjectID)
	return nil
}

func (s *TaskService) invalidateTaskDetail(ctx context.Context, id uint) {
	_ = s.cache.Delete(ctx, cache.Key("tasks", "detail", uintToStr(id)))
}

func (s *TaskService) invalidateTaskList(ctx context.Context, projectID uint) {
	_ = s.cache.DeleteByPattern(ctx, "tasks:list:*")
}
