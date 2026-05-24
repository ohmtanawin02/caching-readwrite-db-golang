ALTER TABLE suppliers
    ADD COLUMN IF NOT EXISTS status     VARCHAR(20)  NOT NULL DEFAULT 'active',
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_suppliers_deleted_at ON suppliers (deleted_at);
CREATE INDEX IF NOT EXISTS idx_suppliers_status     ON suppliers (status);

ALTER TABLE products
    ADD COLUMN IF NOT EXISTS status     VARCHAR(20)  NOT NULL DEFAULT 'active',
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_products_deleted_at ON products (deleted_at);
CREATE INDEX IF NOT EXISTS idx_products_status     ON products (status);
