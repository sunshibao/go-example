package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

var DB *gorm.DB

func GetDatabase() *gorm.DB {
	return DB
}

func main() {

	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:tyd*#2016@tcp(192.168.1.152:3306)/ry_market_examine?charset=utf8mb4&parseTime=True&loc=Local"
	//uri := "root:sun18188@tcp(127.0.0.1:3306)/test13?charset=utf8mb4&parseTime=True&loc=Local"

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
		if err2 == nil && skip < 4300 {
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
	apkId := 0
	sql := "select apk_id from tem_apk where apk_id>808 limit ?,?"
	DB.Raw(sql, skip, limit).Row().Scan(&apkId)

	if apkId == 0 {
		return errors.New("无数据")
	}
	fmt.Printf("apkId:%d-----num:skip:%d \n", apkId, skip)

	sql2 := "select a.apk_id,a.language,a.image_id,b.hd_image_url from oz_apk_image a left join oz_image b ON a.image_id = b.image_id where a.apk_id = ?"

	rows, err := DB.Raw(sql2, apkId).Rows()
	if err != nil {
		log.Printf("oz_image 1 err:%v", err)
		return err
	}

	menuMap := []ImageInfo{}
	for rows.Next() {
		info := ImageInfo{}
		err := DB.ScanRows(rows, &info)
		if err != nil {
			continue
		}

		menuMap = append(menuMap, info)
	}

	if len(menuMap) <= 0 {
		return nil
	}

	//处理img数据
	imgList := make(map[string][]ImageInfo)
	for _, v := range menuMap {
		imgList[v.Language] = append(imgList[v.Language], v)
	}
	for lang, val := range imgList {
		mm := make(map[string][]int) //去重
		for _, v := range val {
			if v.HdImageUrl == "" {
				iconSql4 := "delete from oz_apk_image where image_id = ? and language = ?;"
				DB.Exec(iconSql4, v.ImageId, lang)
				continue
			}
			res, err := http.Get(v.HdImageUrl)
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

			mm[fileCode] = append(mm[fileCode], v.ImageId)

		}
		for _, vv := range mm {
			for kkk, vvv := range vv {
				if kkk == 0 {
					continue
				}
				//写入oz_image表
				iconSql3 := "delete from oz_image where image_id = ? and language = ?;"
				DB.Exec(iconSql3, vvv, lang)

				iconSql4 := "delete from oz_apk_image where image_id = ? and language = ?;"
				DB.Exec(iconSql4, vvv, lang)

			}
		}

	}

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
