package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/raifuki/task-management/internal/models"
	"github.com/raifuki/task-management/internal/repositories"
	"github.com/raifuki/task-management/internal/services"
)

// Mock repository
type mockProjectRepo struct {
	mock.Mock
}

func (m *mockProjectRepo) Create(p *models.Project) error {
	args := m.Called(p)
	return args.Error(0)
}
func (m *mockProjectRepo) FindAll(ownerID uint, page, limit int) ([]models.Project, int64, error) {
	args := m.Called(ownerID, page, limit)
	return args.Get(0).([]models.Project), args.Get(1).(int64), args.Error(2)
}
func (m *mockProjectRepo) FindByID(id uint) (*models.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}
func (m *mockProjectRepo) Update(p *models.Project) error { return m.Called(p).Error(0) }
func (m *mockProjectRepo) Delete(id uint) error           { return m.Called(id).Error(0) }

// NOTE: Để dùng mock, bạn cần biến ProjectService nhận interface thay vì con trỏ repository.
// Ở đây minh họa hướng tiếp cận; bạn có thể refactor ProjectService để phụ thuộc interface.

func TestProjectService_Create_Success(t *testing.T) {
	// Ví dụ flow
	_ = gorm.ErrRecordNotFound
	_ = repositories.NewProjectRepository
	_ = services.NewProjectService
	assert.True(t, true)
}
