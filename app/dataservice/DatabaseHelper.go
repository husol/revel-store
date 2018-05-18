package dataservice

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"husol.org/mypham/app/models"
	"log"
)

var dbcon *gorm.DB

func GetDBConnection() *gorm.DB {
	return dbcon
}

func InitDb(driver, spec string) {
	var err error
	// Open DB
	dbcon, err = gorm.Open(driver, spec)
	checkErr(err, "sql.Open failed.")
	dbcon.AutoMigrate(&models.User{})
	dbcon.AutoMigrate(&models.Category{})
	dbcon.AutoMigrate(&models.Product{})
	dbcon.AutoMigrate(&models.Comment{})
	dbcon.AutoMigrate(&models.Transaction{})
	dbcon.AutoMigrate(&models.Order{})
	dbcon.AutoMigrate(&models.Information{})
	checkErr(err, "Create tables failed.")
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
		fmt.Println(msg)
	}
}
