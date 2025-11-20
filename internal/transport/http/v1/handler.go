package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilindan-dev/CommentTree/internal/usecase"
	"github.com/wb-go/wbf/zlog"
)

// Handler represents the HTTP handler for comment-related endpoints.
type Handler struct {
	service usecase.CommentService
	logger  *zlog.Zerolog
}

// NewHandler creates a new instance of Handler.
func NewHandler(service usecase.CommentService, logger *zlog.Zerolog) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// commentQuery represents the query parameters for fetching comments.
type commentQuery struct {
	Page   int32  `form:"page,default=1"`
	Limit  int32  `form:"limit,default=10"`
	Sort   string `form:"sort,default=desc"`
	Parent string `form:"parent"`
	Query  string `form:"query"`
}

func (h *Handler) CreateComment(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn().Err(err).Msg("Binding failed, invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	comment, err := h.service.CreateComment(c.Request.Context(), req.ParentID, req.Text)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create comment")
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toResponseDTO(*comment))
}

func (h *Handler) GetComments(c *gin.Context) {
	var params commentQuery
	if err := c.ShouldBindQuery(&params); err != nil {
		h.logger.Warn().Err(err).Msg("Binding failed, invalid query parameters")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}
	offset := pageToOffset(params.Page, params.Limit)
	switch {
	case params.Parent != "":
		sortOrder, err := usecase.ParseSortOrder(params.Sort)
		if err != nil {
			h.logger.Warn().Err(err).Msg("Invalid sort order")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort order"})
			return
		}

		comments, err := h.service.GetCommentSubtree(c.Request.Context(), params.Parent, sortOrder)
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to get comment subtree")
			h.handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, toResponseDTOs(comments))
		return
	case params.Query != "":
		comments, err := h.service.SearchComments(c.Request.Context(), params.Query, params.Limit, offset)
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to search comments")
			h.handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, toResponseDTOs(comments))
		return
	default:
		sortOrder, err := usecase.ParseSortOrder(params.Sort)
		if err != nil {
			h.logger.Warn().Err(err).Msg("Invalid sort order")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort order"})
			return
		}
		comments, err := h.service.GetRootsComments(c.Request.Context(), params.Limit, offset, sortOrder)
		if err != nil {
			h.logger.Error().Err(err).Msg("Failed to get root comments")
			h.handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, toResponseDTOs(comments))
	}
}

func (h *Handler) SoftDeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	if commentID == "" {
		h.logger.Warn().Msg("Comment ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "comment ID is required"})
		return
	}
	comment, err := h.service.SoftDeleteComment(c.Request.Context(), commentID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to soft delete comment")
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, toResponseDTO(*comment))
}

func (h *Handler) handleError(c *gin.Context, err error) {
	if errors.Is(err, usecase.ErrTextEmpty) ||
		errors.Is(err, usecase.ErrInvalidUUID) ||
		errors.Is(err, usecase.ErrInvalidSortOrder) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

func pageToOffset(page int32, limit int32) int32 {
	return (page - 1) * limit
}
