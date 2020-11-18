package cx

import (
	"fmt"
	"strings"

	"ditto.co.jp/awsman/model"
	"github.com/jinzhu/gorm"
)

//CognitoUserSearch -
type CognitoUserSearch struct {
	Email  string
	Group  string
	Offset int
	Limit  int
	Sort   string
}

//EnableCognitoUser - update user status to enabled/disabled
func (d *Database) EnableCognitoUser(uuid, status string) error {
	tx := d.DB().Begin()
	enabled := false
	if status == "enabled" {
		enabled = true
	}
	result := d.DB().Model(model.CognitoUser{}).Where("uuid=?", uuid).
		Updates(map[string]interface{}{"enabled": enabled})
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	//commit
	tx.Commit()

	return nil
}

//CreateCognitoUser - update / insert
func (d *Database) CreateCognitoUser(req *model.CognitoUser) error {
	tx := d.DB().Begin()

	//update
	result := d.DB().Model(model.CognitoUser{}).Where("uuid=?", req.UUID).
		Updates(map[string]interface{}{"name": req.Name, "enabled": req.Enabled, "group_name": req.GroupName})
	if result.Error != nil || result.RowsAffected == 0 {
		if gorm.IsRecordNotFoundError(result.Error) || result.RowsAffected == 0 {
			//insert
			result = d.DB().Create(req)
		}
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	//commit
	tx.Commit()

	return nil
}

//DeleteCognitoUser -
func (d *Database) DeleteCognitoUser(uuid string) error {
	tx := d.DB().Begin()

	if uuid != "" {
		sql := fmt.Sprintf("DELETE FROM cognito_users WHERE uuid=%q", uuid)
		result := d.DB().Exec(sql)
		// result := d.DB().Unscoped().Where("uuid=?", uuid).Delete(&model.CognitoUser{})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

	} else {
		result := d.DB().Exec("DELETE FROM cognito_users WHERE id>0")
		// result := d.DB().Unscoped().Where("id>0").Delete(&model.CognitoUser{})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	//commit
	tx.Commit()

	return nil
}

//GetCognitoUser -
func (d *Database) GetCognitoUser(uuid string) (*model.CognitoUser, error) {
	result := model.CognitoUser{
		UUID: uuid,
	}

	res := d.DB().Where("uuid=?", uuid).First(&result)
	if res.Error != nil {
		return nil, res.Error
	}

	return &result, nil
}

//GetCognitoUsers -
func (d *Database) GetCognitoUsers(input *CognitoUserSearch) ([]*model.CognitoUser, error) {
	emailike := fmt.Sprintf("%%%v%%", input.Email)
	result := make([]*model.CognitoUser, 0)
	orderby := ""
	switch strings.ToLower(input.Sort) {
	case "+email":
		orderby = "email"
	case "-email":
		orderby = "email desc"
	case "+enabled":
		orderby = "enabled"
	case "-enabled":
		orderby = "enabled desc"
	case "+group":
		orderby = "group_name"
	case "-group":
		orderby = "group_name desc"
	case "+updatedate":
		orderby = "update_date"
	case "-updatedate":
		orderby = "update_date desc"
	default:
		orderby = "id"
	}

	rows, err := d.DB().Table("cognito_users").
		Order(orderby).
		Where(`(email LIKE ? or ?) and (group_name=? or ?) and deleted_at is null`, emailike, (input.Email == ""), input.Group, (input.Group == "")).
		Offset(input.Offset).Limit(input.Limit).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	offset := input.Offset
	for rows.Next() {
		var record model.CognitoUser
		d.DB().ScanRows(rows, &record)
		record.ROWID = offset + 1
		result = append(result, &record)
		offset++
	}

	return result, nil
}

//GetCognitoUserCount -
func (d *Database) GetCognitoUserCount(input *CognitoUserSearch) (int64, error) {
	emailike := fmt.Sprintf("%%%v%%", input.Email)
	var count int64
	result := d.DB().Table("cognito_users").
		Where(`(email LIKE ? or ?) and (group_name=? or ?) and deleted_at is null`, emailike, (input.Email == ""), input.Group, (input.Group == "")).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

//DeleteCognitoGroup -
func (d *Database) DeleteCognitoGroup(userpool string) error {
	tx := d.DB().Begin()

	if userpool != "" {
		sql := fmt.Sprintf("DELETE FROM cognito_groups WHERE user_pool_id=%q", userpool)
		result := d.DB().Exec(sql)
		// result := d.DB().Where("user_pool_id=?", userpool).Delete(&model.CognitoGroup{})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

	} else {
		result := d.DB().Exec("DELETE FROM cognito_groups WHERE id>0")
		// result := d.DB().Where("id>0").Delete(&model.CognitoGroup{})
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	//commit
	tx.Commit()

	return nil
}

//CreateCognitoGroup - update / insert
func (d *Database) CreateCognitoGroup(req *model.CognitoGroup) error {
	tx := d.DB().Begin()

	//update
	result := d.DB().Model(model.CognitoGroup{}).Where("name=? and user_pool_id=?", req.Name, req.UserPoolID).Updates(req)
	if result.Error != nil || result.RowsAffected == 0 {
		if gorm.IsRecordNotFoundError(result.Error) || result.RowsAffected == 0 {
			//insert
			result = d.DB().Create(req)
		}
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}

	//commit
	tx.Commit()

	return nil
}

//GetCognitoGroup - select
func (d *Database) GetCognitoGroup(userpoolid string) ([]*model.CognitoGroup, error) {
	result := make([]*model.CognitoGroup, 0)

	rows, err := d.DB().Table("cognito_groups").Order("id").
		Where(`(user_pool_id=? or ?)`, userpoolid, (userpoolid == "")).
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record model.CognitoGroup
		d.DB().ScanRows(rows, &record)

		result = append(result, &record)
	}

	return result, nil
}
