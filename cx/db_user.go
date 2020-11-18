package cx

import (
	"ditto.co.jp/awsman/model"
)

//GetUser -
func (d *Database) GetUser(name string) (*model.User, error) {
	result := model.User{
		Name: name,
	}

	res := d.DB().Where("name=?", name).First(&result)
	if res.Error != nil {
		return nil, res.Error
	}

	return &result, nil
}

//CreateUser - update / insert
func (d *Database) CreateUser(req *model.User) error {
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
