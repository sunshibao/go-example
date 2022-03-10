package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/sunshibao/go-utils/util/gconv"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"

	"strings"
)

type TempImg struct {
	Id          int    `gorm:"column:id" db:"id" json:"id" form:"id"`
	ImageId     int    `gorm:"column:image_id" db:"image_id" json:"image_id" form:"image_id"`
	ImageName   string `gorm:"column:image_name" db:"image_name" json:"image_name" form:"image_name"`
	HdImageUrl  string `gorm:"column:hd_image_url" db:"hd_image_url" json:"hd_image_url" form:"hd_image_url"`
	ApkId       int    `gorm:"column:apk_id" db:"apk_id" json:"apk_id" form:"apk_id"`
	PackageName string `gorm:"column:package_name" db:"package_name" json:"package_name" form:"package_name"`
	NewImgUrl   string `gorm:"column:new_img_url" db:"new_img_url" json:"new_img_url" form:"new_img_url"`
}

var DB *sqlx.DB
var CosClient *cos.Client
var ApkHttpClient *http.Client

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

func NewHttpClient() {
	ApkHttpClient = &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   5 * time.Second,
	}
}

func main() {
	NewCosClient()
	NewHttpClient()
	start(0)
}

func start(id int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()
	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:Droi*#2021@tcp(18.197.156.118:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

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
		if err2 == nil {
			skip = 0 + limit*s
			err2 = shell(id, skip, limit)
			s++
		} else {
			break
		}
	}
	return

	//shell(mongodb, database, skip, limit)

}

func shell(id, skip, limit int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()
	apImgSql := "select image_id,hd_image_url,apk_id,package_name from 3temp_img where id>? limit ?,?"
	rows, err := DB.Queryx(apImgSql, id, skip, limit)
	if err != nil {
		fmt.Println("sql 报错 sql:", apImgSql, "--Id:", id, "--skip:", skip, "--limit:", limit)
		return err
	} else {
		for rows.Next() {
			var tempImg TempImg
			err := rows.StructScan(&tempImg)
			if err != nil {
				fmt.Println("sql 报错 sql222:", err)
				continue
			}
			uploadImg(tempImg)
		}
	}
	return nil
}

var realNum int
var priority = 0
var vipChannel = 1    // vip通道
var companyType = 2   // gp 1,ap 2
var apkSourceSign = 5 // 1:17万俄语，2:10万英语，3：10万aptoide,4:10万英语第二次

func uploadImg(tempImg TempImg) string {
	fileType := "img"
	newImageUrl := tempImg.HdImageUrl
	var filenameWithSuffix string
	filenameWithSuffix = path.Base(newImageUrl)
	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix)
	fmt.Println("fileSuffix =", fileSuffix)
	a1 := strings.TrimSuffix(newImageUrl, fileSuffix)
	split := strings.Split(a1, "_")

	if split[1] == "icon" {
		fileType = "icon"
	}

	if newImageUrl != "" {
		baseName := "aptoide_img"
		if fileType == "icon" {
			baseName = "aptoide_icon"
		}
		split := strings.Split(newImageUrl, "honor-api/")
		imgName := gconv.String(tempImg.ApkId) + "-" + tempImg.PackageName + "-" + split[1]
		newName := baseName + "/" + imgName
		newImageUrl = UploadCos(tempImg.ApkId, newName, newImageUrl)

	}
	if newImageUrl != "" {
		//写入oz_image表
		scrSql := "update 3temp_img set new_img_url = ? where image_id = ?;"
		_, err := DB.Exec(scrSql, newImageUrl, tempImg.ImageId)
		fmt.Println("detail ：", tempImg.PackageName, err)
	}
	fmt.Println("2222222", newImageUrl)
	return newImageUrl
}

func UploadCos(id int, name, filePath string) (newFilePath string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()
	c := CosClient

	newFilePath = filePath
	filePath = strings.Replace(filePath, "https", "http", 1)
	resp, err := ApkHttpClient.Get(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReaderSize(resp.Body, 32*1024)

	imgPath := "/Users/sunshibao/Desktop/tempData/"
	fileName := path.Base(filePath)
	file, err := os.Create(imgPath + fileName)
	if err != nil {
		fmt.Println(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d", written)

	if written == int64(0) {
		return ""
	}

	basePath := "http://apk-ry.tt286.com/"
	//本地上传
	_, err = c.Object.PutFromFile(context.Background(), name, imgPath+fileName, nil)
	if err != nil {
		fmt.Println(err)
	} else {
		err := os.Remove(imgPath + fileName)
		if err != nil {
			fmt.Println("删除失败:", imgPath+fileName)
		} else {
			fmt.Println("删除成功:", imgPath+fileName)
			newFilePath = basePath + name
		}
	}
	return newFilePath
}
