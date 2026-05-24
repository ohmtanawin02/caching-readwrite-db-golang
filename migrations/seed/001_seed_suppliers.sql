-- Seed 1,000 suppliers
INSERT INTO suppliers (name, created_at, updated_at)
SELECT
    'Supplier ' || i || ' - ' || (
        ARRAY['TH','US','JP','CN','DE','KR','SG','IN','GB','FR']
    )[(MOD(i, 10) + 1)],
    NOW() - (RANDOM() * INTERVAL '730 days'),
    NOW()
FROM generate_series(1, 1000) AS i
ON CONFLICT DO NOTHING;
