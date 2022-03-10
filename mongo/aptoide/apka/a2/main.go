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
	"sync"

	"github.com/sunshibao/go-utils/util/gconv"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/tencentyun/cos-go-sdk-v5"
)

var CosClient *cos.Client
var DB *sqlx.DB

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			minId := i * 20000
			start(minId)
		}(i)
	}
	wg.Wait()
	//start(0)
}

func start(minId int) {
	//建立连接
	NewCosClient()
	uri := "root:Droi*#2021@tcp(18.197.156.118:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:@tcp(127.0.0.1:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"

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
		if err2 == nil && skip < 20000 {
			skip = limit * s
			err2 = GetApkList(minId, skip, limit)
			s++
		} else {
			break
		}
	}
	return

	//GetApkList(0, 1)
}

type ApkInfo struct {
	ApkId       int64  `db:"apk_id" json:"apk_id"`
	PackageName string `db:"package_name" json:"package_name"`
	FileSize    int64  `db:"file_size" json:"file_size"`
	DownloadUrl string `db:"download_url" json:"download_url"`
}

func GetApkList(minId, skip, limit int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()

	apkInfo := ApkInfo{}
	sql3 := "select apk_id,package_name,file_size,download_url from ap_temp_table where id>? and temp_status = 0 limit ?,? "
	err = DB.Get(&apkInfo, sql3, minId, skip, limit)

	if apkInfo.DownloadUrl == "" {
		return nil
	}

	log.Println("hd_apk_url" + gconv.String(skip))

	resultIcon := CheckImgUrl(apkInfo.ApkId, apkInfo.PackageName, apkInfo.DownloadUrl, apkInfo.FileSize)
	if resultIcon {
		NewUploadImg(apkInfo.DownloadUrl)
	}
	return nil
}

func NewUploadImg(hdImgUrl string) {
	updateNum++
	index := strings.Split(hdImgUrl, "aptoide_apk/")
	s := index[1]
	//fmt.Println(index[0])
	basePath := "https://syncPool.apk.aptoide.com/catappult/"
	newUrl := basePath + s
	newName := "aptoide_apk/" + s
	UploadCos(newName, newUrl)
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

func CheckImgUrl(apkId int64, packageName, filePath string, fileSize int64) bool {
	res, err := http.Get(filePath)
	if err != nil {
		log.Println("CheckImgUrl-====================err:", err)
		return false
	}
	defer res.Body.Close()
	//fmt.Println(packageName+"size============:", res.ContentLength)
	if res.ContentLength < fileSize {
		log.Println(err, "未识别的apk filePath:", filePath)
		delSql := "insert into  oz_apk_temp2 (package_name,download_url) values(?,?)"
		DB.Exec(delSql, packageName, filePath)
		return true
	} else {
		delSql := "update ap_temp_table set temp_status = 1 where apk_id = ?"
		DB.Exec(delSql, apkId)
	}

	return false
}

var updateNum int

func UploadCos(name, filePath string) {
	c := CosClient
	//filePath = strings.Replace(filePath, "https", "http", 1)
	resp, err := http.Get(filePath)
	if err != nil {
		log.Println(err, "获取图片失败:"+filePath)
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReaderSize(resp.Body, 0)

	imgPath := "/data/apkNew/"
	fileName := path.Base(filePath)
	file, err := os.Create(imgPath + fileName)
	if err != nil {
		fmt.Println(err, "url:"+filePath)
		return
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d \n", written)

	//本地上传
	_, err = c.Object.PutFromFile(context.Background(), name, imgPath+fileName, nil)
	if err != nil {
		fmt.Println("333333333", err)
	} else {
		fmt.Println("444444444", "上传成功 url:"+imgPath+fileName)
		err := os.Remove(imgPath + fileName)
		if err != nil {
			fmt.Println("删除失败:", imgPath+fileName)
		} else {
			fmt.Println("删除成功:", imgPath+fileName)
		}
	}
}
