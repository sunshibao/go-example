package main

import (
	"fmt"
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
	package_name := "com.os.airforce"

	// sql.获取一条记录
	apkData := model.ApkModel{}
	err = myDb.Where("package_name = ?", package_name).First(&apkData).Error
	if err != nil {
		log.Error(err)
	}
	fmt.Println("sql:", apkData)

	// sql2.获取多条记录
	apkDatas := []model.ApkModel{}
	err = myDb.Where("id < ?", 3).Find(&apkDatas).Error
	if err != nil {
		log.Error(err)
	}
	fmt.Println("sql2:", apkDatas)

	// sql3.指定要从数据库检索的字段，默认情况下，将选择所有字段;

	apkDatas3 := []model.ApkModel{}
	err = myDb.Select("apk_id,apk_name").Where("id < ?", 3).Find(&apkDatas3).Error
	if err != nil {
		log.Error(err)
	}
	fmt.Println("sql3:", apkDatas3)

	//sql4、scan将结果扫描到另一个结构中
	apkDatas4 := []NewApkInfo{}
	err = myDb.Table("apk").Where("id < ?", 3).Scan(&apkDatas4).Error
	if err != nil {
		log.Error(err)
	}
	fmt.Println("sql4:", apkDatas4)

	// sql5.获取一条记录
	apkData5 := model.ApkModel{}
	err = myDb.Debug().Where("apk_id <10").Order("id desc").First(&apkData5).Error
	if err != nil {
		log.Error(err)
	}
	fmt.Println("sql5:", apkData5)

	//sql6.子查询
	apkData6 := model.ApkModel{}
	err = myDb.Debug().Where("apk_id <?", myDb.Table("apk").Select("count(*)").Where("id<?", 10).QueryExpr()).Order("id desc").First(&apkData6).Error
	if err != nil {
		log.Error(err)
	}
	fmt.Println("sql5:", apkData6)
}

// rows.Scan 和 db.ScanRows的区别， rows.Scan可以映射单个字段，也可以映射整个结构体。 db.ScanRows只能是结构体

//db.Where("amount > ?", DB.Table("orders").Select("AVG(amount)").Where("state = ?", "paid").QueryExpr()).Find(&orders)
// SELECT * FROM "orders"  WHERE "orders"."deleted_at" IS NULL AND (amount > (SELECT AVG(amount) FROM "orders"  WHERE (state = 'paid')));
