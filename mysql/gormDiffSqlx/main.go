package tesst

import (
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ShardingDB struct {
	ID        uint64    `db:"id" gorm:"column:id"`
	DBID      string    `db:"db_id" gorm:"column:db_id"`
	Host      string    `db:"host" gorm:"column:host"`
	Port      int32     `db:"port" gorm:"column:port"`
	User      string    `db:"user" gorm:"column:user"`
	Password  string    `db:"password" gorm:"column:password"`
	Memo      string    `db:"memo" gorm:"column:memo"`
	State     int32     `db:"state" gorm:"column:state"`
	CreatedAt time.Time `db:"create_time" gorm:"column:create_time"`
	UpdatedAt time.Time `db:"update_time" gorm:"column:update_time"`
}

// TableName returns table name of sharding_db.
func (s *ShardingDB) TableName() string {
	return "sharding_db"
}

var (
	dbhostsip  = "127.0.0.1" // IP地址
	dbport     = 3306        // Port
	dbusername = "root"      // 用户名
	dbpassword = "admin"     // 密码
	dbname     = "test"      // 表名
)

func Benchmark(b *testing.B) {
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=UTC",
		dbusername, dbpassword, dbhostsip, dbport, dbname)

	limits := []int{
		5,
		50,
		500,
		10000,
	}

	sqlxDB, _ := sqlx.Connect("mysql", dsn)
	sqlxDB.SetMaxOpenConns(500)
	sqlxDB.SetMaxIdleConns(100)

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NowFunc: func() time.Time { return time.Now().UTC().Round(time.Microsecond) }})
	db, _ := gormDB.DB()
	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(100)

	fmt.Println("=============================== CPU:8 MEM:16G MaxOpenConns:500 MaxOpenConns:100 ====================================")

	for _, lim := range limits {
		lim := lim

		// Benchmark sqlx
		b.Run(fmt.Sprintf("sqlx limit:%d", lim), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				q := fmt.Sprintf("SELECT * FROM sharding_db ORDER BY id LIMIT %d", lim)
				res := []ShardingDB{}
				err := sqlxDB.Select(&res, q)
				if err != nil {
					b.Fatal(err)
				}
			}
		})

		if err != nil {
			panic(err)
		}
		// Benchmark gormDB
		b.Run(fmt.Sprintf("gormDB limit:%d", lim), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for i := 0; i < b.N; i++ {
					var res = []ShardingDB{}
					err := gormDB.Order("id").Limit(lim).Find(&res).Error

					if err != nil {
						b.Fatal(err)
					}
				}
			}
		})

		fmt.Println("==================================================================================================================")
	}
}
