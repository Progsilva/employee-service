package employees

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestService_CreateEmployee(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := NewMockStore(ctrl)
	email := NewMockMailClient(ctrl)
	service := NewService(store, email)

	ctx := context.Background()

	t.Run("fail when employee exists", func(t *testing.T) {
		// given
		connor := newEmployee()

		store.EXPECT().CheckEmployeeExists(ctx, connor.UserName, connor.Email).Return(true, nil)

		// when
		_, err := service.CreateEmployee(ctx, connor)

		// then
		require.ErrorIs(t, err, ErrUsernameOrEmailExists)
	})

	t.Run("fail when department does not exists", func(t *testing.T) {
		// given
		connor := newEmployee()

		store.EXPECT().CheckEmployeeExists(ctx, connor.UserName, connor.Email).Return(false, nil)
		store.EXPECT().GetDepartment(ctx, connor.DepartmentID).Return(nil, ErrDepartmentNotFound)

		// when
		_, err := service.CreateEmployee(ctx, connor)

		// then
		require.ErrorIs(t, err, ErrDepartmentNotFound)
	})

	t.Run("fail when db problem occurs", func(t *testing.T) {
		// given
		connor := newEmployee()
		dbErr := errors.New("some db problem")

		store.EXPECT().CheckEmployeeExists(ctx, connor.UserName, connor.Email).Return(false, nil)
		store.EXPECT().GetDepartment(ctx, connor.DepartmentID).Return(nil, nil)
		store.EXPECT().CreateEmployee(ctx, connor).Return(nil, dbErr)

		// when
		_, err := service.CreateEmployee(ctx, connor)

		// then
		require.ErrorIs(t, err, dbErr)
	})

	t.Run("fail when email host is down", func(t *testing.T) {
		// given
		connor := newEmployee()
		getConnor := newEmployee()
		getConnor.ID = 10
		ID := int64(10)
		emailErr := errors.New("some email problem")

		store.EXPECT().CheckEmployeeExists(ctx, connor.UserName, connor.Email).Return(false, nil)
		store.EXPECT().GetDepartment(ctx, connor.DepartmentID).Return(nil, nil)
		store.EXPECT().CreateEmployee(ctx, connor).Return(&ID, nil)
		store.EXPECT().GetEmployee(ctx, ID).Return(getConnor, nil)
		email.EXPECT().DialAndSendWithContext(ctx, gomock.Any()).Return(emailErr)

		// when
		_, err := service.CreateEmployee(ctx, connor)

		// then
		require.ErrorIs(t, err, emailErr)
	})

	t.Run("create employee success", func(t *testing.T) {
		// given
		connor := newEmployee()
		getConnor := newEmployee()
		getConnor.ID = 10
		ID := int64(10)

		store.EXPECT().CheckEmployeeExists(ctx, connor.UserName, connor.Email).Return(false, nil)
		store.EXPECT().GetDepartment(ctx, connor.DepartmentID).Return(nil, nil)
		store.EXPECT().CreateEmployee(ctx, connor).Return(&ID, nil)
		store.EXPECT().GetEmployee(ctx, ID).Return(getConnor, nil)
		email.EXPECT().DialAndSendWithContext(ctx, gomock.Any()).Return(nil)

		// when
		connorCreated, err := service.CreateEmployee(ctx, connor)

		// then
		require.NoError(t, err)
		require.NotNil(t, connorCreated)
		require.Equal(t, getConnor.ID, connorCreated.ID)
		require.Equal(t, getConnor.Email, connorCreated.Email)
		require.Equal(t, getConnor.Dob, connorCreated.Dob)
		require.Equal(t, getConnor.UserName, connorCreated.UserName)
		require.Equal(t, getConnor.DepartmentID, connorCreated.DepartmentID)
		require.Equal(t, getConnor.Password, connorCreated.Password)
		require.Equal(t, getConnor.Position, connorCreated.Position)
		require.Equal(t, getConnor.LastName, connorCreated.LastName)
		require.Equal(t, getConnor.FirstName, connorCreated.FirstName)
	})
}

func newEmployee() *Employee {
	return &Employee{
		FirstName:    "Connor",
		LastName:     "MacLeod",
		UserName:     "highlander",
		Password:     "whoWantsToLiveForEver",
		Email:        "high_lander@gmail.com",
		Dob:          time.Date(1986, time.August, 29, 0, 0, 0, 0, time.UTC),
		DepartmentID: 2,
		Position:     "Head of Tech",
	}
}
