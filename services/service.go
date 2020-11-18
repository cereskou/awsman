package services

import (
	"fmt"
	"reflect"

	"ditto.co.jp/awsman/api/cloudtrail"
	"ditto.co.jp/awsman/api/code"
	"ditto.co.jp/awsman/api/cognito"
	"ditto.co.jp/awsman/api/elb"
	"ditto.co.jp/awsman/api/rsa"
	"ditto.co.jp/awsman/config"
	"ditto.co.jp/awsman/cx"
	"ditto.co.jp/awsman/disk"
	"ditto.co.jp/awsman/logger"
	"ditto.co.jp/awsman/web"
)

var (
	//database
	_db  *cx.Database
	_rsa *rsa.RSA

	//WebService ...
	WebService web.ServiceInterface
	//CognitoService ...
	CognitoService cognito.ServiceInterface
	//CloudTrailService ...
	CloudTrailService cloudtrail.ServiceInterface
	//CodeService ... CodeCommit/CodeBuild/CodePipeline
	CodeService code.ServiceInterface
	//ELBv2Service ... Load Banlance
	ELBv2Service elb.ServiceInterface
)

// InitService -
func InitService() error {
	logger.Trace("Service InitService")
	conf, _ := config.Load()
	//database
	dbname := conf.Db.Name
	if dbname == "" {
		dbname = fmt.Sprintf("awsman-%v.db", conf.Port)
		conf.Db.AutoVacuum = true
	}

	dir := disk.HomeDir()
	_db = cx.NewDB(dir, dbname, conf.Db.AutoVacuum)
	err := _db.Migration()
	if err != nil {
		return err
	}

	//Rsa Key
	_rsa, err = rsa.NewRSA(conf.Rsa.Private, conf.Rsa.Public)
	if err != nil {
		return err
	}

	//web service
	if nil == reflect.TypeOf(WebService) {
		WebService = web.NewService(_db, _rsa)
	}
	//aws coginito
	if nil == reflect.TypeOf(CognitoService) {
		CognitoService = cognito.NewService(_db, _rsa)
	}
	//cloudtrail
	if nil == reflect.TypeOf(CloudTrailService) {
		CloudTrailService = cloudtrail.NewService(_db, _rsa)
	}
	//codecommid
	if nil == reflect.TypeOf(CodeService) {
		CodeService = code.NewService(_db, _rsa)
	}
	//elbv2
	if nil == reflect.TypeOf(ELBv2Service) {
		ELBv2Service = elb.NewService(_db, _rsa)
	}

	return nil
}

//Close -
func Close() {
	logger.Trace("Service Close")
	WebService.Close()
	CognitoService.Close()
	CloudTrailService.Close()
	CodeService.Close()
	ELBv2Service.Close()

	if _db != nil {
		_db.Close()
	}
}
