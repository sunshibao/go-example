package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
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
		go func(b int) {
			defer wg.Done()
			minId := b * 10000
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
		if err2 == nil && skip < minId+10000 {
			skip = minId + limit*s
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
	ApkId       int64  `db:"apk_id" json:"apk_id"`
	PackageName string `db:"package_name" json:"package_name"`
	FileSize    int64  `db:"file_size" json:"file_size"`
	DownloadUrl string `db:"download_url" json:"download_url"`
}

func GetApkList(skip, limit int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()

	apkInfo := ApkInfo{}
	sql3 := "select apk_id,package_name,file_size,download_url from oz_apk where company_type = 2 and status = 1 limit ?,? "
	err = DB.Get(&apkInfo, sql3, skip, limit)

	if apkInfo.DownloadUrl == "" {
		return nil
	}

	log.Println("hd_apk_url" + gconv.String(skip))

	UploadCos(apkInfo.ApkId, apkInfo.DownloadUrl)

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

func UploadCos(apkId int64, filePath string) bool {
	resp, err := http.Get(filePath)
	if err != nil {
		log.Println(err, "获取图片失败:"+filePath)
		return false
	}
	defer resp.Body.Close()

	reader := bufio.NewReaderSize(resp.Body, 0)

	imgPath := "/data/apkNew/"
	fileName := path.Base(filePath)
	file, err := os.Create(imgPath + fileName)
	if err != nil {
		fmt.Println(err, "url:"+filePath)
		return false
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d \n", written)

	apkFile := imgPath + fileName
	toolsPath := "/data/apksigner.jar"
	cmdStr := fmt.Sprintf("java -jar %s verify -v --print-certs %s ", toolsPath, apkFile)
	fileSha256 := apkSignInfo(cmdStr)
	fmt.Println(apkId, "=============sha256", fileSha256)

	if fileSha256 != "" {
		upSql := "update oz_apk set file_sha256 = ? where apk_id = ? "
		DB.Exec(upSql, fileSha256, apkId)

		err = os.Remove(imgPath + fileName)
		if err != nil {
			fmt.Println("删除失败:", imgPath+fileName)
		} else {
			fmt.Println("删除成功:", imgPath+fileName)
		}

	}

	return true
}

func apkSignInfo(cmsStr string) string {
	cmd := exec.Command("/bin/bash", "-c", cmsStr)
	fmt.Println("apkSignInfo", cmd)
	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("1111111", err)
		return ""
	}
	//执行命令
	if err := cmd.Start(); err != nil {
		fmt.Println("2222222", err)
		return ""
	}
	//读取所有输出
	resultByte, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("3333333", err)
		return ""
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println("4444444", err)
		return ""
	}

	result := string(resultByte)
	startStr := "Signer #1 certificate SHA-256 digest: "
	endStr := "Signer #1 certificate SHA-1 digest:"
	startIndex := strings.Index(result, startStr)
	endIndex := strings.Index(result, endStr)
	if startIndex+len(startStr) < endIndex {
		return result[startIndex+len(startStr) : endIndex-1]
	}
	return ""
}
