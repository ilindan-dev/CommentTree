package v1

import (
	"time"

	"github.com/ilindan-dev/CommentTree/internal/domain/model"
)

// CreateCommentRequest represents the request payload for creating a comment.
type CreateCommentRequest struct {
	ParentID string `json:"parent_id"`
	Text     string `json:"text" binding:"required"`
}

// CommentResponse represents the response payload for a comment.
type CommentResponse struct {
	ID        string    `json:"id"`
	ParentID  *string   `json:"parent_id,omitempty"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDeleted bool      `json:"is_deleted"`
}

// toResponseDTO converts a model.Comment to a CommentResponse DTO.
func toResponseDTO(comment model.Comment) CommentResponse {
	var parentID *string
	if comment.ParentID != nil {
		p := comment.ParentID.String()
		parentID = &p
	}
	return CommentResponse{
		ID:        comment.ID.String(),
		ParentID:  parentID,
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		IsDeleted: comment.IsDeleted(),
	}
}

// toResponseDTOs converts a slice of model.Comment to a slice of CommentResponse DTOs.
func toResponseDTOs(comments []model.Comment) []CommentResponse {
	if len(comments) == 0 {
		return []CommentResponse{}
	}

	responses := make([]CommentResponse, len(comments))
	for i, comment := range comments {
		responses[i] = toResponseDTO(comment)
	}
	return responses
}
