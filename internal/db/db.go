package db

import "github.com/Lerner17/gophermart/internal/db/postgres"

type DB interface {
	RegisterUser(string, string) error
}

func GetDB() DB {
	db := postgres.New()
	return db
}
