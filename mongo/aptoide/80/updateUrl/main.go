package main

import (
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sunshibao/go-utils/util/gconv"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var CosClient *cos.Client
var DB *gorm.DB

func main() {
	//for i := 0; i <= 9; i++ {
	//	minId := i * 10000
	//	go func(id int) {
	//		start(id)
	//	}(minId)
	//}
	start(0)
}

func start(minId int) {
	//建立连接
	uri := "root:Droi*#2021@tcp(127.0.0.1:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"

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
		if err2 == nil && skip < 100808 {
			skip = 0 + limit*s
			err2 = GetApkList(minId, skip, limit)
			s++
		} else {
			break
		}
	}
	return
}

func GetApkList(id, skip, limit int) (err error) {
	sql1 := "select id,ws_id,package,icon,media_screenshots,file_path from ws80_detail where id >? order by id limit ?,? "
	rows, err := DB.Raw(sql1, id, skip, limit).Rows()
	if err != nil {
		return err
	} else {
		for rows.Next() {
			var id int
			var wsId int
			var packageName string
			var icon string
			var image string
			var apkPath string
			err := rows.Scan(&id, &wsId, &packageName, &icon, &image, &apkPath)
			if err != nil {
				continue
			}
			fmt.Println(packageName, "--------", skip)
			if icon != "" {
				baseName := "https://gp-image-1308128293.cos.eu-moscow.myqcloud.com/aptoide_icon/"
				iconName := gconv.String(wsId) + "-" + packageName + "-"
				index := strings.LastIndex(icon, "/")
				fileName := icon[index+1:]
				newUrl := baseName + iconName + fileName
				UploadCos(id, newUrl)
			}
			if image != "" {
				baseName := "https://gp-image-1308128293.cos.eu-moscow.myqcloud.com/aptoide_img/"
				split2 := strings.Split(image, ",")
				imgName := gconv.String(wsId) + "-" + packageName + "-"
				tempUrl := []string{}
				for _, v := range split2 {
					index := strings.LastIndex(v, "/")
					fileName := v[index+1:]
					newUrl := baseName + imgName + fileName
					tempUrl = append(tempUrl, newUrl)
				}
				UploadCos2(id, strings.Join(tempUrl, ","))
			}
			if apkPath != "" {
				baseName := "https://gp-image-1308128293.cos.eu-moscow.myqcloud.com/aptoide_apk/"
				index := strings.LastIndex(apkPath, "/")
				fileName := apkPath[index+1:]
				newUrl := baseName + fileName
				UploadCos3(id, newUrl)
			}
		}
	}
	return nil
}

func UploadCos(id int, newUrl string) {
	sql := `update ws80_detail set icon = ? where id = ?`
	DB.Exec(sql, newUrl, id)
}

func UploadCos2(id int, newUrl string) {
	sql := `update ws80_detail set media_screenshots = ? where id = ?`
	DB.Exec(sql, newUrl, id)
}

func UploadCos3(id int, newUrl string) {
	sql := `update ws80_detail set file_path = ? where id = ?`
	DB.Exec(sql, newUrl, id)
}
