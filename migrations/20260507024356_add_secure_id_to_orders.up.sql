-- Tambahkan UUID ke table orders
ALTER TABLE orders ADD COLUMN secure_id UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE;

-- Buat index agar pencarian via api url sangat cepat
CREATE INDEX idx_orders_secure_id ON orders(secure_id);

-- Tambahkan UUID ke table order_items
ALTER TABLE order_items ADD COLUMN secure_id UUID DEFAULT gen_random_uuid() NOT NULL UNIQUE;