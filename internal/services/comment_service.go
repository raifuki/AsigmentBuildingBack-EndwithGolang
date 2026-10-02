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

func (s *CommentService) Delete(id uint) error {
	return s.repo.Delete(id)
}
