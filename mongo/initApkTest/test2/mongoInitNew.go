package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
	"time"

	"github.com/sunshibao/go-utils/util/gconv"
	"github.com/tencentyun/cos-go-sdk-v5"

	_ "github.com/go-sql-driver/mysql"

	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/jmoiron/sqlx"
)

type Ws80Detail struct {
	Info  Info  `json:"info"`
	Nodes Nodes `json:"nodes"`
}

type Info struct {
	Status string `json:"status"`
}

type Nodes struct {
	Meta Meta `json:"meta"`
}

type Meta struct {
	Data Data `json:"data"`
}

type Data struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Package   string    `json:"package"`
	Size      int       `json:"size"`
	Icon      string    `json:"icon"`
	Graphic   string    `json:"graphic"`
	Added     string    `json:"added"`
	Modified  string    `json:"modified"`
	Updated   string    `json:"updated"`
	Age       Age       `json:"age"`
	Developer Developer `json:"developer"`
	File      File      `json:"file"`
	Media     Media     `json:"media"`
	Stats     Stats     `json:"stats"`
	Appcoins  Appcoins  `json:"appcoins"`
}

type Age struct {
	Pegi string `json:"pegi"`
}

type Developer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Website string `json:"website"`
	Email   string `json:"email"`
	Privacy string `json:"privacy"`
}

type File struct {
	Vername         string   `json:"vername"`
	Vercode         int      `json:"vercode"`
	Md5Sum          string   `json:"md5sum"`
	Filesize        int      `json:"filesize"`
	Added           string   `json:"added"`
	Path            string   `json:"path"`
	Flags           Flags    `json:"flags"`
	UsedFeatures    []string `json:"used_features"`
	UsedPermissions []string `json:"used_permissions"`
}

type Flags struct {
	Votes []Votes `json:"votes"`
}

type Votes struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type Media struct {
	Keywords    []string      `json:"keywords"`
	Description string        `json:"description"`
	News        string        `json:"news"`
	Screenshots []Screenshots `json:"screenshots"`
}

type Screenshots struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type Stats struct {
	Downloads  int `json:"downloads"`
	Pdownloads int `json:"pdownloads"`
}

type Appcoins struct {
	Advertising bool `json:"advertising"`
	Billing     bool `json:"billing"`
}

var DB *sqlx.DB
var CosClient *cos.Client
var realNum int

func main() {
	//uri := "root:Droi*#2021@tcp(18.197.156.118:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	mysqldb, err := sqlx.Connect("mysql", uri)
	if err != nil {
		fmt.Println("mysql连接失败")
		mysqldb.Close()
	}
	DB = mysqldb

	//spec := "0 27 21 * * ?" //每天早上3：00：00执行一次
	//c := cron.New()
	//c.AddFunc(spec, gpCronFunc)
	//c.Start()
	//select {}

	gpCronFunc()
}

func gpCronFunc() {
	wg := sync.WaitGroup{}
	for i := 0; i <= 5; i++ {
		wg.Add(1)
		minId := i * 20000
		go func(id int) {
			defer wg.Done()
			start(id)
		}(minId)
	}
	wg.Wait()
	//start(0)
}

func start(minId int) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("UpLoadApkFile recover success")
		}
	}()
	skip := 0
	limit := 1
	s := 0
	var err2 error
	for {
		if err2 == nil && skip < 20000 {
			skip = 0 + limit*s
			err2 = shell(minId, skip, limit)
			s++
		} else {
			break
		}
	}
	return

}

type GpStruct struct {
	ApkId       int    `json:"apk_id"`
	PackageName string `json:"package_name"`
}

func shell(minId, skip, limit int) (err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("UpLoadApkFile recover success")
		}
	}()

	NewCosClient()
	var wsId int64
	var packageName string
	var updateTime string
	gpSql := `select ws_id,package,modified from ws80_detail where id>? limit ?,? ;`
	err = DB.QueryRow(gpSql, minId, skip, limit).Scan(&wsId, &packageName, &updateTime)

	if err != nil || wsId == 0 {
		log.Println("获取mysql数据失败:", err)
		return err
	}

	url := fmt.Sprintf("https://ws75.aptoide.com/api/7/app/get/store_name=catappult/app_id=%d", wsId)

	resp, err := http.Get(url)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var ws80 Ws80Detail
	json.Unmarshal([]byte(string(body)), &ws80)
	if err != nil {
		log.Printf("获取ap数据失败2 err:%v", err)
		return nil
	}

	if ws80.Info.Status == "FAIL" {
		apkSql2 := "insert into ap_update (ws_id,package,status)values(?,?,?)"
		DB.Exec(apkSql2, wsId, packageName, -1)
		log.Printf("获取ap数据已下架 aptoideId:%d", wsId)
		return nil
	}

	appDetails := ws80.Nodes.Meta.Data

	apUpdateTime := appDetails.Modified
	apCreateTime := appDetails.Added

	if apUpdateTime > updateTime {
		log.Printf("PackageName:%s, updateTime:%s,========== num: %d", appDetails.Package, apUpdateTime, skip)

		apkSql := "insert into ap_update (ws_id,package,create_time,update_time)values(?,?,?,?)"

		_, err := DB.Exec(apkSql, appDetails.ID, appDetails.Package, apCreateTime, apUpdateTime)
		if err != nil {
			log.Println("oz_apk insert fail:", err)
			return nil
		}
	} else {
		log.Printf("PackageName:%s, updateTime:%s ==========", appDetails.Package, apUpdateTime)
		return nil
	}

	return nil
}

//fileType ====  aptoide_icon aptoide_img aptoide_apk
func UpLoadImgFile(wsId int, packageName, fileUrl, fileType string) (pathUrl string, err error) {
	if fileType != "" {
		baseName := fileType + "/" + time.Now().Format("2006-01-02")
		split := strings.Split(fileUrl, "catappult/")
		iconName := gconv.String(wsId) + "-" + packageName + "-" + split[1]
		newName := baseName + "/" + iconName
		UploadCos(newName, fileUrl)
		pathUrl = "http://apk-ry.tt286.com/" + newName
	}
	return pathUrl, nil
}

//fileType ====  aptoide_icon aptoide_img aptoide_apk
func UpLoadApkFile(fileUrl, fileType string) (pathUrl string, err error) {
	if fileType != "" {
		baseName := fileType + "/" + time.Now().Format("2006-01-02")
		split := strings.Split(fileUrl, "catappult/")
		newName := baseName + "/" + split[1]
		UploadCos(newName, fileUrl)
		pathUrl = "http://apk-ry.tt286.com/" + newName
	}
	return pathUrl, nil
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

func UploadCos(name, filePath string) {
	c := CosClient
	filePath = strings.Replace(filePath, "https", "http", 1)
	resp, err := http.Get(filePath)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	reader := bufio.NewReaderSize(resp.Body, 32*1024)

	imgPath := "/data/service/cron/apFile/"

	fileName := path.Base(filePath)
	file, err := os.Create(imgPath + fileName)
	if err != nil {
		log.Printf("保存文件失败 err:%v", err)
	}
	_, err = io.Copy(file, reader)
	if err != nil {
		log.Println("下载文件到本地失败")
	}

	//本地上传
	_, err = c.Object.PutFromFile(context.Background(), name, imgPath+fileName, nil)
	if err != nil {
		log.Println(err)
	} else {
		err := os.Remove(imgPath + fileName)
		if err != nil {
			log.Println("删除失败:", imgPath+fileName)
		} else {
			log.Println("删除成功:", imgPath+fileName)
		}
	}
}
