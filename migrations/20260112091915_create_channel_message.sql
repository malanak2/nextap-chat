-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE "MessageChannel" (
  "ID" SERIAL PRIMARY KEY,
  "channel" INTEGER NOT NULL,
  "message" INTEGER UNIQUE NOT NULL
);


ALTER TABLE "MessageChannel" ADD FOREIGN KEY ("channel") REFERENCES "Channel" ("ID");

ALTER TABLE "MessageChannel" ADD FOREIGN KEY ("message") REFERENCES "Message" ("ID");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE "MessageChannel";
-- +goose StatementEnd
