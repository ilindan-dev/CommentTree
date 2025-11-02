-- +goose UP

-- comments table to store hierarchical comments with full-text search capabilities
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id UUID REFERENCES comments(id),
    text TEXT NOT NULL CHECK ( length(text) > 0 ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    tsv TSVECTOR
);

-- Indexes for performance optimization

-- Index on parent_id for efficient retrieval of child comments
CREATE INDEX idx_comments_parent_id ON comments(parent_id);

-- Index on created_at for efficient sorting by creation time
CREATE INDEX idx_comments_created_at ON comments(created_at);

-- Full-text search index on the tsv column
CREATE INDEX comments_tsv_idx ON comments USING GIN(tsv);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION comments_tsv_trigger() RETURNS TRIGGER AS $$
BEGIN
    NEW.tsv := to_tsvector('russian', coalesce(NEW.text, ''));
    RETURN NEW;
END
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
ON comments FOR EACH ROW EXECUTE PROCEDURE comments_tsv_trigger();

-- +goose DOWN
DROP TRIGGER IF EXISTS tsvectorupdate ON comments;
DROP FUNCTION IF EXISTS comments_tsv_trigger();
DROP TABLE IF EXISTS comments;