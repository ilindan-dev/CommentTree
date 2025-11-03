package postgres

import (
	"github.com/google/uuid"
	"github.com/ilindan-dev/CommentTree/internal/domain/model"
	"github.com/ilindan-dev/CommentTree/internal/storage/postgres/db"
	"github.com/jackc/pgx/v5/pgtype"
)

// dbCommentToDomainComment maps a db.Comment to a model.Comment.
func dbCommentToDomainComment(dbComment *db.Comment) *model.Comment {
	var parentID *uuid.UUID
	if dbComment.ParentID.Valid {
		uid := uuid.UUID(dbComment.ParentID.Bytes)
		parentID = &uid
	}
	return &model.Comment{
		ID:        dbComment.ID.Bytes,
		ParentID:  parentID,
		Text:      dbComment.Text,
		CreatedAt: dbComment.CreatedAt.Time,
		UpdatedAt: dbComment.UpdatedAt.Time,
	}
}

// dbCommentsToDomainComments maps a slice of db.Comment to a slice of model.Comment.
func dbCommentsToDomainComments(dbComments []db.Comment) []model.Comment {
	if dbComments == nil {
		return nil
	}
	domainComments := make([]model.Comment, len(dbComments))
	for i, dbComment := range dbComments {
		domainComments[i] = *dbCommentToDomainComment(&dbComment)
	}
	return domainComments
}

// uuidPtrToPgtypeUUID converts a *uuid.UUID to pgtype.UUID.
func uuidPtrToPgtypeUUID(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{
		Bytes: *id, Valid: true,
	}
}

func pgtypeUUIDToUUIDPtr(pgID pgtype.UUID) *uuid.UUID {
	if !pgID.Valid {
		return nil
	}
	uid := uuid.UUID(pgID.Bytes)
	return &uid
}
