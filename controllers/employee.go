package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/himanshugarg165/Employee-Management/apps/employee"
	"github.com/himanshugarg165/Employee-Management/db"
	"github.com/himanshugarg165/Employee-Management/utils"
	"github.com/labstack/echo/v4"
)

// ListEmployees ...
func (controller *Controller) ListEmployees(c echo.Context) error {
	parameter, err := db.NewParameter(c)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	response, status, err := controller.EmployeeService.List(parameter)
	if err != nil {
		return c.JSON(status, utils.NewError(err))
	}
	return c.JSON(status, response)
}

// GetEmployee ...
func (controller *Controller) GetEmployee(c echo.Context) error {
	employeeID := c.Param("id")
	intEmployeeID, err := strconv.ParseInt(employeeID, 10, 64)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	response, status, err := controller.EmployeeService.GetByID(intEmployeeID)
	if err != nil {
		return c.JSON(status, utils.NewError(err))
	}
	return c.JSON(status, response)
}

// CreateEmployee ...
func (controller *Controller) CreateEmployee(c echo.Context) error {
	requestData := employee.EmployeeCreateRequest{}

	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewValidatorError(err))
	}
	requestData.Name = strings.TrimSpace(requestData.Name)
	requestData.DOB = strings.TrimSpace(requestData.DOB)

	if err := c.Validate(&requestData); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewValidatorError(err))
	}

	response, status, err := controller.EmployeeService.Create(&requestData)
	if err != nil {
		return c.JSON(status, utils.NewError(err))
	}
	return c.JSON(status, response)
}

// UpdateEmployee ...
func (controller *Controller) UpdateEmployee(c echo.Context) error {
	employeeID := c.Param("id")
	intEmployeeID, err := strconv.ParseInt(employeeID, 10, 64)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	requestData := employee.EmployeeUpdateRequest{}
	if err := c.Bind(&requestData); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewValidatorError(err))
	}
	requestData.Name = strings.TrimSpace(requestData.Name)
	requestData.DOB = strings.TrimSpace(requestData.DOB)

	if err := c.Validate(&requestData); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewValidatorError(err))
	}
	requestData.ID = intEmployeeID
	response, status, err := controller.EmployeeService.Update(&requestData)
	if err != nil {
		return c.JSON(status, utils.NewError(err))
	}
	return c.JSON(status, &response)
}
