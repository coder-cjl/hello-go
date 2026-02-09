package mysql

import (
	"hello-go/src/logger"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	database, err := gorm.Open(mysql.Open("root:123456@~!@tcp(localhost:3306)/dev?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		logger.Error("MySQL 连接失败", zap.Error(err))
		return
	}
	Db = database
}

func Insert(table string, data interface{}) error {
	result := Db.Table(table).Create(data)
	return result.Error
}

func Query(table string, conditions interface{}, dest interface{}) error {
	result := Db.Table(table).Where(conditions).Find(dest)
	return result.Error
}

func Update(table string, conditions interface{}, updates interface{}) error {
	result := Db.Table(table).Where(conditions).Updates(updates)
	return result.Error
}

func Delete(table string, conditions interface{}) error {
	result := Db.Table(table).Where(conditions).Delete(nil)
	return result.Error
}
