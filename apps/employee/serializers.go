package employee

import "time"

// Employee ...
type Employee struct {
	ID   int64     `json:"id"`
	Name string    `json:"name"`
	DOB  time.Time `json:"dob"`
}

// EmployeeList ...
type EmployeeList struct {
	Employees []Employee `json:"employees"`
	Count     int64      `json:"totalCount"`
}

// EmployeeCreateRequest ...
type EmployeeCreateRequest struct {
	Name string `json:"name" validate:"required"`
	DOB  string `json:"dob" validate:"required,isValidDOB"`
}

// EmployeeUpdateRequest ...
type EmployeeUpdateRequest struct {
	ID   int64
	Name string `json:"name" validate:"required"`
	DOB  string `json:"dob" validate:"required,isValidDOB"`
}
