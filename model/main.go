// Package model
/**
  @author: zhenxiangjin
  @update_date: 2023/12/13
  @note:
*/
package model

import (
	"coolrun-backend-go/common"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//DeletedAt time.Time `json:"deleted_at"`
}

var DB *gorm.DB

func InitDB() (err error) {
	var db *gorm.DB
	dsn := "xxxx:xxx@tcp(127.0.0.1:3306)/xxxx?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err == nil {
		common.SysLog("Connect to Database successful!")
		DB = db
		// migrate Table User
		err := db.AutoMigrate(&User{})
		if err != nil {
			return err
		}
		// migrate Table Record
		err = db.AutoMigrate(&Record{})
		if err != nil {
			return err
		}
		//if err != nil {
		//	return err
		//}
		//// migrate Table UserWarehouse
		//err = db.AutoMigrate(&UserWarehouse{})
		//if err != nil {
		//	return err
		//}
		//// migrate Table item
		//err = db.AutoMigrate(&Item{})
		//if err != nil {
		//	return err
		//}
	} else {
		common.FatalLog(err)
	}
	return err
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	return err
}
