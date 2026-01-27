-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE "Message" ADD COLUMN "user_id" INTEGER NOT NULL;
UPDATE "Message"
SET
    "user_id" = "UserMessage"."user"
FROM "UserMessage"
WHERE "Message".id = "UserMessage".message;

ALTER TABLE "Message" ADD FOREIGN KEY ("user_id") REFERENCES "User" ("id");

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';


-- +goose StatementEnd
