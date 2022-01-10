package main

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i <= 20; i++ {
		wg.Add(1)
		minId := i * 2000000
		go func(id int) {
			defer wg.Done()
			start(id)
		}(minId)
	}
	wg.Wait()
}

func start(minId int) {
	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:sun18188@tcp(127.0.0.1:3306)/test13?charset=utf8mb4&parseTime=True&loc=Local"

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

type ImgInfo struct {
	ImageId    int    `db:"image_id"`
	HdImageUrl string `db:"hd_image_url"`
}

func shell(minId, skip, limit int) (err error) {
	var imgInfos []ImgInfo
	sql := `select image_id,hd_image_url from oz_image where image_id>? and temp_status = 0 order by image_id limit ?,?;`
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

func execShell(imgInfos []ImgInfo) {
	for _, v := range imgInfos {
		fileNameWithSuffix := path.Base(v.HdImageUrl)
		fileType := path.Ext(fileNameWithSuffix)
		if fileType == "" {
			v.HdImageUrl = v.HdImageUrl + ".png"
		}
		v.HdImageUrl = strings.Replace(v.HdImageUrl, "http://18.177.149.123:8001/pic/", "http://apk-ry.tt286.com/app_img/", 1)
		v.HdImageUrl = strings.Replace(v.HdImageUrl, "https://gp-image-1308128293.cos.eu-moscow.myqcloud.com/", "http://apk-ry.tt286.com/", 1)
		apkSql := "update oz_image set hd_image_url = ?,nhd_image_url = ?,temp_status = 1 where image_id = ?"
		DB.Exec(apkSql, v.HdImageUrl, v.HdImageUrl, v.ImageId)
		fmt.Println("ImageId:", v.ImageId)
	}
}
