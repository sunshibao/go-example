package main

import (
	log "github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"myGo/gorm/dbCon"
	"myGo/gorm/model"
)

// 事务
func main() {
	myDb := dbCon.GetDatabase()
	defer myDb.Close()

	var err error
	//package_name := "com.os.airforce"
	tx := myDb.Begin()
	// 更新单个属性（如果更改）
	apkData := model.ApkModel{}
	err = tx.Model(apkData).Where("apk_id = ?", 1).Update("apk_name", "AAAAAAAAAA").Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
	}

	//使用`map`更新多个属性，只会更新这些更改的字段
	apkData2 := map[string]interface{}{"apk_name": "BBBBBBBBBBBBB"}
	err = tx.Model(apkData).Where("apk_id = ?", 2).Update(apkData2).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
	}

	//使用`struct`更新多个属性，只会更新这些更改的和非空白字段
	type ApkData3 struct {
		ApkName string `json:"apk_name"`
	}
	err = tx.Model(apkData).Where("apk_id2 = ?", 3).Update(ApkData3{ApkName: "CCCCCCCCCCC"}).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
	}
	tx.Commit()

}

