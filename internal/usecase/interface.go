package usecase

import (
	"context"
	"errors"

	"github.com/ilindan-dev/CommentTree/internal/domain/model"
)

// SortOrder represents the order in which comments can be sorted.
type SortOrder string

const (
	// SortOrderAsc indicates ascending order.
	SortOrderAsc SortOrder = "ASC"
	// SortOrderDesc indicates descending order.
	SortOrderDesc SortOrder = "DESC"
)

func (s SortOrder) Asc() SortOrder {
	return SortOrderAsc
}

func (s SortOrder) Desc() SortOrder {
	return SortOrderDesc
}

func (s SortOrder) String() string {
	return string(s)
}

// ErrInvalidSortOrder is returned when an invalid sort order is provided.
var ErrInvalidSortOrder = errors.New("invalid sort order")

// ErrInvalidUUID is returned when an invalid UUID is provided.
var ErrInvalidUUID = errors.New("invalid UUID")

// ErrTextEmpty is returned when the comment text is empty.
var ErrTextEmpty = errors.New("comment text cannot be empty")

// ParseSortOrder parses a string into a SortOrder type.
func ParseSortOrder(order string) (SortOrder, error) {
	switch {
	case "ASC" == order || "asc" == order:
		return SortOrderAsc, nil
	case "DESC" == order || "desc" == order:
		return SortOrderDesc, nil
	default:
		return "", ErrInvalidSortOrder
	}
}

// CommentService defines the interface for comment-related use cases.
type CommentService interface {
	// CreateComment creates a new comment with the given parent ID and text.
	CreateComment(ctx context.Context, parentID string, text string) (*model.Comment, error)

	// GetRootsComments retrieves root comments with pagination and sorting.
	GetRootsComments(ctx context.Context, limit int32, offset int32, sortOrder SortOrder) ([]model.Comment, error)

	// SearchComments searches for comments containing the query string with pagination.
	SearchComments(ctx context.Context, query string, limit int32, offset int32) ([]model.Comment, error)

	// SoftDeleteComment marks a comment as deleted without removing it from the database.
	SoftDeleteComment(ctx context.Context, commentID string) (*model.Comment, error)

	// GetCommentSubtree retrieves the subtree of comments under a root comment with sorting.
	GetCommentSubtree(ctx context.Context, rootCommentID string, sortOrder SortOrder) ([]model.Comment, error)
}
