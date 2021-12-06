package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

type Ws75 struct {
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

type MysqlWs75 struct {
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
}

var DB *gorm.DB

func main() {

	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	mysqldb, err := gorm.Open("mysql", uri)
	if err != nil {
		fmt.Println("mysql连接失败")
		mysqldb.Close()
	}
	DB = mysqldb

	i := 0
	var wg = sync.WaitGroup{}
	for {
		wg.Add(1)
		go func(i int) {
			getApkList(i * 25)
			wg.Done()
		}(i)
		wg.Wait()
		if i*25 > 100000 {
			break
		}
		i++
	}
	return

	//shell(mongodb, database, skip, limit)

}

func getApkList(offset int) (err error) {

	log.Printf("------offset:%d", offset)
	url := fmt.Sprintf("https://ws75.aptoide.com/api/7/listApps/store_name=catappult/offset=%d", offset)

	resp, err := http.Get(url)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var ws75 Ws75
	json.Unmarshal([]byte(string(body)), &ws75)

	if ws75.Datalist.Next >= ws75.Datalist.Total {
		return err
	}
	insertWsData(ws75.Datalist)
	return nil
}

func insertWsData(dataList Datalist) (err error) {
	for _, s := range dataList.List {
		advertising := 0
		billing := 0
		if s.Appcoins.Advertising {
			advertising = 1
		}
		if s.Appcoins.Billing {
			billing = 1
		}
		newMysql := MysqlWs75{
			s.Id,
			s.Name,
			s.Package,
			s.Uname,
			s.Size,
			s.Icon,
			s.Graphic,
			s.Added,
			s.Modified,
			s.Updated,
			s.Uptype,
			s.Store.Id,
			s.Store.Name,
			s.Store.Avatar,
			s.File.Vername,
			s.File.Vercode,
			s.File.Md5sum,
			s.Stats.Downloads,
			s.Stats.Pdownloads,
			advertising,
			billing,
		}
		DB.Table("ws77").Create(newMysql)
	}
	return nil
}
