-- 0002_order mifration

DO $$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
        CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
    END IF;
END
$$;

CREATE TABLE IF NOT EXISTS orders (
    "id"          bigint       GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "number"      bigint       UNIQUE NOT NULL,
    "user_id"     bigint       NOT NULL,
    "created_at"  timestamptz  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "uploaded_at" timestamptz  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "status"      order_status NOT NULL DEFAULT 'NEW'
)
