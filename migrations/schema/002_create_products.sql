CREATE TABLE IF NOT EXISTS products (
    id          BIGSERIAL      PRIMARY KEY,
    name        VARCHAR(255)   NOT NULL,
    price       NUMERIC(12, 2) NOT NULL,
    stock       INT            NOT NULL DEFAULT 0,
    supplier_id BIGINT         REFERENCES suppliers (id) ON DELETE SET NULL,  -- nullable: ไม่มี relation ได้
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_supplier_id ON products (supplier_id);
CREATE INDEX IF NOT EXISTS idx_products_id_cursor   ON products (id);          -- สำหรับ cursor pagination
