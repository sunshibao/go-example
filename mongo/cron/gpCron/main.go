package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/robfig/cron"

	"github.com/sunshibao/go-utils/util/gconv"

	"github.com/disintegration/imageorient"
	"gopkg.in/mgo.v2/bson"

	_ "github.com/go-sql-driver/mysql"

	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2"
)

type BadgeForLegacyRating struct {
	Major string
}

type Ds struct {
	Low int `json:"low"`
}

type Dc struct {
	Low int `json:"low"`
}

type AppInfo struct {
	DeveloperName     string   `bson:"developerName"`
	DeveloperEmail    string   `bson:"developerEmail"`
	DeveloperWebsite  string   `bson:"developerWebsite"`
	VersionCode       int      `bson:"versionCode"`
	VersionString     string   `bson:"versionString"`
	Permission        []string `bson:"permission"`
	PackageName       string   `bson:"packageName"`
	AppType           string   `bson:"appType"`
	CategoryName      string   `bson:"categoryName"`
	RecentChangesHtml string   `bson:"recentChangesHtml"`
	InfoDownloadSize  Ds       `bson:"infoDownloadSize"`
	DownloadCount     Dc       `bson:"downloadCount"`
	InstallNotes      string   `bson:"installNotes"`
	InAppProduct      string   `bson:"inAppProduct"`
}

type Annotations struct {
	PrivacyPolicyUrl     string               `bson:"privacyPolicyUrl"`
	BadgeForLegacyRating BadgeForLegacyRating `bson:"badgeForLegacyRating"`
}

type Detail struct {
	AppDetails AppInfo `bson:"appDetails"`
}

type LangAppInfo struct {
	ID                     string      `bson:"id"`
	Type                   int         `bson:"type"`
	CategoryId             string      `bson:"CategoryId"`
	Title                  string      `bson:"title"`
	Creator                string      `bson:"creator"`
	DescriptionHtml        string      `bson:"descriptionHtml"`
	Icon                   string      `bson:"icon"`
	ImgList                []string    `bson:"imgList"`
	Detail                 Detail      `bson:"details"`
	ShareUrl               string      `bson:"shareUrl"`
	Annotations            Annotations `bson:"annotations"`
	DromotionalDescription string      `bson:"promotionalDescription"`
}

type MongoAppInfo struct {
	Package    string                 `bson:"package"`
	Language   map[string]LangAppInfo `bson:"language"`
	UpdateTime time.Time              `bson:"updateTime"`
	CreateTime time.Time              `bson:"createTime"`
	Status     string                 `bson:"status"`
	PicHost    string                 `bson:"pic_host"`
}

var DB *sqlx.DB
var MongoDB *mgo.Session
var realNum int

func main() {

	mongodb, err := mgo.Dial("mongodb://admin:Droi*#2021@10.0.0.8:27017/?replicaSet=market")
	//mongodb, err := mgo.Dial("mongodb://admin:Droi*#2021@43.131.92.130:27017/?replicaSet=market")
	if err != nil {
		mongodb.Close()
		return
	}
	defer mongodb.Close()
	MongoDB = mongodb

	uri := "root:Droi*#2021@tcp(18.197.156.118:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:@tcp(127.0.0.1:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"

	mysqldb, err := sqlx.Connect("mysql", uri)
	if err != nil {
		fmt.Println("mysql连接失败")
		mysqldb.Close()
	}

	DB = mysqldb

	spec := "0 0 20 * * ?" //每天早上3：00：00执行一次
	c := cron.New()
	c.AddFunc(spec, gpCronFunc)
	c.Start()
	select {}

	//gpCronFunc()
}

func gpCronFunc() {
	wg := sync.WaitGroup{}
	for i := 0; i <= 15; i++ {
		wg.Add(1)
		minId := i * 50000
		go func(id int) {
			defer wg.Done()
			start(id)
		}(minId)
	}
	wg.Wait()
	//start(0)
}

