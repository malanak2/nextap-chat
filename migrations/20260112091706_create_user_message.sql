-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "UserMessage" (
  "ID" SERIAL PRIMARY KEY,
  "user" INTEGER NOT NULL,
  "message" INTEGER UNIQUE NOT NULL
);

ALTER TABLE "UserMessage" ADD FOREIGN KEY ("user") REFERENCES "User" ("id");

ALTER TABLE "UserMessage" ADD FOREIGN KEY ("message") REFERENCES "Message" ("id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE "UserMessage";

-- +goose StatementEnd
