-- +goose Up
-- +goose StatementBegin
INSERT INTO posts (id, post_id, content, created_at)
VALUES
    (1, "f55a0da4-0225-11ee-be56-0242ac120002", "This is my first post. Hello all!", "2022-12-24 10:34:23"),
    (2, "f8a434d0-0225-11ee-be56-0242ac120002", "Yet another post.", "2022-12-25 09:12:53");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM posts WHERE id IN (1, 2);
-- +goose StatementEnd