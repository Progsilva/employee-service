package employees

import (
	"context"
	"fmt"
	"github.com/wneessen/go-mail"
)

type Service struct {
	store Store
	email MailClient
}

func NewService(store Store, email MailClient) *Service {
	return &Service{
		store: store,
		email: email,
	}
}

func (s *Service) Login(ctx context.Context, username, password string) (*Employee, error) {
	return s.store.Login(ctx, username, password)
}

func (s *Service) CreateEmployee(ctx context.Context, employee *Employee) (*Employee, error) {
	exists, err := s.store.CheckEmployeeExists(ctx, employee.UserName, employee.Email);
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameOrEmailExists
	}
	if _, err = s.store.GetDepartment(ctx, employee.DepartmentID); err != nil {
		return nil, err
	}
	ID, err := s.store.CreateEmployee(ctx, employee)
	if err != nil {
		return nil, err
	}
	employee, err = s.store.GetEmployee(ctx, *ID)
	if err != nil {
		return nil, err
	}
	email, err := welcomeEmail(employee.Email)
	if err = s.email.DialAndSendWithContext(ctx, email); err != nil {
		return nil, err
	}
	return employee, nil
}

func welcomeEmail(email string) (*mail.Msg, error) {
	msg := mail.NewMsg()
	if err := msg.From("company@example.com"); err != nil {
		return nil, fmt.Errorf("failed to set From address: %s", err)
	}
	if err := msg.To(email); err != nil {
		return nil, fmt.Errorf("failed to set To address: %s", err)
	}
	msg.Subject("Welcome to the company!")
	msg.SetBodyString(mail.TypeTextPlain, "Welcome!")
	return msg, nil
}

func (s *Service) GetEmployee(ctx context.Context, ID int64) (*Employee, error) {
	return s.store.GetEmployee(ctx, ID)
}

func (s *Service) UpdateEmployee(ctx context.Context, employee *Employee) (*Employee, error) {
	if _, err := s.store.GetDepartment(ctx, employee.DepartmentID); err != nil {
		return nil, err
	}
	if err := s.store.UpdateEmployee(ctx, employee); err != nil {
		return nil, err
	}
	employee, err := s.store.GetEmployee(ctx, employee.ID)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func (s *Service) DeleteEmployee(ctx context.Context, ID int64) error {
	return s.store.DeleteEmployee(ctx, ID)
}

func (s *Service) Employees(ctx context.Context, departmentID *int, sort []string, limit, offset int) ([]*Employee, error) {
	employees, err := s.store.ListEmployees(ctx, departmentID, sort, limit, offset)
	if err != nil {
		return nil, err
	}
	return employees, nil
}
