package model

import (
	"github.com/Progsilva/employee-service/cmd/employees"
	"strings"
	"time"
)

func (e *NewEmployee) ToEmployee(dob time.Time) *employees.Employee {
	return &employees.Employee{
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		UserName:     e.UserName,
		Password:     e.Password,
		Email:        e.Email,
		Dob:          dob,
		DepartmentID: int64(e.DepartmentID),
		Position:     e.Position,
	}
}

func FromEmployee(e *employees.Employee) *Employee {
	return &Employee{
		ID:           int(e.ID),
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		UserName:     e.UserName,
		Password:     e.Password,
		Email:        e.Email,
		Dob:          e.Dob.Format("02-01-2006"),
		DepartmentID: int(e.DepartmentID),
		Position:     e.Position,
	}
}

func FromEmployees(list []*employees.Employee) []*Employee {
	employeeList := make([]*Employee, len(list))
	for i, employee := range list {
		employeeList[i] = FromEmployee(employee)
	}
	return employeeList
}

func (e *UpdateEmployee) ToEmployee(dob time.Time) *employees.Employee {
	return &employees.Employee{
		ID:           int64(e.ID),
		FirstName:    e.FirstName,
		LastName:     e.LastName,
		UserName:     e.UserName,
		Password:     e.Password,
		Email:        e.Email,
		Dob:          dob,
		DepartmentID: int64(e.DepartmentID),
		Position:     e.Position,
	}
}

func Limit(limit *int) int {
	if limit == nil {
		return 10
	}
	l := *limit
	if l < 10 {
		return 10
	}
	if l > 50 {
		return 50
	}
	return l
}

func Offset(offset *int) int {
	if offset == nil {
		return 0
	}
	o := *offset
	if o < 0 {
		return 0
	}
	return o
}

var (
	sortMap = map[string]bool{"FIRST_NAME": true, "EMAIL": true}
)

func FilterSort(sort []*string) []string {
	if len(sort) == 0 {
		return nil
	}
	ss := make([]string, 0)
	for _, s := range sort {
		v := strings.TrimPrefix(strings.TrimSpace(*s), "-")
		v = strings.ToUpper(v)
		if sortMap[v] {
			ss = append(ss, *s)
		}
	}
	return ss
}
