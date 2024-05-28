-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    user_id VARCHAR(255) NOT NULL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    );

CREATE INDEX id_index ON posts USING btree (user_id);

-- +goose Down
DROP TABLE IF EXISTS users;