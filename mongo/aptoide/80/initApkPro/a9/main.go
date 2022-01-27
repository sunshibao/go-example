package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"

	"strings"
)

type Ws80Detail struct {
	Id                  int    `gorm:"column:id" db:"id" json:"id" form:"id"`
	WsId                int    `gorm:"column:ws_id" db:"ws_id" json:"ws_id" form:"ws_id"`                                                 //应用ID
	Name                string `gorm:"column:name" db:"name" json:"name" form:"name"`                                                     //应用名
	Package             string `gorm:"column:package" db:"package" json:"package" form:"package"`                                         //包名
	CateType            string `gorm:"column:cate_type" db:"cate_type" json:"cate_type" form:"cate_type"`                                 //应用分类2级
	Size                string `gorm:"column:size" db:"size" json:"size" form:"size"`                                                     //文件大小
	Icon                string `gorm:"column:icon" db:"icon" json:"icon" form:"icon"`                                                     //icon
	Graphic             string `gorm:"column:graphic" db:"graphic" json:"graphic" form:"graphic"`                                         //icon2
	Added               string `gorm:"column:added" db:"added" json:"added" form:"added"`                                                 //创建时间
	Modified            string `gorm:"column:modified" db:"modified" json:"modified" form:"modified"`                                     //修改时间
	AgePegi             string `gorm:"column:age_pegi" db:"age_pegi" json:"age_pegi" form:"age_pegi"`                                     //年龄限制
	DeveloperName       string `gorm:"column:developer_name" db:"developer_name" json:"developer_name" form:"developer_name"`             //开发商名称
	DeveloperWebsite    string `gorm:"column:developer_website" db:"developer_website" json:"developer_website" form:"developer_website"` //开发商官网
	DeveloperEmail      string `gorm:"column:developer_email" db:"developer_email" json:"developer_email" form:"developer_email"`         //开发商邮箱
	DeveloperPrivacy    string `gorm:"column:developer_privacy" db:"developer_privacy" json:"developer_privacy" form:"developer_privacy"` //开发商隐私政策网址
	FileVername         string `gorm:"column:file_vername" db:"file_vername" json:"file_vername" form:"file_vername"`                     //版本名称
	FileVercode         int    `gorm:"column:file_vercode" db:"file_vercode" json:"file_vercode" form:"file_vercode"`                     //版本号
	FileMd5sum          string `gorm:"column:file_md5sum" db:"file_md5sum" json:"file_md5sum" form:"file_md5sum"`                         //MD5加密
	FileFilesize        int    `gorm:"column:file_filesize" db:"file_filesize" json:"file_filesize" form:"file_filesize"`                 //文件大小
	FileAdded           string `gorm:"column:file_added" db:"file_added" json:"file_added" form:"file_added"`                             //文件创建时间
	FilePath            string `gorm:"column:file_path" db:"file_path" json:"file_path" form:"file_path"`                                 //包下载地址
	FileFlagsVotes      string `gorm:"column:file_flags_votes" db:"file_flags_votes" json:"file_flags_votes" form:"file_flags_votes"`     //标签
	FileUsedFeatures    string `gorm:"column:file_used_features" db:"file_used_features" json:"file_used_features" form:"file_used_features"`
	FileUsedPermissions string `gorm:"column:file_used_permissions" db:"file_used_permissions" json:"file_used_permissions" form:"file_used_permissions"` //文件权限说明
	MediaKeywords       string `gorm:"column:media_keywords" db:"media_keywords" json:"media_keywords" form:"media_keywords"`                             //关键字
	MediaDescription    string `gorm:"column:media_description" db:"media_description" json:"media_description" form:"media_description"`                 //描述信息
	MediaNews           string `gorm:"column:media_news" db:"media_news" json:"media_news" form:"media_news"`                                             //新闻
	MediaScreenshots    string `gorm:"column:media_screenshots" db:"media_screenshots" json:"media_screenshots" form:"media_screenshots"`                 //图片截图
	StatsDownloads      string `gorm:"column:stats_downloads" db:"stats_downloads" json:"stats_downloads" form:"stats_downloads"`                         //下载量
	StatsPdownloads     string `gorm:"column:stats_pdownloads" db:"stats_pdownloads" json:"stats_pdownloads" form:"stats_pdownloads"`                     //真实下载量
	AppcoinsAdvertising string `gorm:"column:appcoins_advertising" db:"appcoins_advertising" json:"appcoins_advertising" form:"appcoins_advertising"`     //是否含有广告
	AppcoinsBilling     string `gorm:"column:appcoins_billing" db:"appcoins_billing" json:"appcoins_billing" form:"appcoins_billing"`                     //是否付费
	DownStatus          int    `gorm:"column:down_status" db:"down_status" json:"down_status" form:"down_status"`                                         //是否已拉取
	ImgPullStatus       int    `gorm:"column:img_pull_status" db:"img_pull_status" json:"img_pull_status" form:"img_pull_status"`                         //图片是否已经拉取
	ApkPullStatus       int    `gorm:"column:apk_pull_status" db:"apk_pull_status" json:"apk_pull_status" form:"apk_pull_status"`                         //包是否已经拉取
	ApkType             int    `gorm:"column:apk_type" db:"apk_type" json:"apk_type" form:"apk_type"`                                                     //0 应用 1 游戏
	SetCreateTime       string `gorm:"column:set_create_time" db:"set_create_time" json:"set_create_time" form:"set_create_time"`
	SetUpdateTime       string `gorm:"column:set_update_time" db:"set_update_time" json:"set_update_time" form:"set_update_time"`
}

