package elb

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//RegisterRoutes -
func (s *Service) RegisterRoutes(e *echo.Echo, prefix string) {
	g := e.Group(prefix)
	//prefix = "/elbv2"
	s.prefix = prefix

	//authoriz
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    s._rsa.GetPublicKey(),
		SigningMethod: "RS512",
		TokenLookup:   "header:Authorization",
	}))

	g.GET("/listelb", s.elbv2ListLoadBalancers)
	g.GET("/listeners", s.elbv2GetListeners)
	g.GET("/rules", s.elbv2GetRules)
}
