-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE changelogs (
    id SERIAL PRIMARY KEY,
    version VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INT NOT NULL REFERENCES users(id)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS changelogs;
