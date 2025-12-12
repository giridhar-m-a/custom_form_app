package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// UnknownFieldsMiddleware sets up strict JSON binding to reject unknown fields
func UnknownFieldsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Enable DisallowUnknownFields for JSON binding
		// This will cause ShouldBindJSON to return an error if the request contains unknown fields
		binding.EnableDecoderDisallowUnknownFields = true

		c.Next()
	}
}
