-- +goose Up
CREATE TABLE IF NOT EXISTS comments
(
    id         VARCHAR(255) PRIMARY KEY,
    content     VARCHAR(2000) NOT NULL,
    author_id VARCHAR(255) NOT NULL,
    post_id VARCHAR(255) NOT NULL,
    parent_comment_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX post_id_index ON comments (post_id);
CREATE INDEX parent_comment_id_index ON comments (parent_comment_id);

-- +goose Down
DROP TABLE IF EXISTS comments;
