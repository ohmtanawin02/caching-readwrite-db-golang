-- Seed 1,000,000 products
--   700,000 records → มี supplier_id  (random จาก suppliers ที่มีอยู่)
--   300,000 records → supplier_id = NULL (ไม่มี relation)
--
-- ใช้ generate_series + batch insert 10,000 ต่อครั้ง เพื่อ performance
-- ประมาณเวลา: ~30-60 วินาทีขึ้นกับ machine

DO $$
DECLARE
    supplier_count BIGINT;
BEGIN
    SELECT COUNT(*) INTO supplier_count FROM suppliers;
    IF supplier_count = 0 THEN
        RAISE EXCEPTION 'Run seed suppliers first (001_seed_suppliers.sql)';
    END IF;
END $$;

-- Batch 1: 700,000 records ที่มี supplier
INSERT INTO products (name, price, stock, supplier_id, created_at, updated_at)
SELECT
    (
        ARRAY[
            'Ultra','Pro','Max','Lite','Plus','Air','Neo','Smart','Eco','Fast'
        ]
    )[(MOD(i, 10) + 1)]
    || ' '
    || (
        ARRAY[
            'Widget','Gadget','Device','Module','Component',
            'Unit','Kit','Pack','Set','Bundle'
        ]
    )[(MOD(i, 10) + 1)]
    || ' #' || i,

    ROUND((RANDOM() * 99900 + 100)::NUMERIC, 2),   -- price: 100 - 100,000
    FLOOR(RANDOM() * 10000)::INT,                   -- stock: 0 - 9,999

    (FLOOR(RANDOM() * (SELECT COUNT(*) FROM suppliers)) + 1)::BIGINT,

    NOW() - (RANDOM() * INTERVAL '365 days'),
    NOW()
FROM generate_series(1, 700000) AS i;

-- Batch 2: 300,000 records ที่ไม่มี supplier (supplier_id = NULL)
INSERT INTO products (name, price, stock, supplier_id, created_at, updated_at)
SELECT
    'Unbranded '
    || (
        ARRAY[
            'Item','Product','Good','Ware','Article',
            'Commodity','Material','Supply','Stock','Part'
        ]
    )[(MOD(i, 10) + 1)]
    || ' #' || i,

    ROUND((RANDOM() * 4900 + 100)::NUMERIC, 2),    -- price: 100 - 5,000 (cheaper unbranded)
    FLOOR(RANDOM() * 500)::INT,                     -- stock: 0 - 499

    NULL,                                           -- ไม่มี supplier

    NOW() - (RANDOM() * INTERVAL '365 days'),
    NOW()
FROM generate_series(700001, 1000000) AS i;
