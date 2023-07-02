package app

import (
	"crypto/tls"
	"github.com/Progsilva/employee-service/cmd/employees"
	"github.com/Progsilva/employee-service/cmd/storage/dbclient"
	"github.com/Progsilva/employee-service/cmd/storage/persistence"
	"github.com/Progsilva/employee-service/cmd/web"
	"github.com/wneessen/go-mail"
	"os"
)

type App struct {
	srv *web.Server
}

func New() (*App, error) {
	// init clients
	dbPool, err := dbclient.Pool()
	if err != nil {
		return nil, err
	}
	store := persistence.New(dbPool)

	mailhogPort := os.Getenv("MAILHOG_HOST")
	email, err := mail.NewClient(
		mailhogPort,
		mail.WithPort(1025),
		mail.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		mail.WithTLSPolicy(mail.NoTLS),
	)

	if err != nil {
		return nil, err
	}
	// init service/business
	service := employees.NewService(store, email)

	// init web layer
	srv, err := web.New(service)
	if err != nil {
		return nil, err
	}

	return &App{
		srv: srv,
	}, nil
}

func (a *App) Start() error {
	return a.srv.Serve()
}
