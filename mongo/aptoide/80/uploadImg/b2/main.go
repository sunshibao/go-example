package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sunshibao/go-utils/util/gconv"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var CosClient *cos.Client
var DB *gorm.DB

//474333
func main() {
	//wg := sync.WaitGroup{}
	//for i := 0; i <= 12; i++ {
	//	wg.Add(1)
	//	minId := i * 50000
	//	go func(id int) {
	//		defer wg.Done()
	//		start(id)
	//	}(minId)
	//}
	//wg.Wait()
	start(0)
}

func start(minId int) {
	//建立连接
	NewCosClient()
	uri := "root:Droi*#2021@tcp(18.197.156.118:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"

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
		if err2 == nil && skip < 50000 {
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
	sql1 := "select image_id,ws_id,package_name,hd_image_url,image_name from ws75_image where image_id>? and img_up_status = 0 and hd_image_url!='' limit ?,?"
	rows, err := DB.Raw(sql1, id, skip, limit).Rows()
	if err != nil {
		return err
	} else {
		for rows.Next() {
			var imageId int
			var wsId int
			var packageName string
			var hdImageUrl string
			var imageName string

			err := rows.Scan(&imageId, &wsId, &packageName, &hdImageUrl, &imageName)
			if err != nil {
				continue
			}
			fmt.Println("ws_id:", wsId, "--------skip:", skip)
			lastInd := strings.LastIndex(imageName, "_")

			if hdImageUrl == "" {
				continue
			}
			if imageName[lastInd+1:] == "Icon" {
				baseName := "aptoide_icon"
				split := strings.Split(hdImageUrl, "catappult/")
				iconName := gconv.String(wsId) + "-" + packageName + "-" + split[1]
				newName := baseName + "/" + iconName
				UploadCos(imageId, newName, hdImageUrl)
			} else {
				baseName := "aptoide_img"
				split := strings.Split(hdImageUrl, "catappult/")
				imgName := gconv.String(wsId) + "-" + packageName + "-" + split[1]
				newName := baseName + "/" + imgName
				UploadCos(imageId, newName, hdImageUrl)
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

func UploadCos(imageId int, name, filePath string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("捕获到的错误：%s\n", r)
		}
	}()
	c := CosClient
	filePath = strings.Replace(filePath, "https", "http", 1)
	resp, err := http.Get(filePath)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReaderSize(resp.Body, 32*1024)

	imgPath := "/data/image/"
	fileName := path.Base(filePath)
	file, err := os.Create(imgPath + fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	log.Printf("Total length: %d", written)

	//本地上传
	_, err = c.Object.PutFromFile(context.Background(), name, imgPath+fileName, nil)
	if err != nil {
		log.Println(err)
		return
	} else {
		sql := `update ws75_image set img_up_status = 1 where image_id = ?`
		DB.Exec(sql, imageId)
		err := os.Remove(imgPath + fileName)
		if err != nil {
			fmt.Println("删除失败:", imgPath+fileName)
		} else {
			fmt.Println("删除成功:", imgPath+fileName)
		}
	}
	return
}
