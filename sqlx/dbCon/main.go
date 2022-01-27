package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// 建立连接
func initDB() (err error) {
	dsn := "ruan:password@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True"
	// Connect包含了open和ping方法，也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}

	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)
	return
}

func main() {
	if err := initDB(); err != nil {
		fmt.Printf("init DB failed, err:%v\n", err)
		return
	}
	fmt.Println("init DB success...")
}

// 基本查询
type user struct { // 结构体打tag，数据库对应字段
	ID   int    `db:"id"` // 同类型放一起对内存地址更有顺序
	Age  int    `db:"age"`
	Name string `db:"name"`
}

// 查询单条数据示例 结构体，Get
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"
	var u user
	err := db.Get(&u, sqlStr, 1) // 结构体接受指针， 查询语句， 占位符值
	if err != nil {
		fmt.Printf("get failed err:%v\n", err)
		return
	}
	fmt.Printf("id:%d name:%s age:%d\n", u.ID, u.Name, u.Age)

}

// 查询多条数据示例 切片，Select
func queryMultiRowDemo() {
	sqlStr := "select id, name, age from user where id > ?"
	var users []user
	err := db.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	fmt.Printf("users:%#v\n", users)
}

// sqlx中的exec方法与原生sql中的exec使用基本一致
// 插入数据 Exec
func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "朝阳", 18)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}

// sqlx中的exec方法与原生sql中的exec使用基本一致
// 更新数据 Exec
func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 18, 1)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}

	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

// sqlx中的exec方法与原生sql中的exec使用基本一致
// 删除数据 Exec
func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 2)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

// NamedQuery方法用来绑定SQL语句与结构体或map中的同名字段, 查询。
// NamedQuery :name
func namedQuery() {
	sqlStr := "SELECT * FROM user WHERE name=:name"
	// 使用map做命名查询
	rows, err := db.NamedQuery(sqlStr, map[string]interface{}{"name": "西瓜"})
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}

	// 使用结构体命名查询，根据结构体字段的 db tag进行映射
	u := user{
		Name: "西瓜",
	}
	rows, err = db.NamedQuery(sqlStr, u)
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}
}

// NamedExec方法用来绑定SQL语句与结构体或map中的同名字段。
// NamedExec :name :age
func insertUserDemo() (err error) {
	sqlStr := "INSERT INTO user (name,age) VALUES (:name,:age)"

	// 使用map做命名插入
	_, err = db.NamedExec(sqlStr,
		map[string]interface{}{
			"name": "西瓜",
			"age":  23,
		})
	if err != nil {
		fmt.Printf("db.NamedExec failed, err:%v\n", err)
		return
	}

	// 使用结构体命名插入，根据结构体字段的 db tag进行映射
	u := user{
		Name: "西瓜",
		Age:  18,
	}
	_, err = db.NamedExec(sqlStr, u)
	if err != nil {
		fmt.Printf("db.NamedExec failed, err:%v\n", err)
		return
	}
	return
}

//对于事务操作，我们可以使用sqlx中提供的db.Beginx()和tx.Exec()方法。
// 事务操作
func transactionDemo() (err error) {
	tx, err := db.Beginx() // 开启事务
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // 回滚后，panic
		} else if err != nil {
			fmt.Println("rollback")
			tx.Rollback() // err为非零，不要改变它
		} else {
			err = tx.Commit() // err is nil，如果Commit返回错误更新错误
			fmt.Println("commit")
		}
	}()

	sqlStr1 := "Update user set age=23 where id=?"

	rs, err := tx.Exec(sqlStr1, 1)
	if err != nil {
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	sqlStr2 := "Update user set age=18 where id=?"
	rs, err = tx.Exec(sqlStr2, 3)
	if err != nil {
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr2 failed")
	}
	return err
}

//sqlx.In
//sqlx.In是sqlx提供的一个非常方便的函数。
//
//sqlx.In的批量插入
//bindvars（绑定变量）
//查询占位符?在内部称为bindvars（查询占位符）,它非常重要。你应该始终使用它们向数据库发送值，因为它们可以防止SQL注入攻击。database/sql不尝试对查询文本进行任何验证；它与编码的参数一起按原样发送到服务器。除非驱动程序实现一个特殊的接口，否则在执行之前，查询是在服务器上准备的。因此bindvars是特定于数据库的:
//
//MySQL中使用?
//PostgreSQL使用枚举的$1、$2等bindvar语法
//SQLite中?和$1的语法都支持
//Oracle中使用:name的语法
//bindvars的一个常见误解是，它们用来在sql语句中插入值。它们其实仅用于参数化，不允许更改SQL语句的结构。例如，使用bindvars尝试参数化列或表名将不起作用

// ？不能用来插入表名（做SQL语句中表名的占位符）
//db.Query("SELECT * FROM ?", "mytable")

// ？也不能用来插入列名（做SQL语句中列名的占位符）
//db.Query("SELECT ?, ? FROM people", "name", "location")

