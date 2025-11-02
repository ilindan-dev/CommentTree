package model

import (
	"time"

	"github.com/google/uuid"
)

// Comment represents a comment entity in the system.
type Comment struct {
	ID        uuid.UUID
	ParentID  *uuid.UUID
	Text      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsRoot checks if the comment is a root comment (i.e., it has no parent).
func (c *Comment) IsRoot() bool {
	return c.ParentID == nil
}

// IsDeleted checks if the comment has been deleted.
func (c *Comment) IsDeleted() bool {
	return c.Text == "[deleted]"
}
