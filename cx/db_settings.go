package cx

import (
	"ditto.co.jp/awsman/model"
)

//GetCognitoSettings -
func (d *Database) GetCognitoSettings() ([]*model.Setting, error) {
	result := make([]*model.Setting, 0)

	rows, err := d.DB().Table("settings").
		Where(`category=? and deleted_at is null`, "cognito").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var record model.Setting
		d.DB().ScanRows(rows, &record)

		result = append(result, &record)
	}

	return result, nil
}

//GetCognitoSetting -
func (d *Database) GetCognitoSetting(key string) (*model.Setting, error) {
	result := model.Setting{
		Key: key,
	}

	res := d.DB().Where("category=? and key=?", "cognito", key).First(&result)
	if res.Error != nil {
		return nil, res.Error
	}

	return &result, nil
}

//CreateUpdateSetting - update
func (d *Database) CreateUpdateSetting(req *model.Setting) error {
	tx := d.DB().Begin()

	//update
	result := d.DB().Model(model.Setting{}).Where("category=? and key=?", req.Category, req.Key).
		Updates(map[string]interface{}{"value": req.Value})
	if result.Error != nil || result.RowsAffected == 0 {
		result = d.DB().Create(req)
		tx.Rollback()
		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}
	}
	//commit
	tx.Commit()

	return nil
}

//CreateSetting - insert
func (d *Database) CreateSetting(req *model.Setting) error {
	tx := d.DB().Begin()

	//update
	result := d.DB().Create(req)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	//commit
	tx.Commit()

	return nil
}
