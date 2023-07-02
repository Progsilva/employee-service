package web

import (
	"context"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Progsilva/employee-service/cmd/web/handler"
	"github.com/Progsilva/employee-service/cmd/web/middleware"
	"net/http"
	"time"

	"github.com/Progsilva/employee-service/cmd/employees"
	"github.com/Progsilva/employee-service/config"
	"github.com/gin-gonic/gin"
	"github.com/ory/graceful"
)

type Server struct {
	srv *http.Server
}

func New(service *employees.Service) (*Server, error) {
	c := &Config{}
	if err := config.Load(c); err != nil {
		return nil, err
	}
	r := gin.Default()

	r.GET("/", playgroundHandler())

	handler.NewLoginHandler(service, c.Secret, c.TokenHourLifespan).Endpoints(r)

	secureGroup := r.Group("", middleware.JwtAuthMiddleware(c.Secret))
	handler.NewGraphqlHandler(service).Endpoints(secureGroup)

	srv := &http.Server{
		Addr:    ":" + c.Port,
		Handler: r,
	}
	return &Server{srv: srv}, nil
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/employee")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) Serve() error {
	return graceful.Graceful(s.listenAndServe, s.shutdown)
}

func (s *Server) listenAndServe() error {
	return s.srv.ListenAndServe()
}

func (s *Server) shutdown(_ context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.srv.Shutdown(ctx)
}
