package dao

import (
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitMysqlDB() (err error) {
	if err := createDatabase(); err != nil {
		return err
	}

	dsn := "root:supersecret@tcp(mysql:3306)/mydb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	dbMigration()

	return nil
}

func dbMigration() {
	db.AutoMigrate(&Book{})

	for i := 0; i < 100; i++ {
		id := strconv.Itoa(i)

		db.Create(&Book{
			Id:        id,
			Name:      "book#" + id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}
}

func createDatabase() error {
	dsn := "root:supersecret@tcp(mysql:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	rootDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	err = rootDb.Exec("CREATE DATABASE if not exists mydb").Error
	if err != nil {
		return err
	}

	return nil
}
