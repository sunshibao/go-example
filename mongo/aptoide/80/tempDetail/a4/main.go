package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

type Ws80Detail struct {
	Info  Info  `json:"info"`
	Nodes Nodes `json:"nodes"`
}

type Info struct {
	Status string `json:"status"`
}

type Nodes struct {
	Meta Meta `json:"meta"`
}

type Meta struct {
	Data Data `json:"data"`
}

type Data struct {
	ID  int `json:"id"`
	Age Age `json:"age"`
}

type Age struct {
	Pegi string `json:"pegi"`
}

var DB *gorm.DB

func main() {

	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	mysqldb, err := gorm.Open("mysql", uri)
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
		if err2 == nil && skip < 20000 {
			skip = 0 + limit*s
			err2 = GetApkList(skip, limit)
			s++
		} else {
			break
		}
	}
	return

}

func GetApkList(skip, limit int) (err error) {
	sql1 := "select ws_id from ws80 where id>60000 limit ?,? "
	rows, err := DB.Raw(sql1, skip, limit).Rows()
	if err != nil {
		return err
	} else {
		for rows.Next() {
			var wsId int
			err := rows.Scan(&wsId)
			if err != nil {
				continue
			}
			fmt.Println("package:apkId:", wsId, "--------skip:", skip)
			GetApkDetail(wsId)
		}
	}
	return nil
}

func GetApkDetail(wsId int) (err error) {

	url := fmt.Sprintf("https://ws75.aptoide.com/api/7/app/get/store_name=catappult/app_id=%d", wsId)

	resp, err := http.Get(url)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var ws80 Ws80Detail
	json.Unmarshal([]byte(string(body)), &ws80)
	if err != nil || ws80.Info.Status != "OK" {
		sql1 := "update ws80 set pull_status = 1 where ws_id = ? "
		DB.Exec(sql1, wsId)
	} else {
		sql2 := "update ws80_detail set age_pegi = ? where ws_id = ? "
		DB.Exec(sql2, ws80.Nodes.Meta.Data.Age.Pegi, wsId)
	}

	return nil
}
