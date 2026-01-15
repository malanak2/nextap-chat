-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE "User" (
  "ID" SERIAL PRIMARY KEY,
  "username" VARCHAR(50) UNIQUE NOT NULL
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE "User";
-- +goose StatementEnd
