-- name: CreateComment :one
INSERT INTO comments (parent_id, text)
VALUES (sqlc.arg(parent_id), sqlc.arg(text))
RETURNING *;

-- name: GetRootCommentsDesc :many
SELECT * FROM comments
WHERE parent_id is NULL
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_) OFFSET sqlc.arg(offset_);

-- name: GetRootCommentsAsc :many
SELECT * FROM comments
WHERE parent_id is NULL
ORDER BY created_at ASC
LIMIT sqlc.arg(limit_) OFFSET sqlc.arg(offset_);

-- name: SearchComments :many
SELECT * FROM comments
WHERE tsv @@ to_tsquery('russian', sqlc.arg(query))
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_) OFFSET sqlc.arg(offset_);

-- name: SoftDeleteComment :one
UPDATE comments
SET text = '[deleted]', updated_at = NOW(), tsv = NULL
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: GetCommentSubtreeDesc :many
WITH RECURSIVE comment_tree AS (
    SELECT
        c.*,
        ARRAY[EXTRACT(EPOCH FROM c.created_at) * -1] AS sort_path
    FROM comments c
    WHERE c.id = sqlc.arg(root_id)
    UNION ALL
    SELECT
        c.*,
        ct.sort_path || (EXTRACT(EPOCH FROM c.created_at) * -1)
    FROM comments c
             JOIN comment_tree ct ON c.parent_id = ct.id
)
SELECT * FROM comment_tree
ORDER BY sort_path;

-- name: GetCommentSubtreeAsc :many
WITH RECURSIVE comment_tree AS (
    SELECT
        c.*,
        ARRAY[c.created_at] AS sort_path
    FROM comments c
    WHERE c.id = sqlc.arg(root_id)
    UNION ALL
    SELECT
        c.*,
        ct.sort_path || c.created_at
    FROM comments c
             JOIN comment_tree ct ON c.parent_id = ct.id
)
SELECT * FROM comment_tree
ORDER BY sort_path;