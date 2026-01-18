BEGIN;

-- Remove trigger
DROP TRIGGER IF EXISTS set_updated_at_user ON "user";

-- Remove columns
ALTER TABLE "user" DROP COLUMN IF EXISTS "clerk_id";
ALTER TABLE "user" DROP COLUMN IF EXISTS "email";
ALTER TABLE "user" DROP COLUMN IF EXISTS "first_name";
ALTER TABLE "user" DROP COLUMN IF EXISTS "last_name";
ALTER TABLE "user" DROP COLUMN IF EXISTS "image_url";
ALTER TABLE "user" DROP COLUMN IF EXISTS "updated_at";

-- Remove default for id
ALTER TABLE "user" ALTER COLUMN "id" DROP DEFAULT;

COMMIT;
