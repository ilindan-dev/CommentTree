package v1

import "github.com/gin-gonic/gin"

// RegisterRoutes registers the comment-related routes to the given Gin router.
func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/v1")
	{
		api.POST("/comments", h.CreateComment)
		api.GET("/comments", h.GetComments)
		api.DELETE("/comments/:id", h.SoftDeleteComment)
	}
}
