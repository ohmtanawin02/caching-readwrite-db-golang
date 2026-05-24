-- จำนวน product ทั้งหมด แยกตาม supplier/ไม่มี supplier
SELECT
    COUNT(*)                       AS total_products,
    COUNT(supplier_id)             AS with_supplier,
    COUNT(*) - COUNT(supplier_id)  AS no_supplier
FROM products;

-- ราคา: min, max, avg, median
SELECT
    MIN(price)                          AS min_price,
    MAX(price)                          AS max_price,
    ROUND(AVG(price)::NUMERIC, 2)       AS avg_price,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY price) AS median_price
FROM products;

-- stock: distribution
SELECT
    COUNT(*) FILTER (WHERE stock = 0)           AS out_of_stock,
    COUNT(*) FILTER (WHERE stock BETWEEN 1 AND 10)  AS critical_stock,
    COUNT(*) FILTER (WHERE stock BETWEEN 11 AND 100) AS low_stock,
    COUNT(*) FILTER (WHERE stock > 100)         AS normal_stock
FROM products;
