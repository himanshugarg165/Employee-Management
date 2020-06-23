package server

import (
	"github.com/himanshugarg165/Employee-Management/apps/employee"
	"github.com/himanshugarg165/Employee-Management/cache"
	"github.com/himanshugarg165/Employee-Management/controllers"
	"github.com/himanshugarg165/Employee-Management/db"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func New() *echo.Echo {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Validator = NewValidator()
	e.Validator.(*Validator).Register()

	sqlDB := db.New()
	db.Migrate(sqlDB)
	e.Use(db.DBMiddleware(sqlDB))

	employeeService := employee.NewEmployeeService(sqlDB)
	ctrl := controllers.NewController()
	ctrl.EmployeeService = employeeService
	ctrl.Register(e.Group("/api"))

	var err error
	cache.Instance, err = cache.NewRedisClient(cache.CacheConfig{URL: "localhost:6379", Password: ""})
	if err != nil {
		panic(err)
	}
	return e
}
