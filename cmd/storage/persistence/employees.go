package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Progsilva/employee-service/cmd/employees"
	"github.com/golang-sql/civil"
	"strings"
)

type Employees struct {
	pool *sql.DB
}

var (
	_ employees.Store = (*Employees)(nil)
)

func NewEmployeePersistence(pool *sql.DB) *Employees {
	return &Employees{pool: pool}
}

func (e *Employees) Login(ctx context.Context, username, password string) (*employees.Employee, error) {
	tsql := `SELECT ID, FIRST_NAME, LAST_NAME, USERNAME, EMAIL, DOB, POSITION, DEPARTMENT_ID
			FROM dbo.Employee 
         	WHERE USERNAME = @USERNAME AND PWDCOMPARE(@PASSWORD, PASSWORD) = 1;`

	stmt, err := e.pool.Prepare(tsql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(
		ctx,
		sql.Named("USERNAME", username),
		sql.Named("PASSWORD", password),
	)

	employee := new(employees.Employee)
	switch err = row.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.UserName,
		&employee.Email, &employee.Dob, &employee.Position, &employee.DepartmentID); err {
	case nil:
		return employee, nil
	default:
		return nil, err
	}
}

func (e *Employees) CreateEmployee(ctx context.Context, employee *employees.Employee) (*int64, error) {
	tsql := `
      INSERT INTO dbo.Employee (FIRST_NAME, LAST_NAME, USERNAME, PASSWORD, EMAIL, DOB, POSITION, DEPARTMENT_ID) 
      VALUES (@FIRST_NAME, @LAST_NAME, @USERNAME, PWDENCRYPT(@PASSWORD),
              @EMAIL, @DOB, @POSITION, @DEPARTMENT_ID);
      select isNull(SCOPE_IDENTITY(), -1);
    `

	stmt, err := e.pool.Prepare(tsql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		ctx,
		sql.Named("FIRST_NAME", employee.FirstName),
		sql.Named("LAST_NAME", employee.LastName),
		sql.Named("USERNAME", employee.UserName),
		sql.Named("PASSWORD", employee.Password),
		sql.Named("EMAIL", employee.Email),
		sql.Named("DOB", civil.DateOf(employee.Dob)),
		sql.Named("POSITION", employee.Position),
		sql.Named("DEPARTMENT_ID", employee.DepartmentID),
	)
	var newID int64
	if err = row.Scan(&newID); err != nil {
		return nil, err
	}
	return &newID, err
}

func (e *Employees) GetEmployee(ctx context.Context, id int64) (*employees.Employee, error) {
	tsql := `SELECT ID, FIRST_NAME, LAST_NAME, USERNAME, EMAIL, DOB, POSITION, DEPARTMENT_ID
			FROM dbo.Employee 
         	WHERE ID = @ID;`

	stmt, err := e.pool.Prepare(tsql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(
		ctx,
		sql.Named("ID", id),
	)

	employee := new(employees.Employee)
	switch err = row.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.UserName,
		&employee.Email, &employee.Dob, &employee.Position, &employee.DepartmentID); err {
	case sql.ErrNoRows:
		return nil, employees.ErrEmployeeNotFound
	case nil:
		return employee, nil
	default:
		return nil, err
	}
}

func (e *Employees) CheckEmployeeExists(ctx context.Context, username, email string) (bool, error) {
	tsql := `SELECT ID
			FROM dbo.Employee 
         	WHERE USERNAME = @USERNAME OR EMAIL = @EMAIL;`

	stmt, err := e.pool.Prepare(tsql)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(
		ctx,
		sql.Named("USERNAME", username),
		sql.Named("EMAIL", email),
	)
	if err != nil {
		return false, err
	}
	for rows.Next() {
		return true, nil
	}
	return false, nil
}

func (e *Employees) UpdateEmployee(ctx context.Context, employee *employees.Employee) error {
	tsql := `
      UPDATE dbo.Employee SET 
      FIRST_NAME = @FIRST_NAME, 
      LAST_NAME = @LAST_NAME, 
      USERNAME = @USERNAME, 
      EMAIL = @EMAIL, 
      DOB = @DOB, 
      POSITION = @POSITION, 
      DEPARTMENT_ID = @DEPARTMENT_ID
      WHERE ID = @ID;
    `

	stmt, err := e.pool.Prepare(tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(
		ctx,
		sql.Named("FIRST_NAME", employee.FirstName),
		sql.Named("LAST_NAME", employee.LastName),
		sql.Named("USERNAME", employee.UserName),
		sql.Named("EMAIL", employee.Email),
		sql.Named("DOB", civil.DateOf(employee.Dob)),
		sql.Named("POSITION", employee.Position),
		sql.Named("DEPARTMENT_ID", employee.DepartmentID),
		sql.Named("ID", employee.ID),
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return employees.ErrEmployeeNotFound
	}
	return nil
}

func (e *Employees) DeleteEmployee(ctx context.Context, ID int64) error {
	tsql := `DELETE FROM dbo.Employee 
         	WHERE ID = @ID;`

	stmt, err := e.pool.Prepare(tsql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.ExecContext(ctx, sql.Named("ID", ID))
	if err != nil {
		return err
	}
	return nil
}

func (e *Employees) GetDepartment(ctx context.Context, ID int64) (*employees.Department, error) {
	tsql := `SELECT ID, NAME
			FROM dbo.Department 
         	WHERE ID = @ID;`

	stmt, err := e.pool.Prepare(tsql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(
		ctx,
		sql.Named("ID", ID),
	)

	department := new(employees.Department)
	switch err = row.Scan(&department.ID, &department.Name); err {
	case sql.ErrNoRows:
		return nil, employees.ErrDepartmentNotFound
	case nil:
		return department, nil
	default:
		return nil, err
	}
}

func (e *Employees) ListEmployees(ctx context.Context,
	departmentID *int, sort []string, limit int, offset int) ([]*employees.Employee, error) {
	tsql := `
		SELECT ID, FIRST_NAME, LAST_NAME, USERNAME, EMAIL, DOB, POSITION, DEPARTMENT_ID 
		FROM dbo.Employee
		WHERE (@DEPARTMENT_ID is null or DEPARTMENT_ID = @DEPARTMENT_ID)
		#ORDER_BY#
		OFFSET @OFFSET ROWS
		FETCH NEXT @LIMIT ROWS ONLY
	;`
	tsql = strings.ReplaceAll(tsql, "#ORDER_BY#", orderBy(sort))

	stmt, err := e.pool.Prepare(tsql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(
		ctx,
		sql.Named("DEPARTMENT_ID", departmentID),
		sql.Named("OFFSET", offset),
		sql.Named("LIMIT", limit),
	)
	if err != nil {
		return nil, err
	}
	list := make([]*employees.Employee, 0)
	for rows.Next() {
		employee := new(employees.Employee)
		if err = rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.UserName,
			&employee.Email, &employee.Dob, &employee.Position, &employee.DepartmentID); err != nil {
			return nil, err
		}
		list = append(list, employee)
	}
	return list, nil
}

func orderBy(sort []string) string {
	if len(sort) == 0 {
		return "ORDER BY FIRST_NAME ASC"
	}
	for i, s := range sort {
		sort[i] = ascOrDesc(s)
	}
	return fmt.Sprintf("ORDER BY %s", strings.Join(sort, ","))
}

func ascOrDesc(s string) string {
	if strings.HasPrefix(s, "-") {
		return fmt.Sprintf(" %s DESC", s[1:])
	}
	return fmt.Sprintf(" %s ASC", s)
}
