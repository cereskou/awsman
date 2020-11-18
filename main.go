//go:generate go run gen/buildnumber.go
//go:generate goversioninfo -64

package main

import (
	"fmt"
	"os"

	"ditto.co.jp/awsman/cmd"
	"ditto.co.jp/awsman/config"
	"ditto.co.jp/awsman/logger"
	"github.com/jessevdk/go-flags"
)

//options -
type options struct {
	Port int `short:"p" long:"port" description:"server port" default:"9898"`
}

// @title AWS Management
// @version 1.0
// @description AWSサーバー管理用
// @license.name ditto.co.jp
// @license.url https://www.ditto.co.jp
// @host localhost:4000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(-1)
	}

	conf, err := config.Load()
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	conf.Version = version

	logger.Infof("Aws assistant %v", version)

	if conf.Aws.AccessKey == "" || conf.Aws.SecretKey == "" {
		fmt.Println("No valid credentials avaliable.")
		os.Exit(-1)
	}

	err = cmd.RunServer()
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}
}