//拼接语句实现批量插入
// BatchInsertUsers 自行构造批量插入的语句
func BatchInsertUsers(users []*user) error {
	// 存放 (?, ?) 的slice
	valueStrings := make([]string, 0, len(users))
	// 存放values的slice
	valueArgs := make([]interface{}, 0, len(users)*2)
	// 遍历users准备相关数据
	for _, u := range users {
		// 此处占位符要与插入值的个数对应
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, u.Name)
		valueArgs = append(valueArgs, u.Age)
	}
	// 自行拼接要执行的具体语句
	stmt := fmt.Sprintf("INSERT INTO user (name, age) VALUES %s",
		strings.Join(valueStrings, ","))
	_, err := db.Exec(stmt, valueArgs...)
	return err
}

// 批量拼接插入
func ExecBatchInsert() {
	u1 := user{Name: "西瓜", Age: 18}
	u2 := user{Name: "yy", Age: 23}
	u3 := user{Name: "kk", Age: 24}
	users := []*user{&u1, &u2, &u3}
	err := BatchInsertUsers(users)
	if err != nil {
		fmt.Printf("BatchInsertUsers failed, err:%v\n", err)
	}
}

//sqlx.In实现批量插入
//前提是需要我们的结构体实现driver.Valuer接口：

func (u user) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}

//sqlx.In实现批量插入代码

// BatchInsertUsers2 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
func BatchInsertUsers2(users []interface{}) error {
	query, args, _ := sqlx.In(
		"INSERT INTO user (name, age) VALUES (?), (?), (?)",
		users..., // 如果arg实现了 driver.Valuer, sqlx.In 会通过调用 Value()来展开它
	)
	fmt.Println(query) // 查看生成的querystring  // INSERT INTO user (name, age) VALUES (?, ?), (?, ?), (?, ?)
	fmt.Println(args)  // 查看生成的args // [西瓜 18 yy 23 kk 24]
	_, err := db.Exec(query, args...)
	return err
}

// sqlx.In实现批量插入代码
func ExecBatchInsertSqlxIn() {
	u1 := user{Name: "西瓜", Age: 18}
	u2 := user{Name: "yy", Age: 23}
	u3 := user{Name: "kk", Age: 24}
	users := []interface{}{u1, u2, u3}
	err := BatchInsertUsers2(users)
	if err != nil {
		fmt.Printf("BatchInsertUsers2 failed, err:%v\n", err)
	}
}

//NamedExec实现批量插入
//注意 ：该功能需github.com/jmoiron/sqlx1.3.1版本以上，sql语句最后不能空格和;

// BatchInsertUsers3 使用NamedExec实现批量插入
func BatchInsertUsers3(users []*user) error {
	_, err := db.NamedExec("INSERT INTO user (name, age) VALUES (:name, :age)", users)
	return err
}

func ExecBatchInsertNamedExec() {
	u1 := user{Name: "西瓜", Age: 18}
	u2 := user{Name: "yy", Age: 23}
	u3 := user{Name: "kk", Age: 24}
	users := []*user{&u1, &u2, &u3}
	err := BatchInsertUsers3(users)
	if err != nil {
		fmt.Printf("BatchInsertUsers3 failed, err:%v\n", err)
	}
}

//sqlx.In的查询
//sqlx查询语句中实现In查询和FIND_IN_SET函数。即实现SELEC * FROM user WHERE id in (3, 2, 1);和SELECT * FROM user WHERE id in (3, 2, 1) ORDER BY FIND_IN_SET(id, '3,2,1');。

//in查询
// QueryByIDs 根据给定ID查询
func QueryByIDs(ids []int) (users []user, err error) {
	// 动态填充id
	query, args, err := sqlx.In("SELECT id, name, age FROM user WHERE id IN (?)", ids)
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)

	err = db.Select(&users, query, args...)
	return
}

func ExecQueryByIDs() {
	users, err := QueryByIDs([]int{1, 6, 3, 4}) // 无法自定义顺序，默认id从小到大顺序
	if err != nil {
		fmt.Printf("QueryByIDs failed, err:%v\n", err)
		return
	}
	for _, user := range users {
		fmt.Printf("user:%#v\n", user)
	}
}

//in查询和FIND_IN_SET函数(给定顺序)
//查询id在给定id集合的数据并维持给定id集合的顺序。

// QueryAndOrderByIDs 按照指定id查询并维护顺序
func QueryAndOrderByIDs(ids []int) (users []user, err error) {
	// 动态填充id
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, fmt.Sprintf("%d", id))
	}
	// FIND_IN_SET维护顺序
	query, args, err := sqlx.In("SELECT id, name, age FROM user WHERE id IN (?) ORDER BY FIND_IN_SET(id, ?)", ids, strings.Join(strIDs, ","))
	if err != nil {
		return
	}

	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)

	err = db.Select(&users, query, args...)
	return
}

func ExecQueryAndOrderByIDs() {
	// 1. 用代码去做排序
	// 2. 让MySQL排序
	fmt.Println("----")
	users, err := QueryAndOrderByIDs([]int{1, 6, 3, 4}) // 维护id查询顺序
	if err != nil {
		fmt.Printf("QueryAndOrderByIDs failed, err:%v\n", err)
		return
	}
	for _, user := range users {
		fmt.Printf("user:%#v\n", user)
	}
}
