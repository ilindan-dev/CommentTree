package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ilindan-dev/CommentTree/internal/domain/model"
)

// CommentRepository defines the interface for comment data operations.
type CommentRepository interface {
	// CreateComment creates a new comment with an optional parent ID.
	CreateComment(ctx context.Context, parentID *uuid.UUID, text string) (*model.Comment, error)

	// GetRootCommentsDesc retrieves root comments in descending order with pagination.
	GetRootCommentsDesc(ctx context.Context, limit int32, offset int32) ([]model.Comment, error)

	// GetRootCommentsAsc retrieves root comments in ascending order with pagination.
	GetRootCommentsAsc(ctx context.Context, limit int32, offset int32) ([]model.Comment, error)

	// SearchComments searches for comments containing the query string with pagination.
	SearchComments(ctx context.Context, query string, limit int32, offset int32) ([]model.Comment, error)

	// SoftDeleteComment marks a comment as deleted without removing it from the database.
	SoftDeleteComment(ctx context.Context, commentID uuid.UUID) (*model.Comment, error)

	// GetCommentSubtreeAsc retrieves the subtree of comments under a root comment in ascending order.
	GetCommentSubtreeAsc(ctx context.Context, rootCommentID uuid.UUID) ([]model.Comment, error)

	// GetCommentSubtreeDesc retrieves the subtree of comments under a root comment in descending order.
	GetCommentSubtreeDesc(ctx context.Context, rootCommentID uuid.UUID) ([]model.Comment, error)
}
