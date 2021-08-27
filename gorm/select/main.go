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
// 原生sql查询

func main() {
	myDb := dbCon.GetDatabase()
	defer myDb.Close()

	var err error
	package_name := "com.os.airforce"
	apkName := ""

	//sql1.查询一条数据 一列或多列   scan
	var apkData ApkData
	sql := `select apk_res_type from apk where package_name = ?`
	err = myDb.Raw(sql, package_name).Scan(&apkData).Error
	if err != nil {
		log.Error(err)
	}
	fmt.Println("sql1:", apkData)

	// sql2.查询一条数据 一列    pluck
	apkRes := []string{}
	sql1 := `select apk_res_type from apk where package_name = ?`
	err = myDb.Raw(sql1, package_name).Pluck("apk_res_type", &apkRes).Error
	if err != nil {
		log.Error(err)
	}
	fmt.Println("sql2:", apkRes)

	// sql3查询多条
	appInfos := make([]NewApkInfo, 0)
	apkSql := `select apk_name,package_name,apk_res_type from apk order by apk_id desc limit 10`
	rows, err := myDb.Raw(apkSql).Rows()

	for rows.Next() {
		appInfo := NewApkInfo{}
		err := myDb.ScanRows(rows, &appInfo)
		if err != nil {
			continue // 此处一定要用continue 处理。要保证next 处理最后一条数据从而调用Rows.Close()把连接放回连接池复用。
		}
		appInfos = append(appInfos, appInfo)
	}
	fmt.Println("sql3:", appInfos)

	var packageNames []string
	var apkResTypes []string
	for rows.Next() { //  不可以处理两次，会导致没数据。
		var packageName string
		var apkResType string
		err := rows.Scan(&packageName, &apkResType)
		if err != nil {
			continue // 此处一定要用continue 处理。要保证next 处理最后一条数据从而调用Rows.Close()把连接放回连接池复用。
		}
		packageNames = append(packageNames, packageName)
		apkResTypes = append(apkResTypes, packageName)
	}
	fmt.Println("sql33:", packageNames, apkResTypes)

	// sql4查询多条
	info4 := make([]NewApkInfo, 0)
	sql4 := `select apk_name,package_name,apk_res_type from apk order by apk_id desc limit 10`
	err = myDb.Raw(sql4).Scan(&info4).Error
	if err != nil {
		log.Error(err)
	}
	fmt.Println("sql4:", info4)

	// sql5 取一个字段

	sql5 := `select apk_name from apk order by apk_id desc limit 1`
	err = myDb.Raw(sql5).Row().Scan(&apkName)
	if err != nil {
		log.Error(err)
	}
	fmt.Println("sql5:", apkName)

	//sql6.查询一条数据 一列或多列   scan
	apkData6 := model.ApkModel{}
	sql6 := `select * from apk where package_name = ?`
	err = myDb.Raw(sql6, package_name).First(&apkData6).Error
	if err != nil {
		log.Error(err)
	}
	fmt.Println("sql6:", apkData6)


	// sql7 修改一条数据

	sql7 := `update apk set apk_name ="123456" where id = 6`
	err = myDb.Exec(sql7).Error
	if err != nil {
		log.Error(err)
	}
}

// rows.Scan 和 db.ScanRows的区别， rows.Scan可以映射单个字段，也可以映射整个结构体。 db.ScanRows只能是结构体
