package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ilindan-dev/CommentTree/internal/domain/model"
	"github.com/ilindan-dev/CommentTree/internal/domain/repository"
	"github.com/ilindan-dev/CommentTree/internal/storage/postgres/db"
	"github.com/wb-go/wbf/zlog"
)

// Ensure commentRepository implements repository.CommentRepository.
var _ repository.CommentRepository = (*commentRepository)(nil)

// commentRepository implements repository.CommentRepository.
type commentRepository struct {
	q      *db.Queries
	logger *zlog.Zerolog
}

// NewCommentRepository creates a new instance of CommentRepository.
func NewCommentRepository(pool db.DBTX, logger *zlog.Zerolog) repository.CommentRepository {
	return &commentRepository{
		q:      db.New(pool),
		logger: logger,
	}
}

// CreateComment creates a new comment in the database.
func (r *commentRepository) CreateComment(ctx context.Context, parentID *uuid.UUID, text string) (*model.Comment, error) {
	params := db.CreateCommentParams{
		ParentID: uuidPtrToPgtypeUUID(parentID),
		Text:     text,
	}

	dbComment, err := r.q.CreateComment(ctx, params)
	if err != nil {
		var parentStr string
		if parentID != nil {
			parentStr = parentID.String()
		} else {
			parentStr = "NULL"
		}
		r.logger.Error().Err(err).Str("parent", parentStr).Msg("failed to create comment")
		return nil, err
	}
	return &model.Comment{
		ID:        dbComment.ID.Bytes,
		ParentID:  pgtypeUUIDToUUIDPtr(dbComment.ParentID),
		Text:      dbComment.Text,
		CreatedAt: dbComment.CreatedAt.Time,
		UpdatedAt: dbComment.UpdatedAt.Time,
	}, nil
}

// GetRootCommentsDesc retrieves root comments in descending order with pagination.
func (r *commentRepository) GetRootCommentsDesc(ctx context.Context, limit int32, offset int32) ([]model.Comment, error) {
	params := db.GetRootCommentsDescParams{
		Limit:  limit,
		Offset: offset,
	}
	dbComments, err := r.q.GetRootCommentsDesc(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int32("limit", limit).Int32("offset", offset).Msg("failed to get root comments desc")
		return nil, err
	}
	domainComments := make([]model.Comment, len(dbComments))
	for i, row := range dbComments {
		domainComments[i] = model.Comment{
			ID:        row.ID.Bytes,
			ParentID:  pgtypeUUIDToUUIDPtr(row.ParentID),
			Text:      row.Text,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		}
	}
	return domainComments, nil
}

// GetRootCommentsAsc retrieves root comments in ascending order with pagination.
func (r *commentRepository) GetRootCommentsAsc(ctx context.Context, limit int32, offset int32) ([]model.Comment, error) {
	params := db.GetRootCommentsAscParams{
		Limit:  limit,
		Offset: offset,
	}
	dbComments, err := r.q.GetRootCommentsAsc(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Int32("limit", limit).Int32("offset", offset).Msg("failed to get root comments asc")
		return nil, err
	}
	domainComments := make([]model.Comment, len(dbComments))
	for i, row := range dbComments {
		domainComments[i] = model.Comment{
			ID:        row.ID.Bytes,
			ParentID:  pgtypeUUIDToUUIDPtr(row.ParentID),
			Text:      row.Text,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		}
	}
	return domainComments, nil
}

// SearchComments searches comments by text with pagination.
func (r *commentRepository) SearchComments(ctx context.Context, query string, limit int32, offset int32) ([]model.Comment, error) {
	params := db.SearchCommentsParams{
		Query:  query,
		Limit:  limit,
		Offset: offset,
	}
	dbComments, err := r.q.SearchComments(ctx, params)
	if err != nil {
		r.logger.Error().Err(err).Str("query", query).Int32("limit", limit).Int32("offset", offset).Msg("failed to search comments")
		return nil, err
	}

	domainComments := make([]model.Comment, len(dbComments))
	for i, row := range dbComments {
		domainComments[i] = model.Comment{
			ID:        row.ID.Bytes,
			ParentID:  pgtypeUUIDToUUIDPtr(row.ParentID),
			Text:      row.Text,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		}
	}
	return domainComments, nil
}

// SoftDeleteComment performs a soft delete on a comment by its ID.
func (r *commentRepository) SoftDeleteComment(ctx context.Context, commentID uuid.UUID) (*model.Comment, error) {
	dbComment, err := r.q.SoftDeleteComment(ctx, uuidPtrToPgtypeUUID(&commentID))
	if err != nil {
		r.logger.Error().Err(err).Str("comment_id", commentID.String()).Msg("failed to soft delete comment")
		return nil, err
	}
	return &model.Comment{
		ID:        dbComment.ID.Bytes,
		ParentID:  pgtypeUUIDToUUIDPtr(dbComment.ParentID),
		Text:      dbComment.Text,
		CreatedAt: dbComment.CreatedAt.Time,
		UpdatedAt: dbComment.UpdatedAt.Time,
	}, nil
}

// GetCommentSubtreeAsc retrieves the subtree of comments starting from the root comment ID in ascending order.
func (r *commentRepository) GetCommentSubtreeAsc(ctx context.Context, rootCommentID uuid.UUID) ([]model.Comment, error) {
	dbComments, err := r.q.GetCommentSubtreeAsc(ctx, uuidPtrToPgtypeUUID(&rootCommentID))
	if err != nil {
		r.logger.Error().Err(err).Str("root_comment_id", rootCommentID.String()).Msg("failed to get comment subtree asc")
		return nil, err
	}
	domainComments := make([]model.Comment, len(dbComments))
	for i := range dbComments {
		domainComments[i] = model.Comment{
			ID:        dbComments[i].ID.Bytes,
			ParentID:  pgtypeUUIDToUUIDPtr(dbComments[i].ParentID),
			Text:      dbComments[i].Text,
			CreatedAt: dbComments[i].CreatedAt.Time,
			UpdatedAt: dbComments[i].UpdatedAt.Time,
		}
	}
	return domainComments, nil
}

// GetCommentSubtreeDesc retrieves the subtree of comments starting from the root comment ID in descending order.
func (r *commentRepository) GetCommentSubtreeDesc(ctx context.Context, rootCommentID uuid.UUID) ([]model.Comment, error) {
	dbComments, err := r.q.GetCommentSubtreeDesc(ctx, uuidPtrToPgtypeUUID(&rootCommentID))
	if err != nil {
		r.logger.Error().Err(err).Str("root_comment_id", rootCommentID.String()).Msg("failed to get comment subtree desc")
		return nil, err
	}
	domainComments := make([]model.Comment, len(dbComments))
	for i := range dbComments {
		domainComments[i] = model.Comment{
			ID:        dbComments[i].ID.Bytes,
			ParentID:  pgtypeUUIDToUUIDPtr(dbComments[i].ParentID),
			Text:      dbComments[i].Text,
			CreatedAt: dbComments[i].CreatedAt.Time,
			UpdatedAt: dbComments[i].UpdatedAt.Time,
		}
	}
	return domainComments, nil
}
