package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/Progsilva/employee-service/cmd/employees"
	"github.com/Progsilva/employee-service/graph"
	"github.com/gin-gonic/gin"
)

type GraphqlHandler struct {
	service *employees.Service
}

func NewGraphqlHandler(s *employees.Service) *GraphqlHandler {
	return &GraphqlHandler{service: s}
}

func (h *GraphqlHandler) Endpoints(g *gin.RouterGroup) {
	g.POST("/employee", h.graphqlHandler())
}

func (h *GraphqlHandler) graphqlHandler() gin.HandlerFunc {
	defaultHandler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Service: h.service}}))
	return func(c *gin.Context) {
		defaultHandler.ServeHTTP(c.Writer, c.Request)
	}
}
