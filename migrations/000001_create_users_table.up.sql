-- 0001_user mifration


CREATE TABLE IF NOT EXISTS users (
    "id"         bigint       GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    "username"   varchar(255) UNIQUE NOT NULL,
    "password"   varchar(60)  NOT NULL,
    "created_at" timestamptz  NOT NULL DEFAULT CURRENT_TIMESTAMP
)
