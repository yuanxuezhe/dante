package pool

import (
	"database/sql"
	"gitee.com/yuanxuezhe/dante/core/conf"
	"gitee.com/yuanxuezhe/ynet/yconnpool"
	"time"
)

var Mysqlpool *yconnpool.ConnPool

func init() {
	var err error
	//fmt.Println(conf.Conf.Mysql.Url)
	Mysqlpool, err = yconnpool.NewConnPool(func() (yconnpool.ConnRes, error) {
		return sql.Open("mysql", conf.Conf.Mysql.Url)
	}, conf.Conf.Mysql.Maxcount, time.Second*100)
	if err != nil {
		panic(err)
	}
}
