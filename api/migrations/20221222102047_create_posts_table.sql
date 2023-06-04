-- +goose Up
-- +goose StatementBegin
CREATE TABLE `posts` (
                        `id`        int NOT NULL AUTO_INCREMENT,
                        `post_id`   varchar(45) DEFAULT NULL,
                        `content`   varchar(255) DEFAULT NULL,
                        created_at  datetime(3) NULL,
                        updated_at  datetime(3) NULL,
                        deleted_at  datetime(3) NULL,
                        PRIMARY KEY (`id`)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd