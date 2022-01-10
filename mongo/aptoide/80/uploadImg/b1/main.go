package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sunshibao/go-utils/util/gconv"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var CosClient *cos.Client
var DB *gorm.DB

func main() {
	//建立连接
	NewCosClient()
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
	sql1 := "select id,ws_id,package,icon,media_screenshots from ws80_detail limit ?,?"
	rows, err := DB.Raw(sql1, skip, limit).Rows()
	if err != nil {
		return err
	} else {

		for rows.Next() {
			var id int
			var wsId int
			var packageName string
			var icon string
			var image string
			err := rows.Scan(&id, &wsId, &packageName, &icon, &image)
			if err != nil {
				continue
			}
			fmt.Println("ws_id:", wsId, "--------skip:", skip)
			if icon != "" {
				baseName := "aptoide_icon"
				split := strings.Split(icon, "catappult/")
				iconName := gconv.String(wsId) + "-" + packageName + "-" + split[1]
				newName := baseName + "/" + iconName
				go UploadCos(id, newName, icon)
			}
			if image != "" {
				baseName := "aptoide_img"
				split2 := strings.Split(image, ",")
				wg := sync.WaitGroup{}
				for _, v := range split2 {
					wg.Add(1)
					split := strings.Split(v, "catappult/")
					imgName := gconv.String(wsId) + "-" + packageName + "-" + split[1]
					newName := baseName + "/" + imgName
					go func() {
						defer wg.Done()
						UploadCos(id, newName, v)
					}()

				}
				wg.Wait()
			}
		}
	}
	return nil
}

func NewCosClient() {
	var secretid string = "AKIDjHZaKn0xc0GJ4ZnlRr0tVqtgCSR9alfK"
	var secretkey string = "QaT5RVIo56qJVQ5TaQzHI2WjeKktmOkO"
	var cosUrl string = "https://gp-image-1308128293.cos.eu-moscow.myqcloud.com"

	u, _ := url.Parse(cosUrl)
	b := &cos.BaseURL{BucketURL: u}
	CosClient = cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretid,
			SecretKey: secretkey,
		},
	})
}

func UploadCos(id int, name, filePath string) {
	c := CosClient
	resp, err := http.Get(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	_, err = c.Object.Put(context.Background(), name, resp.Body, nil)
	if err != nil {
		sql := `update ws80_detail set img_pull_status = -1 where id = ?`
		DB.Exec(sql, id)
		fmt.Println(err)
	}
}
