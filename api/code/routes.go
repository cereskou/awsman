package code

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//RegisterRoutes -
func (s *Service) RegisterRoutes(e *echo.Echo, prefix string) {
	g := e.Group(prefix)
	//prefix = "/code"
	s.prefix = prefix

	//authoriz
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    s._rsa.GetPublicKey(),
		SigningMethod: "RS512",
		TokenLookup:   "header:Authorization",
	}))

	//CodeCommit list repositories
	g.GET("/codecommit/repositories", s.codecommitListRepositories)
}
