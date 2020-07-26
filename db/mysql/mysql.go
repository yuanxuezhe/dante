package mysql

import (
	"database/sql"
	"gitee.com/yuanxuezhe/dante/conf"
	"regexp"
	"strconv"
)

var mysqllist map[string]*sql.DB

func init() {
	mysqllist = make(map[string]*sql.DB, len(conf.Conf.Mysql))
	var err error
	var mysqlDb *sql.DB
	for _, v := range conf.Conf.Mysql {
		mysqlDb, err = sql.Open("mysql", v.Url)
		if err != nil {
			panic(err)
		}

		mysqlDb.SetMaxOpenConns(v.MaxOpenConns)
		mysqlDb.SetMaxIdleConns(v.MaxIdleConns)
		mysqllist[v.Rule] = mysqlDb
	}
}

func GetMysqlDB() *sql.DB {
	for k, v := range mysqllist {
		if VerifyRule(strconv.Itoa(123), k) {
			return v
		}
	}
	return nil
}

func VerifyRule(str string, rule string) bool {
	reg := regexp.MustCompile(rule)
	return reg.MatchString(str)
}
