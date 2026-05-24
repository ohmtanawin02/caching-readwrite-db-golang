-- query ที่ใช้ index บน supplier_id
EXPLAIN ANALYZE
SELECT p.id, p.name, p.price, s.name AS supplier
FROM products p
JOIN suppliers s ON s.id = p.supplier_id
WHERE p.supplier_id = 42
ORDER BY p.price DESC
LIMIT 10;

-- query filter ราคา (ไม่มี index — จะเห็น Seq Scan)
EXPLAIN ANALYZE
SELECT id, name, price
FROM products
WHERE price BETWEEN 5000 AND 10000
ORDER BY price
LIMIT 20;

-- เพิ่ม index บน price แล้วดูผลต่าง
-- CREATE INDEX idx_products_price ON products(price);

-- EXPLAIN ANALYZE
-- SELECT id, name, price
-- FROM products
-- WHERE price BETWEEN 5000 AND 10000
-- ORDER BY price
-- LIMIT 20;
-- (รัน query เดิมอีกครั้งหลัง create index แล้วเปรียบ execution time)

-- ดู index ทั้งหมดในตาราง products
SELECT
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'products';
