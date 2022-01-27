package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/tencentyun/cos-go-sdk-v5"
)

var CosClient *cos.Client
var DB *gorm.DB

func main() {
	NewCosClient()
	UploadCos()
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

func UploadCos() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()
	c := CosClient

	basePath := "http://apk-ry.tt286.com/"
	imgPath := "/Users/sunshibao/Desktop/apkNew/"
	fileName := "com-rkstudio-mountain-3d-driving-challenge-7-47042375-6ac82cbc34394e9a03884f205a56dda9.apk"

	name := "aptoide_apk/" + fileName

	//本地上传
	_, err := c.Object.PutFromFile(context.Background(), name, imgPath+fileName, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("上传成功:url:" + basePath + name)
}
