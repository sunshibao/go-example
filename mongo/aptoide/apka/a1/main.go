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
	for i := 0; i < 9; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			minId := 438676 + i*10000
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
		if err2 == nil && skip < 10000 {
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
	PackageName string `db:"package_name" json:"package_name"`
	DownloadUrl string `db:"download_url" json:"download_url"`
}

func GetApkList(minId, skip, limit int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()

	apkInfo := ApkInfo{}
	sql3 := "select package_name,download_url from oz_apk where apk_id > ? and company_type = 2 order by apk_id limit ?,? "
	err = DB.Get(&apkInfo, sql3, minId, skip, limit)

	if apkInfo.DownloadUrl == "" {
		return nil
	}

	log.Println("hd_apk_url" + gconv.String(skip))

	resultIcon := CheckImgUrl(apkInfo.PackageName, apkInfo.DownloadUrl)
	if !resultIcon {
		NewUploadImg(apkInfo.DownloadUrl)
	}
	return nil
}

func NewUploadImg(hdImgUrl string) {
	updateNum++
	index := strings.Split(hdImgUrl, "aptoide_apk/")
	s := index[1]
	//fmt.Println(index[0])
	basePath := "https://pool.apk.aptoide.com/catappult/"
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

func CheckImgUrl(packageName, filePath string) bool {
	res, err := http.Get(filePath)
	if err != nil {
		return false
	}
	//fmt.Println(packageName+"size============:", res.ContentLength)
	if res.ContentLength < 500000 {
		log.Println(err, "未识别的图片 filePath:", filePath)
		delSql := "insert into  oz_apk_temp2 (package_name,download_url) values(?,?)"
		DB.Exec(delSql, packageName, filePath)
		return false
	}

	return true
}

var updateNum int

func UploadCos(name, filePath string) {
	c := CosClient
	filePath = strings.Replace(filePath, "https", "http", 1)
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
	fmt.Printf("Total length: %d", written)

	//fmt.Println("2222222", name)
	//fmt.Println("4444444", imgPath+fileName)
	//本地上传
	_, err = c.Object.PutFromFile(context.Background(), name, imgPath+fileName, nil)
	if err != nil {
		fmt.Println("333333333", err)
	} else {
		err := os.Remove(imgPath + fileName)
		if err != nil {
			fmt.Println("删除失败:", imgPath+fileName)
		} else {
			fmt.Println("删除成功:", imgPath+fileName)
		}
	}
}
