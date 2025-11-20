-- name: CreateComment :one
INSERT INTO comments (parent_id, text)
VALUES (sqlc.arg(parent_id), sqlc.arg(text))
RETURNING id, parent_id, text, created_at, updated_at;

-- name: GetRootCommentsDesc :many
SELECT id, parent_id, text, created_at, updated_at
FROM comments
WHERE parent_id is NULL
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_) OFFSET sqlc.arg(offset_);

-- name: GetRootCommentsAsc :many
SELECT id, parent_id, text, created_at, updated_at
FROM comments
WHERE parent_id is NULL
ORDER BY created_at ASC
LIMIT sqlc.arg(limit_) OFFSET sqlc.arg(offset_);

-- name: SearchComments :many
SELECT id, parent_id, text, created_at, updated_at
FROM comments
WHERE tsv @@ to_tsquery('russian', sqlc.arg(query))
ORDER BY created_at DESC
LIMIT sqlc.arg(limit_) OFFSET sqlc.arg(offset_);

-- name: SoftDeleteComment :one
UPDATE comments
SET text = '[deleted]', updated_at = NOW(), tsv = NULL
WHERE id = sqlc.arg(id)
RETURNING id, parent_id, text, created_at, updated_at;

-- name: GetCommentSubtreeDesc :many
WITH RECURSIVE comment_tree AS (
    SELECT
        c.id, c.parent_id, c.text, c.created_at, c.updated_at,
        ARRAY[EXTRACT(EPOCH FROM c.created_at) * -1] AS sort_path
    FROM comments c
    WHERE c.id = sqlc.arg(root_id)

    UNION ALL

    SELECT
        c.id, c.parent_id, c.text, c.created_at, c.updated_at,
        ct.sort_path || (EXTRACT(EPOCH FROM c.created_at) * -1)
    FROM comments c
             JOIN comment_tree ct ON c.parent_id = ct.id
)
SELECT id, parent_id, text, created_at, updated_at FROM comment_tree
ORDER BY sort_path;

-- name: GetCommentSubtreeAsc :many
WITH RECURSIVE comment_tree AS (
    SELECT
        c.id, c.parent_id, c.text, c.created_at, c.updated_at,
        ARRAY[c.created_at] AS sort_path
    FROM comments c
    WHERE c.id = sqlc.arg(root_id)

    UNION ALL

    SELECT
        c.id, c.parent_id, c.text, c.created_at, c.updated_at,
        ct.sort_path || c.created_at
    FROM comments c
             JOIN comment_tree ct ON c.parent_id = ct.id
)
SELECT id, parent_id, text, created_at, updated_at FROM comment_tree
ORDER BY sort_path;