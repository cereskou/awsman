package cx

import "ditto.co.jp/awsman/model"

//CreateCloudTrailEvent - create event
func (d *Database) CreateCloudTrailEvent(req *model.CloudTrailEvent) error {
	//
	r := model.CloudTrailEvent{}
	result := d.DB().Where("event_id=?", req.EventID).First(&r)
	if result.RowsAffected != 0 {
		return nil
	}

	tx := d.DB().Begin()
	result = tx.Create(req)
	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}
	//commit
	tx.Commit()

	return nil
}
