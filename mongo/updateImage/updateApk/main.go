package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func main() {
	startId := 521480
	wg := sync.WaitGroup{}
	for i := 0; i <= 8; i++ {
		wg.Add(1)
		minId := i * 10000
		go func(id int) {
			defer wg.Done()
			start(startId + id)
		}(minId)
	}
	wg.Wait()
}

func start(minId int) {
	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:Droi*#2021@tcp(127.0.0.1:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"

	mysqldb, err := sqlx.Connect("mysql", uri)
	if err != nil {
		fmt.Println("mysql连接失败")
		mysqldb.Close()
	}
	DB = mysqldb

	skip := 0
	limit := 100
	s := 0
	var err2 error
	for {
		if err2 == nil {
			skip = limit * s
			err2 = shell(minId, skip, limit)
			s++
		} else {
			break
		}
	}
	return

}

type ApkInfo struct {
	ApkId       int    `db:"apk_id"`
	DownloadUrl string `db:"download_url"`
}

func shell(minId, skip, limit int) (err error) {
	var imgInfos []ApkInfo
	sql := `select apk_id,download_url from oz_apk where apk_id > ? and company_type = 2 and temp_status = 0 order by apk_id limit ?,?;`
	err = DB.Select(&imgInfos, sql, minId, skip, limit)
	if err != nil {
		fmt.Println(err)
	}
	if len(imgInfos) == 0 {
		return errors.New("数据结束")
	}
	execShell(imgInfos)
	return nil
}

func execShell(imgInfos []ApkInfo) {
	for _, v := range imgInfos {
		v.DownloadUrl = strings.Replace(v.DownloadUrl, "https://gp-image-1308128293.cos.eu-moscow.myqcloud.com/", "http://apk-ry.tt286.com/", 1)
		apkSql := "update oz_apk set download_url = ?,temp_status = 1 where apk_id = ?"
		DB.Exec(apkSql, v.DownloadUrl, v.ApkId)
		fmt.Println("ApkId:", v.ApkId)
	}
}
