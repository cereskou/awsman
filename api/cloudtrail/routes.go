package cloudtrail

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//RegisterRoutes -
func (s *Service) RegisterRoutes(e *echo.Echo, prefix string) {
	g := e.Group(prefix)
	//prefix = "/cloudtrail"
	s.prefix = prefix

	//authoriz
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    s._rsa.GetPublicKey(),
		SigningMethod: "RS512",
		TokenLookup:   "header:Authorization",
	}))

	g.GET("/events", s.cloudtrailListEvent)
}
