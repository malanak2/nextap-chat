-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "Message" (
  "id" SERIAL PRIMARY KEY,
  "content" VARCHAR(1000) NOT NULL
);
-- +goose StatementEnd
--
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE "Message";
-- +goose StatementEnd
