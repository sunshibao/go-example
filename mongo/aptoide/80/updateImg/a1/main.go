package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jmoiron/sqlx"

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

var DB *sqlx.DB

func main() {
	//wg := sync.WaitGroup{}
	//for i := 0; i < 9; i++ {
	//	//wg.Add(1)
	//	minId := i * 10000
	//	go func(id int) {
	//		//defer wg.Done()
	//		start(id)
	//	}(minId)
	//}
	////wg.Wait()
	start(0)
}

func start(id int) {
	uri := "root:Droi*#2021@tcp(18.192.114.175:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
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
		if err2 == nil && skip < 10000 {
			skip = 0 + limit*s
			err2 = GetApkList(id, skip, limit)
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

// AP应用图片资源表
type Ws80Image struct {
	ImageId     int    `gorm:"column:image_id" db:"image_id" json:"image_id" form:"image_id"`                     //图片ID
	ImageName   string `gorm:"column:image_name" db:"image_name" json:"image_name" form:"image_name"`             //图片名称
	ImageType   int    `gorm:"column:image_type" db:"image_type" json:"image_type" form:"image_type"`             //图片类型
	ImageWidth  int    `gorm:"column:image_width" db:"image_width" json:"image_width" form:"image_width"`         //图片宽
	ImageHeight int    `gorm:"column:image_height" db:"image_height" json:"image_height" form:"image_height"`     //图片高
	HdImageUrl  string `gorm:"column:hd_image_url" db:"hd_image_url" json:"hd_image_url" form:"hd_image_url"`     //图片访问地址
	NhdImageUrl string `gorm:"column:nhd_image_url" db:"nhd_image_url" json:"nhd_image_url" form:"nhd_image_url"` //图片服务器地址
	Display     int    `gorm:"column:display" db:"display" json:"display" form:"display"`                         //是否显示（0不显示  1显示）跟应用有关的不显示在图片管理列表上
	Status      int    `gorm:"column:status" db:"status" json:"status" form:"status"`                             //0:删除 1:可用
	Creator     int    `gorm:"column:creator" db:"creator" json:"creator" form:"creator"`                         //图片上传人ID
	CreateTime  int64  `gorm:"column:create_time" db:"create_time" json:"create_time" form:"create_time"`         //创建时间
	ModifyTime  int64  `gorm:"column:modify_time" db:"modify_time" json:"modify_time" form:"modify_time"`         //修改时间
	Language    string `gorm:"column:language" db:"language" json:"language" form:"language"`                     //语言
}

func GetApkList(minId, skip, limit int) (err error) {
	var wsId int
	sql1 := "select ws_id from ws80 where id>? and pull_status = 0 limit ?,? "
	err = DB.QueryRow(sql1, minId, skip, limit).Scan(&wsId)
	if err != nil {
		return err
	} else {
		fmt.Println("package:", "--------skip:", skip)
		GetApkDetail(wsId)

	}
	return nil
}

func GetApkDetail(wsId int) (err error) {

	url := fmt.Sprintf("https://ws75.aptoide.com/api/7/app/get/store_name=catappult/app_id=%d", wsId)

	resp, err := http.Get(url)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	var ws80 Ws80Detail
	json.Unmarshal(body, &ws80)
	if err != nil || ws80.Info.Status != "OK" {
		fmt.Println("err:", err, "----", wsId, "ws80.Info.Status:", ws80.Info.Status)
		return nil
	} else {

		insertWsData(ws80.Nodes)
		sql1 := "update ws80 set pull_status = 1 where ws_id = ? "
		DB.Exec(sql1, wsId)
	}

	return nil
}

func insertWsData(nodes Nodes) (err error) {
	newData := nodes.Meta.Data
	for _, v := range newData.Media.Screenshots {
		imgSql := `insert into ws80_image (image_name,image_type,image_width,image_height,hd_image_url,nhd_image_url,language) values (?,?,?,?,?,?,?)`
		exec, _ := DB.Exec(imgSql, nodes.Meta.Data.Package+"_Screenshots", 5, v.Width, v.Height, v.URL, v.URL, "ru")
		lastImgId, _ := exec.LastInsertId()
		fmt.Println("insertWsData", lastImgId)
		apkSql := `insert into ws80_apk_image (apk_id,language,image_id) values (?,?,?) ;`
		DB.Exec(apkSql, nodes.Meta.Data.ID, "ru", lastImgId)
	}
	return nil
}
