package model

import "gorm.io/gorm"

//Setting -
type Setting struct {
	gorm.Model
	Category    string `json:"category" gorm:"index:awsman_setting_key,unique"`
	Key         string `json:"key" gorm:"index:awsman_setting_key,unique"`
	Value       string `json:"value"`
	Description string `json:"description"`
}
