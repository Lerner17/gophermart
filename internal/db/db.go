package db

import (
	"context"

	"github.com/Lerner17/gophermart/internal/db/postgres"
)

type DB interface {
	RegisterUser(context.Context, string, string) error
}

func GetDB() DB {
	db := postgres.New()
	return db
}
