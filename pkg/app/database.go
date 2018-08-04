package app

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
)

type DB struct {
	Conn *gorm.DB
}

type DBConfig struct {
	Drive string
	Args  string
}

func NewDBConfig(user, password string) *DBConfig {
	return &DBConfig{
		Drive: "mysql",
		Args: fmt.Sprintf("%s:%s@tcp(mysqldb)/todo_list?charset=utf8&parseTime=True&loc=Local",
			user, password),
	}
}

func InitDB(config *DBConfig) (database *DB, err error) {
	db := &DB{}
	db.Conn, err = gorm.Open(config.Drive, config.Args)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Migrate(models ...interface{}) []error {
	errs := db.Conn.AutoMigrate(models...).GetErrors()

	return errs
}

func (db *DB) Find(model interface{}) {
	db.Conn.Find(model)
}
