package model

import "gorm.io/gorm"

//CognitoUser -
type CognitoUser struct {
	gorm.Model
	UUID       string `json:"uuid,omitempty" gorm:"index:cognitor_user_name,unique"`
	ROWID      int    `json:"id,omitempty" gorm:"-"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Sub        string `json:"sub,omitempty"`
	Enabled    bool   `json:"enabled"`
	UserStatus string `json:"status,omitempty"`
	GroupName  string `json:"group,omitempty"`
	UserPoolID string `json:"userpoolid,omitempty"`
	CreateDate int64  `json:"createdate,omitempty"`
	UpdateDate int64  `json:"updatedate,omitempty"`
}

//CognitoGroup -
type CognitoGroup struct {
	gorm.Model
	UserPoolID  string `json:"userpoolid" gorm:"index:cognitor_group_idx,unique"`
	Name        string `json:"name" gorm:"index:cognitor_group_idx,unique"`
	Description string `json:"description"`
}
