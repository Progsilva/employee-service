package employees

import "errors"

var (
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrDepartmentNotFound    = errors.New("department not found")
	ErrUsernameOrEmailExists = errors.New("username or email already exists")
)
