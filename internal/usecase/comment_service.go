package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ilindan-dev/CommentTree/internal/domain/model"
	"github.com/ilindan-dev/CommentTree/internal/domain/repository"
	"github.com/wb-go/wbf/zlog"
)

// Ensure commentService implements CommentService.
var _ CommentService = (*commentService)(nil)

// commentService implements CommentService.
type commentService struct {
	commentRepository repository.CommentRepository
	logger            *zlog.Zerolog
}

// NewCommentService creates a new instance of CommentService.
func NewCommentService(commentRepository repository.CommentRepository, logger *zlog.Zerolog) CommentService {
	return &commentService{
		commentRepository: commentRepository,
		logger:            logger,
	}
}

// CreateComment creates a new comment with the given parent ID and text.
func (s *commentService) CreateComment(ctx context.Context, parentID string, text string) (*model.Comment, error) {
	var parentUUID *uuid.UUID
	if text == "" {
		s.logger.Error().Msg("comment text cannot be empty")
		return nil, ErrTextEmpty
	}
	if parentID == "" {
		parentUUID = nil
	} else {
		parseUUID, err := uuid.Parse(parentID)
		if err != nil {
			s.logger.Error().Err(err).Str("parentID", parentID).Msg("failed to parse parent id")
			return nil, ErrInvalidUUID
		}
		parentUUID = &parseUUID
	}
	comment, err := s.commentRepository.CreateComment(ctx, parentUUID, text)
	if err != nil {
		s.logger.Error().Err(err).Str("parentID", parentID).Msg("failed to create comment")
		return nil, err
	}
	return comment, nil
}

// GetRootsComments retrieves root comments with pagination and sorting.
func (s *commentService) GetRootsComments(ctx context.Context, limit int32, offset int32, sortOrder SortOrder) ([]model.Comment, error) {
	switch sortOrder {
	case SortOrderAsc:
		comments, err := s.commentRepository.GetRootCommentsAsc(ctx, limit, offset)
		if err != nil {
			s.logger.Error().Err(err).Stringer("sorterOrder", sortOrder).
				Int32("limit", limit).Int32("offset", offset).Msg("failed to get root comments")
			return nil, err
		}
		return comments, nil
	case SortOrderDesc:
		comments, err := s.commentRepository.GetRootCommentsDesc(ctx, limit, offset)
		if err != nil {
			s.logger.Error().Err(err).Stringer("sorterOrder", sortOrder).
				Int32("limit", limit).Int32("offset", offset).Msg("failed to get root comments")
			return nil, err
		}
		return comments, nil
	default:
		return nil, ErrInvalidSortOrder
	}
}

// SearchComments searches for comments containing the query string with pagination.
func (s *commentService) SearchComments(ctx context.Context, query string, limit int32, offset int32) ([]model.Comment, error) {
	comments, err := s.commentRepository.SearchComments(ctx, query, limit, offset)
	if err != nil {
		s.logger.Error().Err(err).Str("query", query).
			Int32("limit", limit).Int32("offset", offset).Msg("failed to search comments")
		return nil, err
	}
	return comments, nil
}

// SoftDeleteComment marks a comment as deleted without removing it from the database.
func (s *commentService) SoftDeleteComment(ctx context.Context, commentID string) (*model.Comment, error) {
	parsedUUID, err := uuid.Parse(commentID)
	if err != nil {
		s.logger.Error().Err(err).Str("commentID", commentID).Msg("failed to parse comment id")
		return nil, ErrInvalidUUID
	}
	comment, err := s.commentRepository.SoftDeleteComment(ctx, parsedUUID)
	if err != nil {
		s.logger.Error().Err(err).Str("commentID", commentID).Msg("failed to soft delete comment")
		return nil, err
	}
	return comment, nil
}

// GetCommentSubtree retrieves the subtree of comments under a root comment with sorting.
func (s *commentService) GetCommentSubtree(ctx context.Context, rootCommentID string, sortOrder SortOrder) ([]model.Comment, error) {
	commentUUID, err := uuid.Parse(rootCommentID)
	if err != nil {
		s.logger.Error().Err(err).Str("rootCommentID", rootCommentID).Msg("failed to parse root comment id")
		return nil, ErrInvalidUUID
	}
	switch sortOrder {
	case SortOrderAsc:
		comments, err := s.commentRepository.GetCommentSubtreeAsc(ctx, commentUUID)
		if err != nil {
			s.logger.Error().Err(err).Str("rootCommentID", rootCommentID).Msg("failed to get comment subtree")
			return nil, err
		}
		return comments, nil
	case SortOrderDesc:
		comments, err := s.commentRepository.GetCommentSubtreeDesc(ctx, commentUUID)
		if err != nil {
			s.logger.Error().Err(err).Str("rootCommentID", rootCommentID).Msg("failed to get comment subtree")
			return nil, err
		}
		return comments, nil
	default:
		return nil, ErrInvalidSortOrder
	}
}
