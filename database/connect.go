package database

import (
	"collection-format/config"
	"collection-format/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectDB connect to db
func ConnectDB() {
	var err error
	var port string = "localhost"

	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Config("DB_USER"), config.Config("DB_PASSWORD"), port, config.Config("DB_PORT"), config.Config("DB_NAME"))
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")

	DB.AutoMigrate(&model.Info{}, &model.Collection{}, &model.Item{}, &model.Folder{}, &model.Example{})

	DB.AutoMigrate(&model.Response{}, &model.Request{}, &model.Header{}, &model.Body{}, &model.Url{}, &model.Query{})

	fmt.Println("Database Migrated")
}
