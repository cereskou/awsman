package cx

import (
	"fmt"
	"path/filepath"

	"ditto.co.jp/awsman/logger"
	"ditto.co.jp/awsman/model"
	"ditto.co.jp/awsman/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//Database - sqlite3 database
type Database struct {
	Error error
	Name  string
	_db   *gorm.DB
}

//NewDB - cereate a new database
func NewDB(dir, name string, autovacuum bool) *Database {
	if name == "" {
		name = "awsman.db"
	}
	dbase := &Database{
		Error: nil,
	}
	parent := filepath.Join(dir, "awsman")
	err := utils.MakeDirectory(parent)
	if err != nil {
		parent = dir
	}
	name = filepath.Join(parent, name)
	logger.Infof("cache db %v", name)

	dnsfmt := "%v"
	if autovacuum {
		dnsfmt = "file:%v?_auto_vacuum=1"
	}

	dns := fmt.Sprintf(dnsfmt, name)
	db, err := gorm.Open(sqlite.Open(dns), &gorm.Config{})
	if err != nil {
		dbase.Error = err
		dbase.Name = name

		return dbase
	}
	//set
	dbase._db = db
	dbase.Name = name
	dbase.Error = nil

	return dbase
}

//Migration -
func (d *Database) Migration() error {
	return d._db.AutoMigrate(
		&model.CognitoUser{},
		&model.CognitoGroup{},
		&model.User{},
		&model.Setting{},
		&model.CloudTrailEventResource{},
		&model.CloudTrailEvent{},
	)
}

//Close -
func (d *Database) Close() error {
	return nil
}

//DB -
func (d *Database) DB() *gorm.DB {
	return d._db
}

//Begin -
func (d *Database) Begin() *gorm.DB {
	return d._db.Begin()
}

//Rollback -
func (d *Database) Rollback() *gorm.DB {
	return d._db.Rollback()
}

//Commit -
func (d *Database) Commit() *gorm.DB {
	return d._db.Commit()
}
