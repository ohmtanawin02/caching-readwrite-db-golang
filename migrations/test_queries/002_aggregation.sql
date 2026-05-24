-- Top 10 supplier ที่มี product เยอะสุด
SELECT
    s.name                          AS supplier,
    COUNT(p.id)                     AS product_count,
    ROUND(AVG(p.price)::NUMERIC, 2) AS avg_price,
    SUM(p.stock)                    AS total_stock
FROM suppliers s
JOIN products p ON p.supplier_id = s.id
GROUP BY s.id, s.name
ORDER BY product_count DESC
LIMIT 10;

-- แบ่ง product เป็น price tier
SELECT
    CASE
        WHEN price < 1000   THEN '< 1,000'
        WHEN price < 10000  THEN '1,000 - 9,999'
        WHEN price < 50000  THEN '10,000 - 49,999'
        ELSE                     '>= 50,000'
    END                     AS price_tier,
    COUNT(*)                AS product_count,
    ROUND(AVG(price)::NUMERIC, 2) AS avg_price
FROM products
GROUP BY price_tier
ORDER BY MIN(price);

-- supplier ที่มี product มากกว่า 800 ชิ้น (HAVING)
SELECT
    s.name,
    COUNT(p.id) AS product_count
FROM suppliers s
JOIN products p ON p.supplier_id = s.id
GROUP BY s.id, s.name
HAVING COUNT(p.id) > 760
ORDER BY product_count DESC;
