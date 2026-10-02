package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raifuki/task-management/pkg/response"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("🔥 Panic: %v", r)
				response.Error(c, http.StatusInternalServerError, "Internal server error", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}