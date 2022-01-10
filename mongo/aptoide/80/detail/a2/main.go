package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
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
	Developer Developer `json:"developer"`
	File      File      `json:"file"`
	Media     Media     `json:"media"`
	Stats     Stats     `json:"stats"`
	Appcoins  Appcoins  `json:"appcoins"`
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

var DB *gorm.DB

func main() {

	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
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

type Ws80Temp struct {
	WsID     int    `json:"ws_id"`
	Name     string `json:"name"`
	Package  string `json:"package"`
	CateType string `json:"cate_type"`
}

type MysqlWs80Detail struct {
	WsId                int    `gorm:"column:ws_id" db:"ws_id" json:"ws_id" form:"ws_id"`
	Name                string `gorm:"column:name" db:"name" json:"name" form:"name"`
	Package             string `gorm:"column:package" db:"package" json:"package" form:"package"`
	CateType            string `gorm:"column:cate_type" db:"cate_type" json:"cate_type" form:"cate_type"`
	Size                int    `gorm:"column:size" db:"size" json:"size" form:"size"`
	Icon                string `gorm:"column:icon" db:"icon" json:"icon" form:"icon"`
	Graphic             string `gorm:"column:graphic" db:"graphic" json:"graphic" form:"graphic"`
	Added               string `gorm:"column:added" db:"added" json:"added" form:"added"`
	Modified            string `gorm:"column:modified" db:"modified" json:"modified" form:"modified"`
	Updated             string `gorm:"column:updated" db:"updated" json:"updated" form:"updated"`
	DeveloperName       string `gorm:"column:developer_name" db:"developer_name" json:"developer_name" form:"developer_name"`
	DeveloperWebsite    string `gorm:"column:developer_website" db:"developer_website" json:"developer_website" form:"developer_website"`
	DeveloperEmail      string `gorm:"column:developer_email" db:"developer_email" json:"developer_email" form:"developer_email"`
	DeveloperPrivacy    string `gorm:"column:developer_privacy" db:"developer_privacy" json:"developer_privacy" form:"developer_privacy"`
	FileVername         string `gorm:"column:file_vername" db:"file_vername" json:"file_vername" form:"file_vername"`
	FileVercode         int    `gorm:"column:file_vercode" db:"file_vercode" json:"file_vercode" form:"file_vercode"`
	FileMd5sum          string `gorm:"column:file_md5sum" db:"file_md5sum" json:"file_md5sum" form:"file_md5sum"`
	FileFilesize        int    `gorm:"column:file_filesize" db:"file_filesize" json:"file_filesize" form:"file_filesize"`
	FileAdded           string `gorm:"column:file_added" db:"file_added" json:"file_added" form:"file_added"`
	FilePath            string `gorm:"column:file_path" db:"file_path" json:"file_path" form:"file_path"`
	FileFlagsVotes      string `gorm:"column:file_flags_votes" db:"file_flags_votes" json:"file_flags_votes" form:"file_flags_votes"`
	FileUsedFeatures    string `gorm:"column:file_used_features" db:"file_used_features" json:"file_used_features" form:"file_used_features"`
	FileUsedPermissions string `gorm:"column:file_used_permissions" db:"file_used_permissions" json:"file_used_permissions" form:"file_used_permissions"`
	MediaKeywords       string `gorm:"column:media_keywords" db:"media_keywords" json:"media_keywords" form:"media_keywords"`
	MediaDescription    string `gorm:"column:media_description" db:"media_description" json:"media_description" form:"media_description"`
	MediaNews           string `gorm:"column:media_news" db:"media_news" json:"media_news" form:"media_news"`
	MediaScreenshots    string `gorm:"column:media_screenshots" db:"media_screenshots" json:"media_screenshots" form:"media_screenshots"`
	StatsDownloads      int    `gorm:"column:stats_downloads" db:"stats_downloads" json:"stats_downloads" form:"stats_downloads"`
	StatsPdownloads     int    `gorm:"column:stats_pdownloads" db:"stats_pdownloads" json:"stats_pdownloads" form:"stats_pdownloads"`
	AppcoinsAdvertising int    `gorm:"column:appcoins_advertising" db:"appcoins_advertising" json:"appcoins_advertising" form:"appcoins_advertising"`
	AppcoinsBilling     int    `gorm:"column:appcoins_billing" db:"appcoins_billing" json:"appcoins_billing" form:"appcoins_billing"`
	DownStatus          int    `gorm:"column:down_status" db:"down_status" json:"down_status" form:"down_status"`
}

func GetApkList(skip, limit int) (err error) {
	sql1 := "select ws_id,name,package,cate_type from ws80 where id>20000 limit ?,? "
	rows, err := DB.Raw(sql1, skip, limit).Rows()
	if err != nil {
		return err
	} else {
		for rows.Next() {
			var ws80Temp Ws80Temp
			err := DB.ScanRows(rows, &ws80Temp)
			if err != nil {
				continue
			}
			fmt.Println("package:", ws80Temp.Package, "--------skip:", skip)
			GetApkDetail(ws80Temp)
		}
	}
	return nil
}

func GetApkDetail(temp Ws80Temp) (err error) {

	url := fmt.Sprintf("https://ws75.aptoide.com/api/7/app/get/store_name=catappult/app_id=%d", temp.WsID)

	resp, err := http.Get(url)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var ws80 Ws80Detail
	json.Unmarshal([]byte(string(body)), &ws80)
	if err != nil || ws80.Info.Status != "OK" {
		sql1 := "update ws80 set pull_status = 1 where ws_id = ? "
		DB.Exec(sql1, temp.WsID)
	} else {
		insertWsData(ws80.Nodes, temp)
	}

	return nil
}

func insertWsData(nodes Nodes, ws80Temp Ws80Temp) (err error) {
	newData := nodes.Meta.Data
	//downUrl := nodes.Meta.Data.File.Path
	votesSlice := []string{}
	for _, v := range newData.File.Flags.Votes {
		votesSlice = append(votesSlice, v.Type)
	}

	imgSlice := []string{}
	for _, v := range newData.Media.Screenshots {
		imgSlice = append(imgSlice, v.URL)
	}

	fileFlagsVotes := strings.Join(votesSlice, ",")
	fileUsedFeatures := strings.Join(newData.File.UsedFeatures, ",")
	fileUsedPermissions := strings.Join(newData.File.UsedPermissions, ",")
	mediaKeywords := strings.Join(newData.Media.Keywords, ",")
	mediaScreenshots := strings.Join(imgSlice, ",")

	appcoinsAdvertising := 0
	appcoinsBilling := 0
	if newData.Appcoins.Advertising {
		appcoinsAdvertising = 1
	}
	if newData.Appcoins.Billing {
		appcoinsBilling = 1
	}

	newMysql := MysqlWs80Detail{
		newData.ID,
		newData.Name,
		newData.Package,
		ws80Temp.CateType,
		newData.Size,
		newData.Icon,
		newData.Graphic,
		newData.Added,
		newData.Modified,
		newData.Updated,
		newData.Developer.Name,
		newData.Developer.Website,
		newData.Developer.Email,
		newData.Developer.Privacy,
		newData.File.Vername,
		newData.File.Vercode,
		newData.File.Md5Sum,
		newData.File.Filesize,
		newData.File.Added,
		newData.File.Path,
		fileFlagsVotes,
		fileUsedFeatures,
		fileUsedPermissions,
		mediaKeywords,
		newData.Media.Description,
		newData.Media.News,
		mediaScreenshots,
		newData.Stats.Downloads,
		newData.Stats.Pdownloads,
		appcoinsAdvertising,
		appcoinsBilling,
		0,
	}
	DB.Table("ws80_detail").Create(newMysql)
	//go DownApk(downUrl, ws80Temp.WsID)
	return nil
}

func DownApk(apkUrl string, wsId int) {

	apkPath := "/Users/sunshibao/Desktop/apkDown/"

	fileName := path.Base(apkUrl)

	res, err := http.Get(apkUrl)
	if err != nil {
		sql1 := "update ws80_detail set down_status = 1 where ws_id = ? "
		DB.Exec(sql1, wsId)
		fmt.Println("A error occurred!---", apkUrl)
		return
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	file, err := os.Create(apkPath + fileName)
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)

	written, _ := io.Copy(writer, reader)
	fmt.Printf("Total length: %d", written)
}
