package mysqlpool

import (
	"database/sql"
	"gitee.com/yuanxuezhe/ynet/yconnpool"
	"time"
)

var Mysqlpool *yconnpool.ConnPool

func init() {
	Mysqlpool, _ = yconnpool.NewConnPool(func() (yconnpool.ConnRes, error) {
		return sql.Open("mysql", "root:1@tcp(192.168.3.25:3306)/dante?parseTime=true")
	}, 100, time.Second*100)
}
