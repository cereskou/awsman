package cognito

import (
	"ditto.co.jp/awsman/api/rsa"
	"ditto.co.jp/awsman/awss"
	"ditto.co.jp/awsman/config"
	"ditto.co.jp/awsman/cx"
)

//Service struct
type Service struct {
	prefix string
	_aws   *awss.Service
	_db    *cx.Database
	_rsa   *rsa.RSA
}

//NewService -
func NewService(db *cx.Database, r *rsa.RSA) *Service {
	svc := &Service{
		_aws: awss.New(),
		_db:  db,
		_rsa: r,
	}

	conf, _ := config.Load()
	//cognito setting
	settings, err := db.GetCognitoSettings()
	if err == nil {
		for _, kv := range settings {
			if kv.Key == "userpool" {
				conf.Cognito.UserPoolID = kv.Value
			}
		}
	}

	return svc
}

//Close -
func (s *Service) Close() {

}

//Aws -get aws ervice
func (s *Service) Aws() *awss.Service {
	return s._aws
}

//DB -
func (s *Service) DB() *cx.Database {
	return s._db
}
