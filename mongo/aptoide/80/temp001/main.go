package main

import (
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sqlx.DB

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i <= 9; i++ {
		wg.Add(1)
		minId := i * 10000
		go func(id int) {
			defer wg.Done()
			start(id)
		}(minId)
	}
	wg.Wait()
}

func start(id int) {

	//uri := "usr_dev:6RqfI^G^QaFLh@eqk*Z@tcp(data-sql1.ry.cn:3306)/ry_market?charset=utf8mb4&parseTime=True&loc=Local"
	uri := "root:Droi*#2021@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&loc=Local"
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

	var wsId int
	apDetailSql := `select ws_id from ws80 where id>? and apk_type = 1 limit ?,?`
	err = DB.Get(&wsId, apDetailSql, id, skip, limit)
	fmt.Println(wsId, "-------", skip)
	upSql := `update ws80_detail set apk_type = 1 where ws_id = ?`
	DB.Exec(upSql, wsId)

	return nil
}
