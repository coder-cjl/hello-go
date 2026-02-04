package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db2 *gorm.DB

func init() {
	database, err := gorm.Open(mysql.Open("root:123456@~!@tcp(localhost:3306)/dev?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		Log.Error("open mysql failed")
		return
	}
	Db2 = database
}

type MyGormSQL struct{}

func go_gorm_autoMigrate() {
	err := Db2.AutoMigrate(&Person{})
	if err != nil {
		Log.Error("AutoMigrate failed:", err)
		return
	}
	Log.Info("AutoMigrate succeeded")
}

func go_gorm_ts1() {
	err := Db2.Transaction(func(tx *gorm.DB) error {
		p := Person{
			Username: "Anna1",
			Sex:      "F",
			Email:    "anan@qq.com",
		}
		if err := tx.Table(p.TableName()).Create(&p).Error; err != nil {
			return err
		}
		Log.Info("inserted person id:", p.UserId)
		return nil
	})
	if err != nil {
		Log.Error("Transaction failed:", err)
		return
	}
	Log.Info("Transaction committed successfully")
}

// 查询数据
func go_gorm_ts2() {
	var persons []Person
	result := Db2.Table(Person{}.TableName()).Where(&Person{UserId: 1001}).Find(&persons)
	if result.Error != nil {
		Log.Error("Query failed:", result.Error)
		return
	}
	for _, p := range persons {
		Log.Info("Person:", p)
	}
}

func (m MyGormSQL) Test() {
	// go_gorm_autoMigrate()
	// go_gorm_ts1()
	go_gorm_ts2()
}
