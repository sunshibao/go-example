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

	"github.com/sunshibao/go-utils/util/gconv"

	"github.com/disintegration/imageorient"

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
	//wg := sync.WaitGroup{}
	//for i := 0; i < 9; i++ {
	//	wg.Add(1)
	//	go func(i int) {
	//		defer wg.Done()
	//		minId := i * 500
	//		start(minId)
	//	}(i)
	//}
	//wg.Wait()

	start(0)
}

func start(minId int) {
	//建立连接
	NewCosClient()
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
		if err2 == nil {
			skip = limit * s
			err2 = GetApkList(skip, limit)
			s++
		} else {
			break
		}
	}
	return

	//GetApkList(0, 1)
}

type ApkInfo struct {
	ApkId       int    `db:"apk_id"`
	AptoideId   int    `db:"aptoide_id"`
	PackageName string `db:"package_name"`
	JpgIconId   int    `db:"jpg_icon_id"`
}

func GetApkList(skip, limit int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()

	var jpgIconUrl string
	sql3 := "select imgUrl from temp_img_url2 order by id limit ?,? "
	err = DB.Get(&jpgIconUrl, sql3, skip, limit)
	if err != nil {
		log.Println(err)
		return err
	}

	if jpgIconUrl == "" {
		return nil
	}

	log.Println("hd_image_url" + gconv.String(skip))

	resultIcon := CheckImgUrl(jpgIconUrl)
	if !resultIcon {
		NewUploadImg(jpgIconUrl)
	}
	return nil
}

func NewUploadImg(hdImgUrl string) {
	updateNum++
	index := strings.LastIndex(hdImgUrl, "-")
	s := hdImgUrl[index+1:]
	basePath := "http://syncPool.img.aptoide.com/catappult/"
	newUrl := basePath + s
	newName := strings.Split(hdImgUrl, ".com/")
	UploadCos(newName[1], newUrl)
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

func CheckImgUrl(filePath string) bool {
	res, err := http.Get(filePath)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	_, _, err = imageorient.Decode(res.Body)
	if err != nil {
		log.Println(err, "未识别的图片 filePath:", filePath)
		delSql := "insert into  temp_img_url3 (imgUrl) values(?)"
		DB.Exec(delSql, filePath)
		return false
	}
	return true
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

	reader := bufio.NewReaderSize(resp.Body, 32*1024)

	imgPath := "/data/image/"
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
		}
	}
}
