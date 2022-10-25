package db

import (
	"context"

	"github.com/Lerner17/gophermart/internal/db/postgres"
	"github.com/Lerner17/gophermart/internal/models"
)

type DB interface {
	RegisterUser(context.Context, string, string) error
	LoginUser(string, string) (int, error)

	CreateOrder(context.Context, models.Order) error
	GetOrders(context.Context, int) ([]models.Order, error)
	GetUserBalance(context.Context, int) (models.Balance, error)

	CreateTransaction(context.Context, int, string, int) error
}

func GetDB() DB {
	db := postgres.New()
	return db
}
