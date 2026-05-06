CREATE TABLE IF NOT EXISTS products (
    id           UUID PRIMARY KEY,
    name         TEXT NOT NULL,
    description  TEXT NOT NULL DEFAULT '',
    price_cents  BIGINT NOT NULL CHECK (price_cents > 0),
    status       TEXT NOT NULL CHECK (status IN ('ACTIVE', 'INACTIVE')),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_created_at ON products (created_at DESC);
