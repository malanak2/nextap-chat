-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "MessageChannel" (
  "id" SERIAL PRIMARY KEY,
  "channel_id" INTEGER NOT NULL,
  "message_id" INTEGER UNIQUE NOT NULL
);


ALTER TABLE "MessageChannel" ADD FOREIGN KEY ("channel_id") REFERENCES "Channel" ("id");

ALTER TABLE "MessageChannel" ADD FOREIGN KEY ("message_id") REFERENCES "Message" ("id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS "MessageChannel";
-- +goose StatementEnd
