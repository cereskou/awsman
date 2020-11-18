package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"ditto.co.jp/awsman/config"
	"ditto.co.jp/awsman/logger"
	"ditto.co.jp/awsman/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "ditto.co.jp/awsman/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//RunServer -
func RunServer() error {
	logger.Debug("RunServer")

	if err := services.InitService(); err != nil {
		return err
	}
	defer services.Close()

	conf, _ := config.Load()

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Index: "index.html",
	}))
	e.Use(createLogMiddleware())

	//static
	e.Static("/", "host")

	//swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	//Router
	services.WebService.RegisterRoutes(e, "/user")
	services.CognitoService.RegisterRoutes(e, "/cognito")
	services.CloudTrailService.RegisterRoutes(e, "/cloudtrail")
	services.CodeService.RegisterRoutes(e, "/code")
	services.ELBv2Service.RegisterRoutes(e, "/elbv2")

	//Run the server
	go func() {
		address := fmt.Sprintf("%v:%v", conf.Host, conf.Port)
		e.Logger.Fatal(e.Start(address))
	}()

	//Quit
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Debug("Close")
	if err := e.Close(); err != nil {
		//e.Logger.Fatal(err)
		logger.Error(err)
		//return err
	}

	logger.Debug("Shutdown")
	if err := e.Shutdown(ctx); err != nil {
		// e.Logger.Fatal(err)
		logger.Error(err)
		// return err
	}

	return nil
}

func createLogMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqid := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 16)
			request := c.Request()
			request.Header.Add("reqid", reqid)
			//before
			logger.Tracef("[%v] %v %v START", reqid, request.Method, request.RequestURI)
			//action
			err := next(c)
			//after
			logger.Tracef("[%v] %v %v END", reqid, request.Method, request.RequestURI)
			return err
		}
	}
}
