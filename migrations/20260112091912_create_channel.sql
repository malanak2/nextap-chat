-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "Channel" (
  "ID" SERIAL PRIMARY KEY,
  "name" VARCHAR(50) UNIQUE NOT NULL
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE "Channel";
-- +goose StatementEnd
