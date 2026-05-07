-- Hapus index terlebih dahulu
DROP INDEX IF EXISTS idx_orders_secure_id;

-- Hapus column, best practice dari child ke parent
ALTER TABLE order_items DROP COLUMN IF EXISTS secure_id;
ALTER TABLE orders DROP COLUMN IF EXISTS secure_id;