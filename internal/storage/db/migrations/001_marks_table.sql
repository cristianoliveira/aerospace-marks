-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS marks (
    window_id INTEGER NOT NULL,
    mark TEXT NOT NULL UNIQUE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS marks;
-- +goose StatementEnd
