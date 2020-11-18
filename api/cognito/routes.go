package cognito

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//RegisterRoutes -
func (s *Service) RegisterRoutes(e *echo.Echo, prefix string) {
	g := e.Group(prefix)
	//prefix = "/cognito"
	s.prefix = prefix

	//authoriz
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:    s._rsa.GetPublicKey(),
		SigningMethod: "RS512",
		TokenLookup:   "header:Authorization",
	}))

	//cognito
	g.GET("/user/list", s.cognitoListUser)
	g.PUT("/user/:uuid", s.cognitoUpdateUser)
	g.POST("/user", s.cognitoCreateUser)
	g.DELETE("/user/:uuid", s.cognitoDeleteUser)
	g.GET("/user/sync", s.cognitoSyncUser)
	g.POST("/user/enable/:uuid", s.cognitoEnableUser)
	g.POST("/user/setpassword/:uuid", s.cognitoSetPassword)

	//setting
	g.GET("/setting", s.cognitoSettings)
	g.POST("/setting", s.cognitoSaveSettings)
}
