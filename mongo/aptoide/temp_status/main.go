package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

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
	sql1 := "select ws_id,name,package,cate_type from ws78 limit ?,? "
	rows, err := DB.Raw(sql1, skip, limit).Rows()
	if err != nil {
		return err
	} else {
		for rows.Next() {
			var ws78Temp Ws78Temp
			err := DB.ScanRows(rows, &ws78Temp)
			if err != nil {
				continue
			}
			fmt.Println("package:", ws78Temp.Package, "--------skip:", skip)
			GetApkDetail(ws78Temp)
		}
	}
	return nil
}

func GetApkDetail(temp Ws78Temp) (err error) {

	url := fmt.Sprintf("https://ws75.aptoide.com/api/7/app/get/store_name=catappult/app_id=%d", temp.WsID)

	resp, err := http.Get(url)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var ws75 Ws75Detail
	json.Unmarshal([]byte(string(body)), &ws75)
	if err == nil {
		insertWsData(ws75.Nodes, temp)
	} else {
		sql1 := "update ws78 set pull_status = 1 where ws_id = ? "
		DB.Exec(sql1, temp.WsID)
	}

	return nil
}

func insertWsData(nodes Nodes, ws78Temp Ws78Temp) (err error) {
	downUrl := nodes.Meta.Data.File.Path
	newMysql := MysqlWs75Detail{
		ws78Temp.WsID,
		ws78Temp.Name,
		ws78Temp.Package,
		ws78Temp.CateType,
		downUrl,
	}
	DB.Table("ws78_detail").Create(newMysql)
	go DownApk(downUrl, ws78Temp.WsID)
	return nil
}

func DownApk(apkUrl string, wsId int) {

	apkPath := "/Users/sunshibao/Desktop/apkDown/"

	fileName := path.Base(apkUrl)

	res, err := http.Get(apkUrl)
	if err != nil {
		sql1 := "update ws78_detail set down_status = 1 where ws_id = ? "
		DB.Exec(sql1, wsId)
		fmt.Println("A error occurred!---", apkUrl)
		return
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	file, err := os.Create(apkPath + fileName)
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d", written)
}
