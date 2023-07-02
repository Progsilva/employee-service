package employees

import "time"

type Employee struct {
	ID           int64
	FirstName    string
	LastName     string
	UserName     string
	Password     string
	Email        string
	Dob          time.Time
	DepartmentID int64
	Position     string
}

type Department struct {
	ID   int64
	Name string
}
