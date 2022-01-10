package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

type Ws80 struct {
	Datalist Datalist `json:"datalist"`
}

type Datalist struct {
	Total  int    `json:"total"`
	Count  int    `json:"count"`
	offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Next   int    `json:"next"`
	List   []List `json:"list"`
}
type List struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Package  string   `json:"package"`
	Uname    string   `json:"uname"`
	Size     int      `json:"size"`
	Icon     string   `json:"icon"`
	Graphic  string   `json:"graphic"`
	Added    string   `json:"added"`
	Modified string   `json:"modified"`
	Updated  string   `json:"updated"`
	Uptype   string   `json:"uptype"`
	Store    Store    `json:"store"`
	File     File     `json:"file"`
	Stats    Stats    `json:"stats"`
	Appcoins Appcoins `json:"appcoins"`
}

type Store struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type File struct {
	Vername string `json:"vername"`
	Vercode int    `json:"vercode"`
	Md5sum  string `json:"md5sum"`
}

type Stats struct {
	Downloads  int `json:"downloads"`
	Pdownloads int `json:"pdownloads"`
}

type Appcoins struct {
	Advertising bool `json:"advertising"`
	Billing     bool `json:"billing"`
}

type MysqlWs80 struct {
	WsId            int    `gorm:"column:ws_id" db:"ws_id" json:"ws_id" form:"ws_id"`
	Name            string `gorm:"column:name" db:"name" json:"name" form:"name"`
	Package         string `gorm:"column:package" db:"package" json:"package" form:"package"`
	Uname           string `gorm:"column:uname" db:"uname" json:"uname" form:"uname"`
	Size            int    `gorm:"column:size" db:"size" json:"size" form:"size"`
	Icon            string `gorm:"column:icon" db:"icon" json:"icon" form:"icon"`
	Graphic         string `gorm:"column:graphic" db:"graphic" json:"graphic" form:"graphic"`
	Added           string `gorm:"column:added" db:"added" json:"added" form:"added"`
	Modified        string `gorm:"column:modified" db:"modified" json:"modified" form:"modified"`
	Updated         string `gorm:"column:updated" db:"updated" json:"updated" form:"updated"`
	Uptype          string `gorm:"column:uptype" db:"uptype" json:"uptype" form:"uptype"`
	StoreId         int    `gorm:"column:store_id" db:"store_id" json:"store_id" form:"store_id"`
	StoreName       string `gorm:"column:store_name" db:"store_name" json:"store_name" form:"store_name"`
	StoreAvatar     string `gorm:"column:store_avatar" db:"store_avatar" json:"store_avatar" form:"store_avatar"`
	FileVername     string `gorm:"column:file_vername" db:"file_vername" json:"file_vername" form:"file_vername"`
	FileVercode     int    `gorm:"column:file_vercode" db:"file_vercode" json:"file_vercode" form:"file_vercode"`
	FileMd5sum      string `gorm:"column:file_md5sum" db:"file_md5sum" json:"file_md5sum" form:"file_md5sum"`
	StatsDownloads  int    `gorm:"column:stats_downloads" db:"stats_downloads" json:"stats_downloads" form:"stats_downloads"`
	StatsPdownloads int    `gorm:"column:stats_pdownloads" db:"stats_pdownloads" json:"stats_pdownloads" form:"stats_pdownloads"`
	Advertising     int    `gorm:"column:advertising" db:"advertising" json:"advertising" form:"advertising"`
	Billing         int    `gorm:"column:billing" db:"billing" json:"billing" form:"billing"`
	AptoideType     int    `gorm:"column:aptoide_type" db:"aptoide_type" json:"aptoide_type" form:"aptoide_type"`
	CateType        string `gorm:"column:cate_type" db:"cate_type" json:"cate_type" form:"cate_type"`
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

	GetApkCate()

}

func GetApkCate() (err error) {

	GetApkList("games")
	return nil
}

var ApkListTotal int

func GetApkList(category string) (err error) {
	url := fmt.Sprintf("https://ws75.aptoide.com/api/7/listApps/store_name=catappult/group_name=%s/offset=0", category)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var ws80 Ws80
	json.Unmarshal([]byte(string(body)), &ws80)
	ApkListTotal = ws80.Datalist.Total

	skip := 0
	limit := 25
	offset := 0
	for {
		offset = skip * limit
		skip++
		if offset < ApkListTotal {
			GetListAll(category, offset)
		} else {
			break
		}
	}

	return nil
}

func GetListAll(category string, offset int) {
	url := fmt.Sprintf("https://ws75.aptoide.com/api/7/listApps/store_name=catappult/group_name=%s/offset=%d", category, offset)
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var ws80 Ws80
	json.Unmarshal([]byte(string(body)), &ws80)
	ApkListTotal = ws80.Datalist.Total
	updateWsData(ws80.Datalist.List)

}

func updateWsData(list []List) {
	for _, v := range list {
		fmt.Println("ws_id:", v.Id)
		sql := "update ws80 set apk_type =1 where ws_id = ?"
		DB.Exec(sql, v.Id)
	}
}
