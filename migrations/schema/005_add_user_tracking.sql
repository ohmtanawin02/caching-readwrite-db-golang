ALTER TABLE products
    ADD COLUMN IF NOT EXISTS created_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS updated_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE suppliers
    ADD COLUMN IF NOT EXISTS created_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS updated_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_products_created_by_user  ON products  (created_by_user_id);
CREATE INDEX IF NOT EXISTS idx_products_updated_by_user  ON products  (updated_by_user_id);
CREATE INDEX IF NOT EXISTS idx_suppliers_created_by_user ON suppliers (created_by_user_id);
CREATE INDEX IF NOT EXISTS idx_suppliers_updated_by_user ON suppliers (updated_by_user_id);
