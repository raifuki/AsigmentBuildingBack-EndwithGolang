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

type ProjectHandler struct{ svc *services.ProjectService }

func NewProjectHandler(s *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{svc: s}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var in services.ProjectInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "Invalid body", err.Error())
		return
	}
	if errs := validator.Validate(in); errs != nil {
		response.BadRequest(c, "Validation failed", errs)
		return
	}
	uid := c.GetUint(middlewares.CtxUserID)
	p, err := h.svc.Create(c.Request.Context(), uid, in)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "Project created", p)
}

func (h *ProjectHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	projects, total, err := h.svc.List(c.Request.Context(), 0, page, limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "OK", gin.H{
		"items": projects,
		"meta":  gin.H{"page": page, "limit": limit, "total": total},
	})
}

func (h *ProjectHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	p, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "OK", p)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var in services.ProjectInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "Invalid body", err.Error())
		return
	}
	if errs := validator.Validate(in); errs != nil {
		response.BadRequest(c, "Validation failed", errs)
		return
	}
	uid := c.GetUint(middlewares.CtxUserID)
	role := c.GetString(middlewares.CtxRole)
	p, err := h.svc.Update(c.Request.Context(), uint(id), uid, in, role == "admin")
	if err != nil {
		response.Error(c, http.StatusForbidden, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Updated", p)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uid := c.GetUint(middlewares.CtxUserID)
	role := c.GetString(middlewares.CtxRole)
	if err := h.svc.Delete(c.Request.Context(), uint(id), uid, role == "admin"); err != nil {
		response.Error(c, http.StatusForbidden, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusOK, "Deleted", nil)
}
