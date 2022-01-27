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

	"gopkg.in/mgo.v2/bson"

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
	DeveloperWebsite  string   `Bson:"developerWebsite"`
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
	Status     int                    `bson:"status"`
}

var DB *sqlx.DB
var realNum int

func GetDatabase() *sqlx.DB {
	return DB
}
func main() {
	wg := sync.WaitGroup{}
	for i := 0; i <= 12; i++ {
		wg.Add(1)
		minId := i * 20000
		go func(id int) {
			defer wg.Done()
			start(id)
		}(minId)
	}
	wg.Wait()
}

func start(minId int) {

	database := "package"
	mongodb, err := mgo.Dial("mongodb://admin:Droi*#2021@43.131.69.147:27017,43.131.92.130:27017,162.62.196.12:27017/?replicaSet=market")
	if err != nil {
		mongodb.Close()
		return
	}
	defer mongodb.Close()

	uri := "root:Droi*#2021@tcp(18.192.114.175:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"

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
		if err2 == nil && skip < 20000 {
			skip = 0 + limit*s
			err2 = shell(mongodb, database, minId, skip, limit)
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
	var gpStruct GpStruct
	gpSql := `select apk_id,package_name from oz_apk where apk_id > ? and status != -1 and company_type = 1 limit ?,? ;`
	err = DB.QueryRowx(gpSql, minId, skip, limit).Scan(&gpStruct)
	if err != nil || gpStruct.PackageName == "" {
		fmt.Println("获取gp数据失败:", err)
		return err
	}
	//查询mongodb
	c := mongodb.DB(database).C("info")
	var mongoAppInfos []MongoAppInfo

	c.Find(bson.M{"package": gpStruct.PackageName}).All(&mongoAppInfos)

	if len(mongoAppInfos) <= 0 {
		fmt.Println("mongo没找到数据:", gpStruct.PackageName)
		return err
	}

	for k, s := range mongoAppInfos {
		v := LangAppInfo{}
		appDetails := AppInfo{}
		if gpStruct.ApkId < 175199 {
			v = s.Language["ru_RU"]
			appDetails = s.Language["ru_RU"].Detail.AppDetails
		}

		statusApk := s.Status //0下线 1正常
		apkType := 0          // 0应用 1 游戏
		cateName := appDetails.CategoryName

		if appDetails.AppType == "GAME" {
			apkType = 1
			cateName = "GAME_" + cateName
		}
		log.Printf("PackageName1111111:%s-----num:%d\n", s.Package, k+skip)

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

		realNum++
		log.Printf("PackageName333333333:%s-----num:%d\n", s.Package, realNum)
		if realNum > 400 {
			return err
		}

		apkSql := "update oz_apk set apk_name = ?,version_code = ?,version_name = ?,download_url = ?,company = ?,developer_email=?,file_size = ?,developer_website = ?,download_num = ? ,apk_res_type=?,apk_type=?,status =?,gp_down_url =?,age_limit =? where apk_id = ?"
		appResult, err := DB.Exec(apkSql, v.Title, appDetails.VersionCode, appDetails.VersionString, v.ShareUrl, appDetails.DeveloperName, 1, appDetails.InfoDownloadSize.Low, appDetails.DownloadCount.Low, cateName, apkType, statusApk, v.ShareUrl, v.Annotations.BadgeForLegacyRating.Major)

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
			err = insertDB(s, "en_US", apkId) //俄语
			if err != nil {
				continue
			}
		}

	}

	return nil
}

// 多语言明细表插入
func insertDB(s MongoAppInfo, lang string, apkId int64) error {
	v := s.Language[lang]
	appDetails := s.Language[lang].Detail.AppDetails

	h, _ := time.ParseDuration("-1h")
	newIcon := "https://gp-image-1308128293.cos.eu-moscow.myqcloud.com/app_img/img_" + s.UpdateTime.Add(8*h).Format("2006-01-02") + "/" + s.Package

	if lang == "zh_CN" {
		lang = "zh_cn"
	} else if lang == "en_US" {
		lang = "en"
	} else if lang == "ru_RU" {
		lang = "ru"
	} else if lang == "be_BY" {
		lang = "be"
	} else if lang == "uk_UA" {
		lang = "uk"
	} else if lang == "kk_KZ" {
		lang = "kk"
	} else if lang == "ka_GE" {
		lang = "ka"
	} else if lang == "uz_UZ" {
		lang = "uz"
	} else {
		lang = "zh_cn"
	}

	cateName := appDetails.CategoryName
	if appDetails.AppType == "GAME" {
		cateName = "GAME_" + cateName
	}

	statusApkDesc := 1
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
	//写入oz_image表
	imgSql := "INSERT INTO oz_image (`image_name`, `hd_image_url`, `nhd_image_url`,`image_type`,`language`,`status`) VALUES (?,?,?,?,?,?);"
	_, err = DB.Exec(imgSql, s.Package+"_imgUrl", v.Icon, v.Icon, 50, lang, statusImg)

	if err != nil {
		log.Printf("oz_image 2 err:%v", err)
		return err
	}

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
	apkDescSql := "insert into oz_apk_desc (apk_id,package_name,apk_name_lang,description,app_permission_desc,app_permission_url,ver_upt_des,language,png_icon_id,jpg_icon_id,status,in_app_product,install_notes,apk_res_type)values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err = DB.Exec(apkDescSql, apkId, s.Package, v.Title, v.DescriptionHtml, permission, v.Annotations.PrivacyPolicyUrl, appDetails.RecentChangesHtml, lang, iconId, iconId, statusApkDesc, inAppProduct, installNotes, cateName)

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
