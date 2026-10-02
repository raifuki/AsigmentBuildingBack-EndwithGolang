package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{Success: true, Message: message, Data: data})
}

func Error(c *gin.Context, status int, message string, errs interface{}) {
	c.JSON(status, APIResponse{Success: false, Message: message, Errors: errs})
}

func BadRequest(c *gin.Context, msg string, errs interface{}) {
	Error(c, http.StatusBadRequest, msg, errs)
}

func Unauthorized(c *gin.Context, msg string) {
	Error(c, http.StatusUnauthorized, msg, nil)
}

func Forbidden(c *gin.Context, msg string) {
	Error(c, http.StatusForbidden, msg, nil)
}

func NotFound(c *gin.Context, msg string) {
	Error(c, http.StatusNotFound, msg, nil)
}

func InternalError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, msg, nil)
}