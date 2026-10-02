package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raifuki/task-management/internal/services"
	"github.com/raifuki/task-management/pkg/response"
	"github.com/raifuki/task-management/pkg/validator"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}

// Register godoc
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body services.RegisterInput true "Register body"
// @Success 201 {object} response.APIResponse
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var in services.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}
	if errs := validator.Validate(in); errs != nil {
		response.BadRequest(c, "Validation failed", errs)
		return
	}
	user, err := h.svc.Register(in)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}
	response.Success(c, http.StatusCreated, "Register success", user)
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body services.LoginInput true "Login body"
// @Success 200 {object} response.APIResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var in services.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}
	if errs := validator.Validate(in); errs != nil {
		response.BadRequest(c, "Validation failed", errs)
		return
	}
	token, user, err := h.svc.Login(in)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Login success", gin.H{
		"token": token,
		"user":  user,
	})
}
