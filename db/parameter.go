package db

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	defaultLimit = 20
	defaultPage  = 1
)

// Parameter can be used to add support for pagination, sorting and filtering in db queries
type Parameter struct {
	limit int
	page  int
}

// NewParameter ...
func NewParameter(context echo.Context) (*Parameter, error) {
	parameter := &Parameter{}

	if err := parameter.initialize(context); err != nil {
		return nil, err
	}

	return parameter, nil
}

func (parameter *Parameter) initialize(context echo.Context) error {
	limit, err := validate(context.QueryParam("perPage"))
	if err != nil {
		return fmt.Errorf("Invalid per page value: %s", err.Error())
	}
	if limit < 0 {
		limit = defaultLimit
	}
	page, err := validate(context.QueryParam("page"))
	if err != nil {
		return fmt.Errorf("Invalid page value: %s", err.Error())
	}
	if page <= 0 {
		page = defaultPage
	}

	parameter.limit = int(math.Max(0, math.Min(10000, float64(limit))))

	parameter.page = int(math.Max(1, float64(page)))

	return nil
}

func validate(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return -1, nil
	}

	num, err := strconv.Atoi(s)
	if err != nil {
		return -1, err
	}

	return num, nil
}

// Offset ...
func (parameter Parameter) Offset() int {
	return parameter.limit * (parameter.page - 1)
}

// Limit ...
func (parameter *Parameter) Limit() int {
	return parameter.limit
}
