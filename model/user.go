package model

import "gorm.io/gorm"

//User -
type User struct {
	gorm.Model
	Name     string   `json:"name" gorm:"index:awsman_user_name,unique"`
	Password string   `json:"password"`
	Avatar   string   `json:"avatar,omitempty"`
	Role     string   `json:"-" gorm:"roles"`
	Roles    []string `json:"roles" gorm:"-"`
}
