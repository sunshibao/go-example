package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func main() {
	uri := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	mysqldb, err := sqlx.Open("mysql", uri)
	if err != nil {
		fmt.Println("mysql连接失败")
		mysqldb.Close()
	}
	DB = mysqldb

	skip := 0
	limit := 1
	s := 0
	var err2 error
	for {
		if err2 == nil {
			skip = 0 + limit*s
			err2 = shell(skip, limit)
			s++
		} else {
			break
		}
	}
	return

	//shell(mongodb, database, skip, limit)

}

type GpInfo struct {
	Id        int    `json:"id"`
	GpDownNum string `json:"gp_down_num"`
	GpRanking int    `json:"gp_ranking"`
}

//10万
//ObjectId("61755372275289742a211b8e")
func shell(skip, limit int) (err error) {
	var packageName []string

	apkSql := "select package from ws80_copy limit ?,?"
	DB.Select(&packageName, apkSql, skip, limit)

	for _, v := range packageName {
		gpInfo := GpInfo{}
		sql2 := "select id,gp_down_num,gp_ranking from ws80_copy where package_name = ?"
		DB.Get(&gpInfo, sql2, v)
		if gpInfo.Id == 0 {
			continue
		}
		//写入oz_apk表
		log.Printf("PackageName:%s-----num:%d\n", v, skip)
		apkSql3 := "update ws80 set gp_down_num = ? ,gp_ranking = ? where package = ?"
		DB.Exec(apkSql3, gpInfo.GpDownNum, gpInfo.GpRanking, v)
	}
	return nil
}
