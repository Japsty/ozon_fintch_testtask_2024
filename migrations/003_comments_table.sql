-- +goose Up
CREATE TABLE IF NOT EXISTS comments
(
    id VARCHAR(255) PRIMARY KEY,
    content VARCHAR(2000) NOT NULL,
    author_id VARCHAR(255) NOT NULL,
    post_id VARCHAR(255) NOT NULL,
    parent_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX comments_post_id_index ON comments USING HASH (post_id);
CREATE INDEX comments_parent_id_index ON comments USING HASH (parent_id);

-- +goose Down
DROP TABLE IF EXISTS comments;