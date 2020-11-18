package cognito

import (
	"github.com/labstack/echo/v4"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	RegisterRoutes(e *echo.Echo, prefix string)
	Close()
}
