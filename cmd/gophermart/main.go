package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Lerner17/gophermart/internal/db"
	"github.com/Lerner17/gophermart/internal/handlers"
	"github.com/Lerner17/gophermart/internal/models"
	"github.com/labstack/echo/v4"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5"
)

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

	// // FIXME: Do not hardcode sql dependency
	// if errors.Is(err, sql.ErrNoRows) {
	// 	code = http.StatusNotFound
	// }

	// FIXME: Do not hardcode postgres dependency
	// var pqError pq.Error
	// if errors.As(err, &pqError) {
	// 	const pgConstraintViolationError = "23505"
	// 	if pqError.Code == pgConstraintViolationError {
	// 		code = http.StatusBadRequest
	// 	}
	// }

	var msgerr interface{ Message() string }
	if errors.As(err, &msgerr) {
		message = msgerr.Message()
	}

	// var vErr core.ValidationError
	// if errors.As(err, &vErr) {
	// 	response = vErr
	// }

	// unknown error
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

func main() {
	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler
	db := db.GetDB()

	migragte(e) // Migrate migrations

	e.POST("/api/user/register", handlers.Registration(db))
	e.Logger.Fatal(e.Start(":5000"))
}

func migragte(e *echo.Echo) {
	db, err := sql.Open("postgres", "postgres://shroten:shroten@localhost:5432/shroten?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	{ // DB Migrations
		const MigrationVersion = 1
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
			// e.Logger.Infof("detected migration version missmatch current [%v] but need [%v]. Start migration...", ver, MigrationVersion)
			if err := m.Migrate(MigrationVersion); err != nil {
				panic(fmt.Errorf("could not apply migrations: %w", err))
			}
		}
	}
}
