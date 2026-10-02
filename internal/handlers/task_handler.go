package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/raifuki/task-management/internal/services"
	"github.com/raifuki/task-management/pkg/response"
	"github.com/raifuki/task-management/pkg/validator"
)

type TaskHandler struct{ svc *services.TaskService }

func NewTaskHandler(s *services.TaskService) *TaskHandler { return &TaskHandler{svc: s} }

func (h *TaskHandler) Create(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("projectId"))
	var in services.TaskInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "Invalid body", err.Error())
		return
	}
	if errs := validator.Validate(in); errs != nil {
		response.BadRequest(c, "Validation failed", errs)
		return
	}
	t, err := h.svc.Create(c.Request.Context(), uint(projectID), in)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusCreated, "Task created", t)
}

func (h *TaskHandler) List(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Query("project_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	tasks, total, err := h.svc.List(c.Request.Context(), uint(projectID), status, page, limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "OK", gin.H{
		"items": tasks,
		"meta":  gin.H{"page": page, "limit": limit, "total": total},
	})
}

func (h *TaskHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	t, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "OK", t)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var in services.TaskInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "Invalid body", err.Error())
		return
	}
	if errs := validator.Validate(in); errs != nil {
		response.BadRequest(c, "Validation failed", errs)
		return
	}
	t, err := h.svc.Update(c.Request.Context(), uint(id), in)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Updated", t)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Deleted", nil)
}
