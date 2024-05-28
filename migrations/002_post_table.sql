-- +goose Up
CREATE TABLE IF NOT EXISTS posts
(
    id         VARCHAR(255) PRIMARY KEY ,
    title       VARCHAR(255) NOT NULL,
    content     VARCHAR(2000) NOT NULL,
    commentAllowed BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX id_index ON posts USING btree (id);

-- +goose Down
DROP TABLE IF EXISTS posts;
