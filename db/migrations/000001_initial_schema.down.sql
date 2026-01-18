BEGIN;

-- Drop triggers
DROP TRIGGER IF EXISTS set_updated_at_expense ON "expense";
DROP TRIGGER IF EXISTS set_updated_at_category ON "category";

-- Drop the function
DROP FUNCTION IF EXISTS update_updated_at_column;

-- Drop foreign key constraints
ALTER TABLE "expense" DROP CONSTRAINT IF EXISTS "expense_user_id_fkey";
ALTER TABLE "expense" DROP CONSTRAINT IF EXISTS "expense_category_id_fkey";
ALTER TABLE "category" DROP CONSTRAINT IF EXISTS "category_user_id_fkey";

-- Drop tables
DROP TABLE IF EXISTS "expense";
DROP TABLE IF EXISTS "category";
DROP TABLE IF EXISTS "user";

COMMIT;