func start(minId int) {
	database := "package"
	skip := 0
	limit := 1
	s := 0
	var err2 error
	for {
		if err2 == nil && skip <= 50000 {
			skip = 0 + limit*s
			err2 = shell(MongoDB, database, minId, skip, limit)
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

func shell(mongodb *mgo.Session, database string, minId, skip, limit int) (err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("UpLoadApkFile recover success")
		}
	}()

	var apkId int
	var packageName string
	var versionCode int
	gpSql := `select apk_id,package_name,version_code from oz_apk where apk_id > ? and company_type = 1 limit ?,? ;`
	err = DB.QueryRow(gpSql, minId, skip, limit).Scan(&apkId, &packageName, &versionCode)

	if err != nil || packageName == "" {
		fmt.Println("获取gp数据失败:", err)
		return err
	}

	//查询mongodb
	c := mongodb.DB(database).C("info")
	var mongoAppInfos []MongoAppInfo

	c.Find(bson.M{"package": packageName}).All(&mongoAppInfos)

	if len(mongoAppInfos) <= 0 {
		fmt.Println("mongo没找到数据:", packageName)
		return err
	}
	var newApkId int
	var newVersionCode int
	gpSql2 := `select apk_id, version_code from oz_apk where company_type = 1 and package_name = ? order by version_code desc limit 1 ;`
	err = DB.QueryRow(gpSql2, packageName).Scan(&newApkId, &newVersionCode)
	if err != nil {
		log.Println("检查版本失败")
		return nil
	}

	for k, s := range mongoAppInfos {

		if s.PicHost == "" {
			continue
		}
		v := LangAppInfo{}
		appDetails := AppInfo{}

		v = s.Language["ru_RU"]
		appDetails = s.Language["ru_RU"].Detail.AppDetails

		gpVersion := appDetails.VersionCode
		gpUpdateTime := s.UpdateTime.Format("2006-01-02 15:04:05")
		gpCreateTime := s.CreateTime.Format("2006-01-02 15:04:05")

		if gpVersion > newVersionCode {
			log.Printf("PackageName:%s, updateTime:%s,========== num: %d", s.Package, gpUpdateTime, k+skip)
		} else {
			log.Printf("PackageName:%s, updateTime:%s ==========", s.Package, gpUpdateTime)
			continue
		}

		statusApk := gconv.Int64(s.Status) // 0下线 1正常
		apkType := 0                       // 0应用 1 游戏
		cateName := appDetails.CategoryName

		if appDetails.AppType == "GAME" {
			apkType = 1
			cateName = "GAME_" + cateName
		}

		//写入oz_apk表
		// 去标签
		v.DescriptionHtml = TrimHtml(v.DescriptionHtml)
		appDetails.RecentChangesHtml = TrimHtml(appDetails.RecentChangesHtml)

		if apkId > 175199 {
			// 过滤其他语言
			fl := filterLanguage(v.DescriptionHtml)
			if fl != true {
				fmt.Println("oz_apk 过滤其他语言 fail:", err)
				continue
			}
		}

		// 过滤敏感词
		fs := filterSensitiveWord(v.DescriptionHtml)
		if fs != true {
			fmt.Println("oz_apk 过滤敏感词 fail:", err)
			continue
		}

		if v.Annotations.BadgeForLegacyRating.Major != "" {
			tempAgeLimit := strings.TrimRight(v.Annotations.BadgeForLegacyRating.Major, "+")
			newAgeLimit := strings.TrimLeft(tempAgeLimit, "Rated for ")
			v.Annotations.BadgeForLegacyRating.Major = newAgeLimit
		}

		apkSql := "insert into oz_apk (package_name,apk_name,version_code,version_name,download_url,file_size,download_num,company,company_type,developer_email,developer_website,apk_res_type,apk_type,status,gp_down_url,age_limit,create_time,modify_time,set_create_time,set_update_time)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		appResult, err := DB.Exec(apkSql, s.Package, v.Title, appDetails.VersionCode, appDetails.VersionString, v.ShareUrl, appDetails.InfoDownloadSize.Low, appDetails.DownloadCount.Low, appDetails.DeveloperName, 1, appDetails.DeveloperEmail, appDetails.DeveloperWebsite, cateName, apkType, statusApk, v.ShareUrl, v.Annotations.BadgeForLegacyRating.Major, gpCreateTime, gpUpdateTime, time.Now(), time.Now())

		if err != nil {
			fmt.Println("oz_apk insert fail:", err)
			continue
		}

		apkId, _ := appResult.LastInsertId()

		permission := strings.Join(appDetails.Permission, ",")
		apkPerSql := "insert into oz_apk_permission (apk_id,permission)values(?,?)"
		_, err = DB.Exec(apkPerSql, apkId, permission)
		if err != nil {
			log.Printf("apkPerSql err:%v", err)
			return err
		}

		err = insertDB(s, "ru_RU", apkId, gpCreateTime, gpUpdateTime) //俄语
		if err != nil {
			continue
		}
	}

	return nil
}

// 多语言明细表插入
func insertDB(s MongoAppInfo, lang string, apkId int64, gpCreateTime, gpUpdateTime string) error {
	v := s.Language[lang]
	appDetails := s.Language[lang].Detail.AppDetails

	//h, _ := time.ParseDuration("-1h") 放到俄罗斯环境不用-8小时
	newIcon := "http://apk-ry.tt286.com" + s.PicHost

	lang = "ru"

	statusImg := 1
	// 更换图片连接
	v.Icon = strings.Replace(v.Icon, "https://play-lh.googleusercontent.com", newIcon, 1) + ".png"

	lang = "ru"

	//写入oz_image表
	iconSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`,`status`) VALUES (?,?,?,?,?,?);"
	iconResult, err := DB.Exec(iconSql, s.Package+"_icon", v.Icon, v.Icon, 50, lang, statusImg)
	if err != nil {
		log.Printf("oz_image 1 err:%v", err)
		return err
	}

	iconId, _ := iconResult.LastInsertId()

	inAppProduct := 0
	installNotes := 0
	if appDetails.InstallNotes != "" {
		installNotes = 1
	}
	if appDetails.InAppProduct != "" {
		inAppProduct = 1
	}

	// 去标签
	v.DescriptionHtml = TrimHtml(v.DescriptionHtml)
	appDetails.RecentChangesHtml = TrimHtml(appDetails.RecentChangesHtml)

	// 添加应用子表
	permission := strings.Join(appDetails.Permission, ",")
	apkDescSql := "insert into oz_apk_desc (apk_id,package_name,apk_name_lang,description,app_permission_desc,app_permission_url,ver_upt_des,language,png_icon_id,jpg_icon_id,in_app_product,install_notes,create_time,modify_time,set_create_time,set_update_time)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err = DB.Exec(apkDescSql, apkId, s.Package, v.Title, v.DescriptionHtml, permission, v.Annotations.PrivacyPolicyUrl, appDetails.RecentChangesHtml, lang, iconId, iconId, inAppProduct, installNotes, gpCreateTime, gpUpdateTime, time.Now(), time.Now())

	if err != nil {
		log.Printf("oz_apk_desc err:%v", err)
		return err
	}

	//图片去重，写入oz_image ,oz_apk_image
	//DelImage(apkId, newIcon, lang, s.Package, v.ImgList)
	//添加图片
	for _, val := range v.ImgList {
		val = strings.Replace(val, "https://play-lh.googleusercontent.com", newIcon, 1) + ".png"
		//写入oz_image表
		scrSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`) VALUES (?,?,?,?,?);"
		screResult, err := DB.Exec(scrSql, s.Package+"_Screenshots", val, val, 50, lang)

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

// 图片去重
func DelImage(apkId int64, newIcon, lang, package2 string, imageList []string) {
	mm := make(map[string][]string) //去重
	for _, v := range imageList {
		if v == "" {
			continue
		}
		v = strings.Replace(v, "https://play-lh.googleusercontent.com", newIcon, 1) + ".png"

		res, err := http.Get(v)
		if err != nil {
			fmt.Println("A error occurred!")
			continue
		}
		// defer后的为延时操作，通常用来释放相关变量
		defer res.Body.Close()

		pix, err := ioutil.ReadAll(res.Body)
		if err != nil {
			continue
		}
		fileCode := bytesToHexString(pix)

		mm[fileCode] = append(mm[fileCode], v)

	}
	for _, vv := range mm {
		for kkk, vvv := range vv {
			if kkk == 0 {
				//写入oz_image表
				scrSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`) VALUES (?,?,?,?,?);"
				screResult, err := DB.Exec(scrSql, package2+"_Screenshots", vvv, vvv, 50, lang)

				if err != nil {
					log.Printf("oz_image 3 err:%v", err)
					continue
				}
				//写入oz_apk_image表
				imageId, _ := screResult.LastInsertId()
				ssSql := "INSERT INTO oz_apk_image (`apk_id`, `image_id`,`language`) VALUES (?,?,?);"
				_, err = DB.Exec(ssSql, apkId, imageId, lang)
				if err != nil {
					log.Printf("oz_apk_image err:%v", err)
					continue
				}
			}
		}
	}
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

// 过滤标签
func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

// 正则过滤语言 只留英语
func filterLanguage(src string) bool {

	var filter = regexp.MustCompile("^[a-zA-Z0-9\\s\\./:●~!@#$%^&*\\+\\-(){}|<>=✔【】:\\\"?'：；‘’“”，。,、\\]\\[`《》]+$").MatchString

	return filter(src)
}

// 过滤敏感词
func filterSensitiveWord(src string) bool {
	sensitiveWord := []string{"Putin", "Xi Jinping", "communism", "Mao Zedong", "Shit", "blockhead", "Racist", "Nazi", "oppositionist", "Snout", "Muzzle", "Fucke", "stupid kakhah", "Bitch", "fuck", "fuck you", "Criminals and terrorists", "Muslims and terrorists", "Anti-Part", "Anti-Communist", "Smear China", "Slander the country", "Heroin", "pornography", "prostitute", "Sell oneself", "Pervert", "Asshole", "Yousuck", "kick ass", "bastard", "stupid jerk", "dick", "stupid idlot", "freak", "whore", "asshole", "Damn you", "fuck you", "Nerd", "bitch", "son of bitch", "suck for you SB", "Playing with fire ", "Pervert", "stupid", "idiot", "go to hell", "Shut up", "Bullshit", "God damn it", "SOB", "Drug dealing", "Dark we"}

	for _, v := range sensitiveWord {
		any := strings.Contains(src, v)
		if any == true {
			fmt.Println(v)
			return false
		}
	}
	return true
}

func getImgWH(url string) (width, height int) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return 0, 0
	}
	img, _, err := imageorient.Decode(resp.Body)

	if err != nil || img == nil {
		return 0, 0
	}
	return img.Bounds().Dx(), img.Bounds().Dy()
}
