package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/raifuki/task-management/internal/middlewares"
	"github.com/raifuki/task-management/internal/services"
	"github.com/raifuki/task-management/pkg/response"
	"github.com/raifuki/task-management/pkg/validator"
)

type CommentHandler struct{ svc *services.CommentService }

func NewCommentHandler(s *services.CommentService) *CommentHandler {
	return &CommentHandler{svc: s}
}

func (h *CommentHandler) Create(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.Param("taskId"))
	var in services.CommentInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "Invalid body", err.Error())
		return
	}
	if errs := validator.Validate(in); errs != nil {
		response.BadRequest(c, "Validation failed", errs)
		return
	}
	uid := c.GetUint(middlewares.CtxUserID)
	cm, err := h.svc.Create(uint(taskID), uid, in)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusCreated, "Comment created", cm)
}

func (h *CommentHandler) List(c *gin.Context) {
	taskID, _ := strconv.Atoi(c.Param("taskId"))
	cm, err := h.svc.ListByTask(uint(taskID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "OK", cm)
}

func (h *CommentHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.Delete(uint(id)); err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Deleted", nil)
}
