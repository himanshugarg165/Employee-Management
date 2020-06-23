package db

import (
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
)

const (
	dsnFormat = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	host      = "localhost"
	port      = 5432
	user      = "postgres"
	password  = "password"
	dbname    = "employee_db"
)

func New() *sql.DB {
	dbDSN := fmt.Sprintf(dsnFormat, host, port, user, password, dbname)
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	log.Infof(`Successfully connected using DSN: "%s"`, dbDSN)
	return db
}

// DBMiddleware ...
func DBMiddleware(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("DB", db)
			return next(c)
		}
	}
}

// Migrate ...
func Migrate(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS employees (
		ID  SERIAL PRIMARY KEY,
		NAME varchar(100) NOT NULL,
		DOB timestamp with time zone NOT NULL,
		ACTIVE boolean NOT NULL
	);`)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	_, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_active_employee on employees(active);")
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		panic(err)
	}
	log.Infof("Successfully migrated the database!!")
}
