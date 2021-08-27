package main

import (
	log "github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"myGo/gorm/dbCon"
	"myGo/gorm/model"
)

type NewApkInfo struct {
	ApkName     string `gorm:"column:apk_name" json:"apk_name"`         //应用名称
	PackageName string `gorm:"column:package_name" json:"package_name"` //应用包名
	ApkResType  string `gorm:"column:apk_res_type" json:"apk_res_type"` //应用资源类型

}

type ApkData struct {
	ApkResType string
}

func main() {
	myDb := dbCon.GetDatabase()
	defer myDb.Close()

	var err error
	//package_name := "com.os.airforce"

	// 更新单个属性（如果更改）
	apkData := model.ApkModel{}
	err = myDb.Model(apkData).Where("apk_id = ?", 1).Update("apk_name", "11111111").Error
	if err != nil {
		log.Error(err)
	}

	//使用`map`更新多个属性，只会更新这些更改的字段
	apkData2 := map[string]interface{}{"apk_name": "222222"}
	err = myDb.Model(apkData).Where("apk_id = ?", 2).Update(apkData2).Error
	if err != nil {
		log.Error(err)
	}

	//使用`struct`更新多个属性，只会更新这些更改的和非空白字段
	type ApkData3 struct {
		ApkName string `json:"apk_name"`
	}
	err = myDb.Model(apkData).Where("apk_id = ?", 3).Update(ApkData3{ApkName: "333333"}).Error
	if err != nil {
		log.Error(err)
	}

	//更新选择的字段
	//如果您只想在更新时更新或忽略某些字段，可以使用Select(只改select的字段), Omit(改除了omit的字段)
	apkData4 := map[string]interface{}{"apk_name": "222222","apk_res_type":"222222"}
	err = myDb.Model(apkData).Select("apk_res_type").Where("apk_id = ?", 4).Update(apkData4).Error
	if err != nil {
		log.Error(err)
	}
	apkData5 := map[string]interface{}{"apk_name": "222222","apk_res_type":"222222"}
	err = myDb.Model(apkData).Omit("apk_res_type").Where("apk_id = ?", 5).Update(apkData5).Error
	if err != nil {
		log.Error(err)
	}
}

// rows.Scan 和 db.ScanRows的区别， rows.Scan可以映射单个字段，也可以映射整个结构体。 db.ScanRows只能是结构体
