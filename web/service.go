package web

import (
	"ditto.co.jp/awsman/api/rsa"
	"ditto.co.jp/awsman/config"
	"ditto.co.jp/awsman/cx"
	"ditto.co.jp/awsman/model"
	"ditto.co.jp/awsman/utils"
)

//Service struct
type Service struct {
	prefix string
	_db    *cx.Database
	_rsa   *rsa.RSA
}

//NewService -
func NewService(db *cx.Database, r *rsa.RSA) *Service {
	svc := &Service{
		_db:  db,
		_rsa: r,
	}

	//default root account
	conf, _ := config.Load()
	user := model.User{
		Name:     conf.Account.Name,
		Password: utils.HashPassword(conf.Account.Password),
		Avatar:   conf.Account.Avatar,
		Role:     "admin",
	}
	result := db.DB().Model(model.User{}).Where("name=?", user.Name).
		Updates(map[string]interface{}{"password": user.Password, "avatar": user.Avatar})
	if result.Error != nil || result.RowsAffected == 0 {
		result = db.DB().Create(&user)
		if result.Error != nil {
			return nil
		}
	}

	return svc
}

//Close -
func (s *Service) Close() {

}

//DB -
func (s *Service) DB() *cx.Database {
	return s._db
}
