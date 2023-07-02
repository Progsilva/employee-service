package employees

import (
	"context"

	go_mail "github.com/wneessen/go-mail"
)

//go:generate mockgen -source=interfaces.go -destination=interfaces_mock.go -package=employees

type Store interface {
	Login(ctx context.Context, username, password string) (*Employee, error)

	CheckEmployeeExists(ctx context.Context, username, email string) (bool, error)
	CreateEmployee(ctx context.Context, e *Employee) (*int64, error)
	GetEmployee(ctx context.Context, id int64) (*Employee, error)
	UpdateEmployee(ctx context.Context, employee *Employee) error
	DeleteEmployee(ctx context.Context, ID int64) error
	ListEmployees(ctx context.Context, departmentID *int, sort []string, limit int, offset int) ([]*Employee, error)

	GetDepartment(ctx context.Context, id int64) (*Department, error)
}

type MailClient interface {
	DialAndSendWithContext(ctx context.Context, ml ...*go_mail.Msg) error
}
