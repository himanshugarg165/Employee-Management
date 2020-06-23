package server

import (
	"time"

	"github.com/himanshugarg165/Employee-Management/constants"
	"github.com/labstack/gommon/log"
	"gopkg.in/go-playground/validator.v9"
)

func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func (v *Validator) Register() error {
	return v.validator.RegisterValidation("isValidDOB", isValidDOB)
}

func isValidDOB(fl validator.FieldLevel) bool {
	inputDOB := fl.Field().Interface()
	dob, err := time.Parse(constants.DOB_FORMAT, inputDOB.(string))
	if err != nil {
		log.Infof("Invalid DOB input: %s", inputDOB)
		return false
	}
	// Add check for DOB to be more than 20 years if not throw error
	if dob.After(time.Now().AddDate(-20, 0, 0)) {
		log.Infof("DOB should be 20 years before the current date")
		return false
	}
	return true
}
