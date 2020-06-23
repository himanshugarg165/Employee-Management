package controllers

import (
	"github.com/himanshugarg165/Employee-Management/apps/employee"
	"github.com/labstack/echo/v4"
)

// Controller ...
type Controller struct {
	EmployeeService *employee.EmployeeService
}

// NewController ...
func NewController() *Controller {
	return &Controller{}
}

// Register registers all the API routes
func (controller *Controller) Register(v1 *echo.Group) {
	employees := v1.Group("/employee")
	employees.POST("", controller.CreateEmployee)
	employees.GET("", controller.ListEmployees)
	employees.GET("/:id", controller.GetEmployee)
	employees.PUT("/:id", controller.UpdateEmployee)
}
