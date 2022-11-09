package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Lerner17/gophermart/internal/auth"
	"github.com/Lerner17/gophermart/internal/config"
	"github.com/Lerner17/gophermart/internal/consumer"
	"github.com/Lerner17/gophermart/internal/db"
	"github.com/Lerner17/gophermart/internal/handlers"
	"github.com/Lerner17/gophermart/internal/models"
	"github.com/Lerner17/gophermart/internal/queue"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func parsArgs(c *config.Config) {
	serverAddressPtr := flag.String("a", "", "")
	DatabaseDsnPtr := flag.String("d", "", "")
	AccrualSystemAddressPtr := flag.String("r", "", "")
	flag.Parse()

	if *serverAddressPtr != "" {
		c.ServerAddress = *serverAddressPtr
	}

	if *AccrualSystemAddressPtr != "" {
		c.AccrualSystemAddress = *AccrualSystemAddressPtr
	}

	if *DatabaseDsnPtr != "" {
		c.DatabaseDsn = *DatabaseDsnPtr
	}
}

func customHTTPErrorHandler(err error, ctx echo.Context) {
	var code = http.StatusInternalServerError
	var response interface{}
	var message string

	if err, ok := err.(*echo.HTTPError); ok {
		code = err.Code
	}

	var codeerr interface{ HTTPCode() int }
	if errors.As(err, &codeerr) {
		code = codeerr.HTTPCode()
	}

	var msgerr interface{ Message() string }
	if errors.As(err, &msgerr) {
		message = msgerr.Message()
	}

	ctx.Logger().Error(err)

	if err := ctx.JSON(code, models.Response{
		Code:     code,
		Success:  false,
		Message:  message,
		Response: response,
	}); err != nil {
		ctx.Logger().Error(err)
	}
}

func migragte(e *echo.Echo, cfg *config.Config) {
	db, err := sql.Open("postgres", cfg.DatabaseDsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	{ // DB Migrations
		const MigrationVersion = 3
		mDriver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			panic(fmt.Errorf("could not instantiate db instance for migrations: %w", err))
		}
		m, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			"postgres", mDriver,
		)
		if err != nil {
			panic(fmt.Errorf("could not instantiate migrate instance for migrations: %w", err))
		}
		ver, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			panic(fmt.Errorf("could not get current migration version: %w", err))
		}
		if dirty {
			panic("detected dirty migration, please resolve it manually")
		}
		if MigrationVersion != ver {
			if err := m.Migrate(MigrationVersion); err != nil {
				panic(fmt.Errorf("could not apply migrations: %w", err))
			}
		}
	}
}

func main() {
	cfg := config.Instance
	parsArgs(cfg)
	fmt.Println(cfg)
	e := echo.New()
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.HTTPErrorHandler = customHTTPErrorHandler
	db := db.GetDB()

	migragte(e, cfg) // Migrate migrations
	authGroup := e.Group("")
	authGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:                  &models.JwtCustomClaims{},
		SigningKey:              []byte(cfg.JWTSecretKey),
		TokenLookup:             "cookie:access-token",
		ErrorHandlerWithContext: auth.JWTErrorChecker,
	}))
	e.POST("/api/user/register", handlers.Registration(db))
	e.POST("/api/user/login", handlers.LoginHandler(db))

	authGroup.POST("/api/user/orders", handlers.CreateOrderHandler(db))
	authGroup.GET("/api/user/orders", handlers.GetOrdersHandler(db))

	authGroup.GET("/api/user/balance", handlers.BalanceHandler(db))

	authGroup.POST("/api/user/balance/withdraw", handlers.WithdrawHandler(db))
	authGroup.GET("/api/user/withdrawals", handlers.GetWithdrawListHandler(db))

	orders, err := RestoreQueue(db)
	if err != nil {
		e.Logger.Infof("Could not restore orders queue dump: %v", err)
	} else {
		queue.FullfilQueue(orders)
	}

	go func() {
		if err := e.Start(cfg.ServerAddress); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	go func() {
		consumer.ProcessOrderBounce(e.Logger, db)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	e.Logger.Info("Received an interrupt, stopping gophermart…")

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

type DBRestorer interface {
	GetNewOrders(context.Context) ([]models.Order, error)
}

func RestoreQueue(db DBRestorer) ([]models.Order, error) {

	ctx := context.Background()

	results, err := db.GetNewOrders(ctx)
	if err != nil {
		panic(err)
	}

	return results, nil
}
