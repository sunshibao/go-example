package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var DB *gorm.DB

func GetDatabase() *gorm.DB {
	return DB
}

func main() {

	uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"

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
		if err2 == nil && skip < 10000 {
			skip = 0 + limit*s
			err2 = shell(skip, limit)
			s++
		} else {
			break
		}
	}
	return

	//shell(mongodb, database, skip, limit)

}

type ImageInfo struct {
	ApkId      int    `json:"apk_id"`
	Language   string `json:"language"`
	ImageId    int    `json:"image_id"`
	HdImageUrl string `json:"hd_image_url"`
}

//ObjectId("6177e34b275289742a6cf720")
func shell(skip, limit int) (err error) {
	imgPath := "/Users/sunshibao/Desktop/apkImage/"
	imageUrl := ""
	sql := "select hd_image_url from oz_image where temp_status3 = 3 and image_id>3691655 limit ?,?"
	DB.Raw(sql, skip, limit).Row().Scan(&imageUrl)

	if imageUrl == "" {
		return errors.New("无数据")
	}
	fmt.Printf("imageUrl:%s-----num:skip:%d \n", imageUrl, skip)

	fileName := path.Base(imageUrl)

	res, err := http.Get(imageUrl)
	if err != nil {
		fmt.Println("A error occurred!")
		return
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	file, err := os.Create(imgPath + fileName)
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d", written)

	return nil
}

// 获取前面结果字节的二进制
func bytesToHexString(src []byte) string {
	res := bytes.Buffer{}
	if src == nil || len(src) <= 0 {
		return ""
	}
	temp := make([]byte, 0)
	for _, v := range src {
		sub := v & 0xFF
		hv := hex.EncodeToString(append(temp, sub))
		if len(hv) < 2 {
			res.WriteString(strconv.FormatInt(int64(0), 10))
		}
		res.WriteString(hv)
	}
	return res.String()
}
