package dbCon

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func init() {
	uri := "root:sun18188@tcp(120.26.218.217:3306)/test13?charset=utf8mb4&parseTime=True&loc=Local"
	mysqldb, err := gorm.Open("mysql", uri)
	if err != nil {
		fmt.Errorf("mysql连接失败")
		mysqldb.Close()
		return
	}
	DB = mysqldb
}
func GetDatabase() *gorm.DB {
	return DB
}
