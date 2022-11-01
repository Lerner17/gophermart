
CREATE TABLE IF NOT EXISTS transactions (
    "id"          bigint       GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "user_id"     bigint       NOT NULL,
    "order_id"    bigint       NOT NULL,
    "amount"      numeric      NOT NULL
)
