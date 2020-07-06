package pool

import (
	"dante/core/conf"
	"database/sql"
	"gitee.com/yuanxuezhe/ynet/yconnpool"
	"time"
)

var Mysqlpool *yconnpool.ConnPool

func init() {
	Mysqlpool, _ = yconnpool.NewConnPool(func() (yconnpool.ConnRes, error) {
		return sql.Open("mysql", conf.Conf.Mysql.Url)
	}, conf.Conf.Mysql.Maxcount, time.Second*100)
}
