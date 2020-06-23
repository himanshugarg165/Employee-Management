package employee

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/himanshugarg165/Employee-Management/cache"
	"github.com/himanshugarg165/Employee-Management/constants"
	"github.com/himanshugarg165/Employee-Management/db"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
)

var cacheTimeout time.Duration = 60 * time.Second

// GetEmployeeCacheKey returns the cache key
func GetEmployeeCacheKey(employeeID int64) string {
	return fmt.Sprintf("emp:%d", employeeID)
}

// EmployeeService ...
type EmployeeService struct {
	db *sql.DB
}

// NewEmployeeService ...
func NewEmployeeService(db *sql.DB) *EmployeeService {
	return &EmployeeService{
		db: db,
	}
}

// List ...
func (employeeService *EmployeeService) List(parameter *db.Parameter) (*EmployeeList, int, error) {
	// 1. Fetch the total number of the employees
	var count int64
	query := "SELECT count(id) FROM employees where active = $1"
	err := employeeService.db.QueryRow(query, true).Scan(&count)
	if err != nil {
		log.Errorf("Error while fetching the total employees count")
		return nil, http.StatusInternalServerError, err
	}
	employees := make([]Employee, 0)

	// 2. Fetch the list of ids of the requested employees
	query = "SELECT id FROM employees where active = $1 ORDER BY ID"

	// TODO: Just a workaround for pagination, need to look into better options
	limit := parameter.Limit()
	if limit >= 0 {
		query = fmt.Sprintf("%s LIMIT %d", query, limit)
	}
	offset := parameter.Offset()
	if offset >= 0 {
		query = fmt.Sprintf("%s OFFSET %d", query, offset)
	}
	rows, err := employeeService.db.Query(query, true)
	if err != nil {
		log.Errorf("Error while fetching the ids of the request employees")
		return nil, http.StatusInternalServerError, err
	}
	// 3. Iterate over the result set and list all the ids which are not found in the cache.
	empIDToEmployeeMap := make(map[int64]Employee)
	var employeeIDs []int64
	var missingEmployeeIDs []int64
	for rows.Next() {
		var employee Employee
		var employeeID int64
		if err := rows.Scan(&employeeID); err != nil {
			log.Errorf("Error while scanning the current employee ID in the result set")
			return nil, http.StatusInternalServerError, err
		}
		employeeIDs = append(employeeIDs, employeeID)
		err := cache.Instance.Get(GetEmployeeCacheKey(employee.ID), &employee)
		if err != nil {
			if err != cache.ErrCacheMiss {
				log.Errorf("Error while accessing the cache to get an employee")
				return nil, http.StatusInternalServerError, err
			}
			missingEmployeeIDs = append(missingEmployeeIDs, employeeID)
		} else {
			empIDToEmployeeMap[employeeID] = employee
		}
	}
	rows.Close()

	// 4. Fetch the details of the employees which were not present in the cache.
	if len(missingEmployeeIDs) > 0 {
		query = "SELECT id, name, dob FROM employees where id = ANY($1)"
		rows, err := employeeService.db.Query(query, pq.Int64Array(missingEmployeeIDs))
		if err != nil {
			log.Errorf("Error while fetching the details of the employees not found in the cache")
			return nil, http.StatusInternalServerError, err
		}
		defer rows.Close()
		for rows.Next() {
			var employee Employee
			if err := rows.Scan(&employee.ID, &employee.Name, &employee.DOB); err != nil {
				log.Errorf("Error while scanning the current employee details in the result set")
				return nil, http.StatusInternalServerError, err
			}
			err = cache.Instance.Set(GetEmployeeCacheKey(employee.ID), employee, cacheTimeout)
			if err != nil {
				log.Errorf("Error while writing the current employee to the cache")
				return nil, http.StatusInternalServerError, err
			}
			empIDToEmployeeMap[employee.ID] = employee
		}
	}
	for _, employeeID := range employeeIDs {
		employees = append(employees, empIDToEmployeeMap[employeeID])
	}
	return &EmployeeList{Employees: employees, Count: count}, http.StatusOK, nil
}

// GetByID ...
func (employeeService *EmployeeService) GetByID(employeeID int64) (*Employee, int, error) {
	var employee Employee

	err := cache.Instance.Get(GetEmployeeCacheKey(employee.ID), &employee)
	if err != nil {
		if err != cache.ErrCacheMiss {
			log.Errorf("Error while fetching the current employee from the cache")
			return nil, http.StatusInternalServerError, err
		} else {
			query := "SELECT id, name, dob FROM employees where id = $1 and active = $2"
			err = employeeService.db.QueryRow(query, employeeID, true).
				Scan(&employee.ID, &employee.Name, &employee.DOB)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, http.StatusNotFound, errors.New(constants.ErrEmployeeRecordNotFound)
				}
				log.Errorf("Error while fetching the current employee from the DB")
				return nil, http.StatusInternalServerError, err
			}
			err = cache.Instance.Set(GetEmployeeCacheKey(employee.ID), employee, cacheTimeout)
			if err != nil {
				log.Errorf("Error while writing the current employee to the cache")
				return nil, http.StatusInternalServerError, err
			}
		}
	}
	return &employee, http.StatusOK, nil
}

// Create ...
func (employeeService *EmployeeService) Create(requestData *EmployeeCreateRequest) (*Employee, int, error) {
	dob, err := time.Parse(constants.DOB_FORMAT, requestData.DOB)
	if err != nil {
		log.Errorf("Error while parsing the employee's DOB in the required format")
		return nil, http.StatusInternalServerError, err
	}

	var lastInsertID int64
	query := "INSERT INTO employees( name, dob, active ) VALUES( $1, $2, $3 ) RETURNING id"
	err = employeeService.db.QueryRow(query, requestData.Name, dob, true).Scan(&lastInsertID)

	if err != nil {
		log.Errorf("Error while insert a new employee in the DB")
		return nil, http.StatusInternalServerError, err
	}
	employee := Employee{
		ID:   lastInsertID,
		Name: requestData.Name,
		DOB:  dob,
	}
	err = cache.Instance.Set(GetEmployeeCacheKey(employee.ID), employee, cacheTimeout)
	if err != nil {
		log.Errorf("Error while writing the newly created employee to the cache")
		return nil, http.StatusInternalServerError, err
	}

	return &employee, http.StatusCreated, nil
}

// Update ...
func (employeeService *EmployeeService) Update(requestData *EmployeeUpdateRequest) (*Employee, int, error) {
	dob, err := time.Parse(constants.DOB_FORMAT, requestData.DOB)
	if err != nil {
		log.Errorf("Error while parsing the employee's DOB in the required format")
		return nil, http.StatusInternalServerError, err
	}

	employee, status, err := employeeService.GetByID(requestData.ID)
	if err != nil {
		return nil, status, err
	}

	query := "Update employees set name = $1, dob = $2 where id = $3 and active = $4"
	result, err := employeeService.db.Exec(query, requestData.Name, dob, employee.ID, true)
	if err != nil {
		log.Errorf("Error while updating details of an employee in the DB")
		return nil, http.StatusInternalServerError, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected <= 0 {
		log.Errorf("Error while updating details of an employee in the DB")
		return nil, http.StatusInternalServerError, errors.New(constants.ErrRecordNotUpdated)
	}
	employee.Name = requestData.Name
	employee.DOB = dob

	err = cache.Instance.Set(GetEmployeeCacheKey(employee.ID), employee, cacheTimeout)
	if err != nil {
		log.Errorf("Error while writing the updated employee record to the cache")
		return nil, http.StatusInternalServerError, err
	}

	return employee, http.StatusOK, nil
}
