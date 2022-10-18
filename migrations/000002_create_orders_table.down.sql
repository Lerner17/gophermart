-- 0002_order mifration

DELETE TYPE IF EXISTS order_status;

DROP TABLE IF EXISTS orders;