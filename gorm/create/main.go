package main

import (
	log "github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"myGo/gorm/dbCon"
	"myGo/gorm/model"
	"time"
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
	apkData.ApkId = 12
	apkData.ApkName = "creat1111111"
	apkData.ApkResType = "creat1111111"
	apkData.CreateTime = time.Now()
	apkData.ModifyTime = time.Now()

	err = myDb.Create(apkData).Error
	if err != nil {
		log.Error(err)
	}


}

// rows.Scan 和 db.ScanRows的区别， rows.Scan可以映射单个字段，也可以映射整个结构体。 db.ScanRows只能是结构体
