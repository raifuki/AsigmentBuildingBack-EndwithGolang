package services

import (
	"errors"

	"github.com/raifuki/task-management/internal/models"
	"github.com/raifuki/task-management/internal/repositories"
)

type CommentService struct {
	repo     *repositories.CommentRepository
	taskRepo *repositories.TaskRepository
}

func NewCommentService(r *repositories.CommentRepository, tr *repositories.TaskRepository) *CommentService {
	return &CommentService{repo: r, taskRepo: tr}
}

type CommentInput struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

func (s *CommentService) Create(taskID, userID uint, in CommentInput) (*models.Comment, error) {
	t, err := s.taskRepo.FindByID(taskID)
	if err != nil || t == nil {
		return nil, errors.New("task not found")
	}
	c := &models.Comment{Content: in.Content, TaskID: taskID, UserID: userID}
	return c, s.repo.Create(c)
}

func (s *CommentService) ListByTask(taskID uint) ([]models.Comment, error) {
	return s.repo.FindByTask(taskID)
}

func (s *CommentService) Get(id uint) (*models.Comment, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, errors.New("comment not found")
	}
	return c, nil
}

func (s *CommentService) Update(id, userID uint, in CommentInput, isAdmin bool) (*models.Comment, error) {
	c, err := s.repo.FindByID(id)
	if err != nil || c == nil {
		return nil, errors.New("comment not found")
	}
	if c.UserID != userID && !isAdmin {
		return nil, errors.New("forbidden: you can only edit your own comments")
	}
	c.Content = in.Content
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CommentService) Delete(id uint) error {
	return s.repo.Delete(id)
}
