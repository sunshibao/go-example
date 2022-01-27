package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 9; i++ {
		wg.Add(1)
		minId := i * 10000
		go func(id int) {
			defer wg.Done()
			start(id)
		}(minId)
	}
	wg.Wait()
	//start(0)
}

func start(id int) {
	uri := "root:Droi*#2021@tcp(18.197.156.118:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

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
		if err2 == nil && skip < 10000 {
			skip = 0 + limit*s
			err2 = GetApkList(id, skip, limit)
			s++
		} else {
			break
		}
	}
	return

}

func GetApkList(minId, skip, limit int) (err error) {
	var wsId int
	var packageName string
	var icon string
	var imgList string
	sql1 := "select ws_id,package,icon,media_screenshots from ws80_detail_copy where id>? limit ?,? "
	err = DB.QueryRow(sql1, minId, skip, limit).Scan(&wsId, &packageName, &icon, &imgList)
	if err != nil {
		return err
	} else {
		fmt.Println("package:", "--------skip:", skip)
		GetApkDetail(wsId, packageName, icon, imgList)
	}
	return nil
}

func GetApkDetail(wsId int, packageName, icon, imgList string) (err error) {

	imgSql := `insert into ws75_image (ws_id,package_name,image_name,image_type,image_width,image_height,hd_image_url,nhd_image_url,language) values (?,?,?,?,?,?,?,?,?)`
	exec, _ := DB.Exec(imgSql, wsId, packageName, packageName+"_Icon", 5, 0, 0, icon, icon, "ru")
	lastImgId, _ := exec.LastInsertId()
	fmt.Println("insertWsData", lastImgId)
	apkSql := `insert into ws75_apk_image (apk_id,language,image_id) values (?,?,?) ;`
	DB.Exec(apkSql, wsId, "ru", lastImgId)

	split2 := strings.Split(imgList, ",")

	for _, v := range split2 {
		imgSql := `insert into ws75_image (ws_id,package_name,image_name,image_type,image_width,image_height,hd_image_url,nhd_image_url,language) values (?,?,?,?,?,?,?,?,?)`
		exec, _ := DB.Exec(imgSql, wsId, packageName, packageName+"_Screenshots", 5, 0, 0, v, v, "ru")
		lastImgId, _ := exec.LastInsertId()
		fmt.Println("insertWsData", lastImgId)
		apkSql := `insert into ws75_apk_image (apk_id,language,image_id) values (?,?,?) ;`
		DB.Exec(apkSql, wsId, "ru", lastImgId)
	}

	return nil

}