var DB *sqlx.DB

func main() {
	start(80000)
}

func start(id int) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()
	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:Droi*#2021@tcp(127.0.0.1:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:sun18188@tcp(127.0.0.1:3306)/test13?charset=utf8mb4&parseTime=True&loc=Local"
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
	apDetailSql := "select * from ws80_detail where id>? limit ?,?"
	rows, err := DB.Queryx(apDetailSql, id, skip, limit)
	if err != nil {
		fmt.Println("sql 报错 sql:", apDetailSql, "--Id:", id, "--skip:", skip, "--limit:", limit)
		return err
	} else {
		for rows.Next() {
			var ws80Detail Ws80Detail
			err := rows.StructScan(&ws80Detail)
			if err != nil {
				fmt.Println("sql 报错 sql222:", err)
				continue
			}
			InsertApkPro(ws80Detail)
		}
	}
	return nil
}

var realNum int

func InsertApkPro(detail Ws80Detail) {

	// 检查包是否存在
	var apkId int64
	sql := `select apk_id from oz_apk where package_name = ? and company_type = 2`
	DB.Get(&apkId, sql, detail.Package)

	if apkId > 0 {

	} else {
		//写入oz_apk表
		realNum++
		log.Printf("PackageName:%s-----num:%d\n", detail.Package, realNum)
		detail.AgePegi = strings.Replace(detail.AgePegi, "PEGI-", "", 1)
		cateName := detail.CateType
		if detail.ApkType == 1 {
			cateName = "GAME_" + cateName
		}

		apkSql := "insert into oz_apk (package_name,apk_name,version_code,version_name,download_url,company,developer_email,developer_website,company_type,file_size,download_num,apk_res_type,apk_type,status,seo_key,age_limit,in_app_product,install_notes,app_permission_desc,app_permission_url,create_time,modify_time)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		appResult, err := DB.Exec(apkSql, detail.Package, detail.Name, detail.FileVercode, detail.FileVername, detail.FilePath, detail.DeveloperName, detail.DeveloperEmail, detail.DeveloperWebsite, 2, detail.FileFilesize, detail.StatsPdownloads, cateName, detail.ApkType, 1, detail.MediaKeywords, detail.AgePegi, detail.AppcoinsBilling, detail.AppcoinsAdvertising, detail.FileUsedPermissions, detail.FileUsedFeatures, detail.Added, detail.Modified)

		if err != nil {
			log.Printf("oz_apk err:%v", err)
		} else {
			apkId, _ = appResult.LastInsertId()

			// 添加权限

			apkPerSql := "insert into oz_apk_permission (apk_id,permission)values(?,?)"
			_, err = DB.Exec(apkPerSql, apkId, detail.FileUsedPermissions)

			if err != nil {
				log.Printf("apkPerSql err:%v", err)

			}
			err = insertDB(detail, apkId) //俄语
			if err != nil {
				log.Printf("insertDB err:%v", err)
			}
		}
	}
}

// 多语言明细表插入
func insertDB(detail Ws80Detail, apkId int64) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("捕获到的错误：%s\n", r)
		}
	}()

	statusImg := 1
	lang := "ru"

	//写入oz_image表
	iconSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`,`status`) VALUES (?,?,?,?,?,?);"
	iconResult, err := DB.Exec(iconSql, detail.Package+"_icon", detail.Icon, detail.Icon, 50, lang, statusImg)
	if err != nil {
		log.Printf("oz_image 1 err:%v", err)
		return err
	}
	iconId, _ := iconResult.LastInsertId()
	//写入oz_image表
	imgSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`,`status`) VALUES (?,?,?,?,?,?);"
	_, err = DB.Exec(imgSql, detail.Package+"_imgUrl", detail.Icon, detail.Icon, 50, lang, statusImg)

	if err != nil {
		log.Printf("oz_image 2 err:%v", err)
		return err
	}

	inAppProduct := detail.AppcoinsBilling
	installNotes := detail.AppcoinsAdvertising

	// 添加应用子表
	apkDescSql := "insert into oz_apk_desc (apk_id,package_name,apk_name_lang,description,app_permission_desc,app_permission_url,ver_upt_des,language,png_icon_id,jpg_icon_id,in_app_product,install_notes)values(?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err = DB.Exec(apkDescSql, apkId, detail.Package, detail.Name, detail.MediaDescription, detail.FileUsedPermissions, detail.FileUsedFeatures, "", lang, iconId, iconId, inAppProduct, installNotes)

	if err != nil {
		log.Printf("oz_apk_desc err:%v", err)
		return err
	}
	//添加图片
	imgList := strings.Split(detail.MediaScreenshots, ",")
	for _, val := range imgList {
		//写入oz_image表
		scrSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`) VALUES (?,?,?,?,?);"
		screResult, err := DB.Exec(scrSql, detail.Package+"_Screenshots", val, val, 50, lang)
		if err != nil {
			log.Printf("oz_image 3 err:%v", err)
			return err
		}
		//写入oz_apk_image表
		imageId, _ := screResult.LastInsertId()
		ssSql := "INSERT INTO oz_apk_image (`apk_id`, `image_id`,`language`) VALUES (?,?,?);"
		_, err = DB.Exec(ssSql, apkId, imageId, lang)
		if err != nil {
			log.Printf("oz_apk_image err:%v", err)
			return err
		}
	}
	return nil
}
