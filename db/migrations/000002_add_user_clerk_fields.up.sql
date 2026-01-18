BEGIN;

-- Add clerk_id and other fields to user table
ALTER TABLE "user" ADD COLUMN "clerk_id" TEXT;
ALTER TABLE "user" ADD COLUMN "email" TEXT;
ALTER TABLE "user" ADD COLUMN "first_name" TEXT;
ALTER TABLE "user" ADD COLUMN "last_name" TEXT;
ALTER TABLE "user" ADD COLUMN "image_url" TEXT;
ALTER TABLE "user" ADD COLUMN "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP;

-- Set clerk_id and email as unique
ALTER TABLE "user" ADD CONSTRAINT "user_clerk_id_key" UNIQUE ("clerk_id");
ALTER TABLE "user" ADD CONSTRAINT "user_email_key" UNIQUE ("email");

-- Update ID to have a default value
ALTER TABLE "user" ALTER COLUMN "id" SET DEFAULT gen_random_uuid();

-- Add trigger for updated_at
CREATE TRIGGER set_updated_at_user
BEFORE UPDATE ON "user"
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Make clerk_id and email NOT NULL
-- First, ensure no NULL values exist (for production safety)
-- These placeholders should be replaced with actual values or migration should be run
-- only after ensuring all users have been backfilled via Clerk sync
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM "user" WHERE "clerk_id" IS NULL) THEN
        RAISE EXCEPTION 'Cannot set clerk_id NOT NULL: existing rows have NULL values. Please backfill data first.';
    END IF;
    IF EXISTS (SELECT 1 FROM "user" WHERE "email" IS NULL) THEN
        RAISE EXCEPTION 'Cannot set email NOT NULL: existing rows have NULL values. Please backfill data first.';
    END IF;
END $$;

ALTER TABLE "user" ALTER COLUMN "clerk_id" SET NOT NULL;
ALTER TABLE "user" ALTER COLUMN "email" SET NOT NULL;

COMMIT;
