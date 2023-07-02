package handler

import (
	"github.com/Progsilva/employee-service/cmd/employees"
	"github.com/Progsilva/employee-service/cmd/web/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"time"
)

type LoginHandler struct {
	service           *employees.Service
	secret            string
	tokenHourLifespan int
}

func NewLoginHandler(s *employees.Service, secret string, tokenHourLifespan int) *LoginHandler {
	return &LoginHandler{
		service:           s,
		secret:            secret,
		tokenHourLifespan: tokenHourLifespan,
	}
}

func (h *LoginHandler) Endpoints(g *gin.Engine) {
	g.POST("/login", h.login)
}

func (h *LoginHandler) login(c *gin.Context) {
	loginParams := &Login{}
	if err := c.ShouldBindJSON(loginParams); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	token, err := h.loginCheck(c, loginParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "The username or password is not correct",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *LoginHandler) loginCheck(c *gin.Context, login *Login) (string, error) {
	employee, err := h.service.Login(c, login.Username, login.Password)
	if err != nil {
		return "", err
	}

	token, err := h.generateToken(employee.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (h *LoginHandler) generateToken(ID int64) (string, error) {
	claims := middleware.CustomClaims{
		UserID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        strconv.FormatInt(ID, 10),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.secret))
}
