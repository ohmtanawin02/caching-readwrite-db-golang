-- CTE: หา supplier ที่มี avg price สูงสุด แล้วดึง product ของ supplier นั้น
WITH top_supplier AS (
    SELECT
        supplier_id,
        ROUND(AVG(price)::NUMERIC, 2) AS avg_price
    FROM products
    WHERE supplier_id IS NOT NULL
    GROUP BY supplier_id
    ORDER BY avg_price DESC
    LIMIT 1
)
SELECT
    p.id,
    p.name,
    p.price,
    p.stock,
    s.name AS supplier,
    ts.avg_price AS supplier_avg_price
FROM products p
JOIN top_supplier ts ON ts.supplier_id = p.supplier_id
JOIN suppliers s ON s.id = p.supplier_id
ORDER BY p.price DESC
LIMIT 20;

-- CTE: stock alert — product ที่ stock ต่ำกว่า avg ของ supplier ตัวเอง
WITH supplier_avg_stock AS (
    SELECT
        supplier_id,
        AVG(stock) AS avg_stock
    FROM products
    WHERE supplier_id IS NOT NULL
    GROUP BY supplier_id
)
SELECT
    p.name,
    p.stock,
    ROUND(sas.avg_stock::NUMERIC, 0) AS supplier_avg_stock,
    s.name AS supplier
FROM products p
JOIN supplier_avg_stock sas ON sas.supplier_id = p.supplier_id
JOIN suppliers s ON s.id = p.supplier_id
WHERE p.stock < sas.avg_stock * 0.1   -- stock ต่ำกว่า 10% ของ avg supplier
ORDER BY p.stock ASC
LIMIT 20;

-- Subquery: product ที่ราคาสูงกว่า avg ราคาทั้งหมด
SELECT id, name, price
FROM products
WHERE price > (SELECT AVG(price) FROM products)
ORDER BY price DESC
LIMIT 20;
