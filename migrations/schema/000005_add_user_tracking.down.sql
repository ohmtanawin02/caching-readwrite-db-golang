DROP INDEX IF EXISTS idx_suppliers_updated_by_user;
DROP INDEX IF EXISTS idx_suppliers_created_by_user;
DROP INDEX IF EXISTS idx_products_updated_by_user;
DROP INDEX IF EXISTS idx_products_created_by_user;

ALTER TABLE suppliers
    DROP COLUMN IF EXISTS created_by_user_id,
    DROP COLUMN IF EXISTS updated_by_user_id;

ALTER TABLE products
    DROP COLUMN IF EXISTS created_by_user_id,
    DROP COLUMN IF EXISTS updated_by_user_id;
