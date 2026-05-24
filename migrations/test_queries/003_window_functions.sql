
-- rank ราคาสูงสุดของแต่ละ supplier (top 3 per supplier)
SELECT *
FROM (
    SELECT
        s.name                              AS supplier,
        p.name                              AS product,
        p.price,
        RANK() OVER (
            PARTITION BY p.supplier_id
            ORDER BY p.price DESC
        )                                   AS price_rank
    FROM products p
    JOIN suppliers s ON s.id = p.supplier_id
) ranked
WHERE price_rank <= 3
ORDER BY supplier, price_rank;

-- running total stock ของ product (เรียงตาม id)
SELECT
    id,
    name,
    stock,
    SUM(stock) OVER (ORDER BY id ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS running_total_stock
FROM products
LIMIT 20;

-- เปรียบเทียบราคาแต่ละ product กับค่าเฉลี่ยของ supplier ตัวเอง
SELECT
    p.name,
    s.name                                          AS supplier,
    p.price,
    ROUND(AVG(p.price) OVER (PARTITION BY p.supplier_id)::NUMERIC, 2) AS supplier_avg_price,
    ROUND((p.price - AVG(p.price) OVER (PARTITION BY p.supplier_id))::NUMERIC, 2) AS diff_from_avg
FROM products p
JOIN suppliers s ON s.id = p.supplier_id
ORDER BY ABS(p.price - AVG(p.price) OVER (PARTITION BY p.supplier_id)) DESC
LIMIT 20;
