-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE "Message" ADD COLUMN "user_id" INTEGER;
UPDATE public."Message"
SET "user_id" = um."user"
    FROM public."UserMessage" AS um
WHERE public."Message"."id" = um."message";

ALTER TABLE "Message" ADD CONSTRAINT "fk_message_user" FOREIGN KEY ("user_id") REFERENCES "User" ("id");

ALTER TABLE "Message" ALTER COLUMN "user_id" SET NOT NULL;

DROP TABLE "UserMessage";

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
ALTER TABLE "Message" DROP CONSTRAINT IF EXISTS "fk_message_user";
ALTER TABLE "Message" DROP COLUMN IF EXISTS "user_id";

-- +goose StatementEnd
