DROP INDEX IF EXISTS idx_products_status;
DROP INDEX IF EXISTS idx_products_deleted_at;
ALTER TABLE products
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS deleted_at;

DROP INDEX IF EXISTS idx_suppliers_status;
DROP INDEX IF EXISTS idx_suppliers_deleted_at;
ALTER TABLE suppliers
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS deleted_at;
