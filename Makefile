# ==========================================
# Database Migration (golang-migrate)
# schema ทั้งหมดจัดการผ่าน migrate — ไม่มี auto-create บน docker up อีกแล้ว
# รันที่ PRIMARY (9191) เสมอ — replica (9192) เป็น read-only
# ติดตั้ง CLI: brew install golang-migrate
# ==========================================

MIGRATE_PATH := migrations/schema
DB_URL       := postgres://postgres:postgres@localhost:9191/product_db?sslmode=disable

# seed (รัน manual ผ่าน psql — ไม่ใช่ส่วนของ migrate)
PGPASSWORD   := postgres
PG_HOST      := localhost
PG_PORT      := 9191
PG_USER      := postgres
PG_DB        := product_db
PSQL         := PGPASSWORD=$(PGPASSWORD) psql -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d $(PG_DB)

.PHONY: migrate-up migrate-down migrate-down-all migrate-version migrate-force migrate-create migrate-drop seed

## รัน migration ที่ค้างทั้งหมด
migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up

## rollback 1 step (เช่น make migrate-down n=2 เพื่อ rollback 2 step)
migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down $(or $(n),1)

## rollback ทั้งหมด (ถาม y/n)
migrate-down-all:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down

## ดู version ปัจจุบัน
migrate-version:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" version

## แก้ dirty state ให้เป็น version v (เช่น make migrate-force v=5)
migrate-force:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" force $(v)

## สร้าง migration ใหม่ (เช่น make migrate-create name=create_orders)
migrate-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq $(name)

## ลบทุก object ใน DB (อันตราย — dev เท่านั้น)
migrate-drop:
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" drop -f

## seed ข้อมูล (suppliers ก่อน แล้วตามด้วย products — เพราะ products มี FK ไป suppliers)
seed:
	$(PSQL) -f migrations/seed/001_seed_suppliers.sql
	$(PSQL) -f migrations/seed/002_seed_products.sql
