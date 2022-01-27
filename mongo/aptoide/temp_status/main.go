package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

type Ws75Detail struct {
	Nodes Nodes `json:"nodes"`
}

type Nodes struct {
	Meta Meta `json:"meta"`
}

type Meta struct {
	Data Data `json:"data"`
}

type Data struct {
	File File `json:"file"`
}

type File struct {
	Path string `json:"path"`
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
		if err2 == nil && skip < 1600 {
			skip = 0 + limit*s
			err2 = GetApkList(skip, limit)
			s++
		} else {
			break
		}
	}
	return

}

type Ws78Temp struct {
	WsID     int    `json:"ws_id"`
	Name     string `json:"name"`
	Package  string `json:"package"`
	CateType string `json:"cate_type"`
}

type MysqlWs75Detail struct {
	WsID     int    `json:"ws_id"`
	Name     string `json:"name"`
	Package  string `json:"package"`
	CateType string `json:"cate_type"`
	DownUrl  string `json:"down_url"`
}

func GetApkList(skip, limit int) (err error) {
	sql1 := "select ws_id,down_url from ws78_detail limit ?,? "
	rows, err := DB.Raw(sql1, skip, limit).Rows()
	if err != nil {
		return err
	} else {
		for rows.Next() {
			wsId := 0
			downUrl := ""
			err := rows.Scan(&wsId, &downUrl)
			if err != nil {
				continue
			}
			DownApk(downUrl, wsId)
		}
	}
	return nil
}

func DownApk(apkUrl string, wsId int) {
	split := strings.Split(apkUrl, "catappult/")
	fmt.Println(split[0])
	fmt.Println(split[1])
	path := fmt.Sprintf("/Users/sunshibao/Desktop/apkDown2/%s", split[1])
	b, err := PathExists(path)
	if err != nil {
		fmt.Printf("PathExists(%s),err(%v)\n", path, err)
	}
	if b {
		fmt.Printf("path %s 存在\n", path)
	} else {
		sql1 := "update ws78_detail set down_status = -1 where ws_id = ? "
		DB.Exec(sql1, wsId)
	}
}

/*
   判断文件或文件夹是否存在
   如果返回的错误为nil,说明文件或文件夹存在
   如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
   如果返回的错误为其它类型,则不确定是否在存在
*/
func PathExists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
