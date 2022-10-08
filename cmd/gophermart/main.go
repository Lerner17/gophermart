package main

import (
	"github.com/Lerner17/gophermart/internal/db"
	"github.com/Lerner17/gophermart/internal/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	db := db.GetDB()

	e.POST("/api/user/register", handlers.Registration(db))
	e.Logger.Fatal(e.Start(":5000"))
}
