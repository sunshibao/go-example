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

	"github.com/sunshibao/go-utils/util/gconv"

	_ "github.com/go-sql-driver/mysql"

	"strings"

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
}

var DB *sqlx.DB
var realNum int

func GetDatabase() *sqlx.DB {
	return DB
}
func main() {
	wg := sync.WaitGroup{}
	for i := 0; i <= 10; i++ {
		minId := i * 40000
		wg.Add(1)
		go func(minId int) {
			defer wg.Done()
			start(minId)
		}(minId)
	}
	wg.Wait()
}

func start(minId int) {

	database := "package"
	mongodb, err := mgo.Dial("mongodb://admin:Droi*#2021@10.0.0.8:27017/?replicaSet=market")
	if err != nil {
		mongodb.Close()
		return
	}
	defer mongodb.Close()

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
		if err2 == nil && skip < minId+40000 {
			skip = minId + limit*s
			err2 = shell(mongodb, database, skip, limit)
			s++
		} else {
			break
		}
	}
	return

}

var gpNum = 0

func shell(mongodb *mgo.Session, database string, skip, limit int) (err error) {

	//查询mongodb
	c := mongodb.DB(database).C("info_220126")
	var mongoAppInfos []MongoAppInfo

	c.Find(nil).Sort("_id").Skip(skip).Limit(limit).All(&mongoAppInfos)

	if len(mongoAppInfos) <= 0 {
		fmt.Println("mongo没找到数据")
		return err
	}

	for k, s := range mongoAppInfos {

		if gpNum >= 100000 {
			log.Println("数据已达到10万=================当前包名:", s.Package)
			return err
		}

		v := LangAppInfo{}
		appDetails := AppInfo{}

		v = s.Language["en_US"]
		appDetails = s.Language["en_US"].Detail.AppDetails

		statusApk := gconv.Int(s.Status) //0下线 1正常
		apkType := 0                     // 0应用 1 游戏
		cateName := appDetails.CategoryName

		if appDetails.AppType == "GAME" {
			apkType = 1
			cateName = "GAME_" + cateName
		}

		log.Printf("PackageName:%s-----num:%d\n", s.Package, k+skip)
		//写入oz_apk表
		// 去标签
		v.DescriptionHtml = TrimHtml(v.DescriptionHtml)
		appDetails.RecentChangesHtml = TrimHtml(appDetails.RecentChangesHtml)

		// 过滤其他语言
		fl := filterLanguage(v.DescriptionHtml)

		if fl != true {
			continue
		}

		// 过滤敏感词
		fs := filterSensitiveWord(v.DescriptionHtml)
		if fs != true {
			continue
		}
		gpUpdateTime := s.UpdateTime.Format("2006-01-02 15:04:05")
		gpCreateTime := s.CreateTime.Format("2006-01-02 15:04:05")

		if v.Annotations.BadgeForLegacyRating.Major != "" {
			tempAgeLimit := strings.TrimRight(v.Annotations.BadgeForLegacyRating.Major, "+")
			newAgeLimit := strings.TrimLeft(tempAgeLimit, "Rated for ")
			v.Annotations.BadgeForLegacyRating.Major = newAgeLimit
		}

		apkSql := "insert into oz_apk (package_name,apk_name,version_code,version_name,download_url,file_size,download_num,company,company_type,developer_email,developer_website,apk_res_type,apk_type,status,gp_down_url,age_limit,create_time,modify_time,apk_source_sign,set_create_time,set_update_time)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		appResult, err := DB.Exec(apkSql, s.Package, v.Title, appDetails.VersionCode, appDetails.VersionString, v.ShareUrl, appDetails.InfoDownloadSize.Low, appDetails.DownloadCount.Low, appDetails.DeveloperName, 1, appDetails.DeveloperEmail, appDetails.DeveloperWebsite, cateName, apkType, statusApk, v.ShareUrl, v.Annotations.BadgeForLegacyRating.Major, gpCreateTime, gpUpdateTime, 4, time.Now(), time.Now())

		if err != nil {
			log.Printf("oz_apk err:%v", err)
			continue
		} else {

			apkId, _ := appResult.LastInsertId()
			// 添加权限
			permission := strings.Join(appDetails.Permission, ",")
			apkPerSql := "insert into oz_apk_permission (apk_id,permission)values(?,?)"
			_, err = DB.Exec(apkPerSql, apkId, permission)
			if err != nil {
				log.Printf("apkPerSql err:%v", err)
				return err
			}
			err = insertDB(s, "en_US", apkId, gpCreateTime, gpUpdateTime) //俄语
			if err != nil {
				continue
			}
		}
	}

	return nil
}

// 多语言明细表插入
func insertDB(s MongoAppInfo, lang string, apkId int64, gpCreateTime, gpUpdateTime string) error {
	v := s.Language[lang]
	appDetails := s.Language[lang].Detail.AppDetails

	//h, _ := time.ParseDuration("-1h") 放到俄罗斯环境不用-8小时
	newIcon := "http://apk-ry-tx.tt286.com/app_img/img_" + s.UpdateTime.Format("2006-01-02") + "/" + s.Package

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

	gpNum++
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
