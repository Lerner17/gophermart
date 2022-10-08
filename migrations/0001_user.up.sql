-- 0001_user mifration


CREATE TABLE IF NOT EXISTS `users` (
    `id` serial NOT NULL,
    `username` varchar(255) UNIQUE NOT NULL,
    `password` varchar(60) NOT NULL,
    `created_at` timestamp NOT NULL DEFAULT now()
)
