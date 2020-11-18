package web

import (
	"github.com/labstack/echo/v4"
)

//RegisterRoutes -
func (s *Service) RegisterRoutes(e *echo.Echo, prefix string) {
	g := e.Group(prefix)

	s.prefix = prefix

	//user
	g.GET("/info", s.UsersInfo)
	g.POST("/login", s.Login)
	g.POST("/logout", s.Logout)
	g.POST("", s.CreateUser)
}
