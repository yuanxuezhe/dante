package tables

import (
	"database/sql"
	"errors"
	"fmt"
	. "gitee.com/yuanxuezhe/dante/server/util/pool"
	"log"
	"runtime"
)

//`goodsid`   BIGINT(20)  default 0       PRIMARY KEY COMMENT '编号',
//`goodsname`       VARCHAR(128) NOT NULL COMMENT '名称',
//`type`      int default 0               COMMENT '商品类型',
//`source`    int default 0               COMMENT '来源',
//`url`     VARCHAR(128) NOT NULL          COMMENT '链接',
//`imgurl`    VARCHAR(128) NOT NULL        COMMENT '图片链接',
//`brand`     int default 0               COMMENT '品牌',
//`status`    int default 0     		  COMMENT '状态',
//`date`      int default 0               COMMENT '日期',
//`time`      int default 0               COMMENT '时间'

type Goods struct {
	Goodsid   int32  `json:"goodsid"`   //编号
	Goodsname string `json:"goodsname"` //名称
	Type      int    `json:"type"`      //商品类型
	Source    string `json:"source"`    //来源
	Url       string `json:"url"`       //链接
	Imgurl    string `json:"imgurl"`    //图片链接
	Brand     int    `json:"brand"`     //品牌
	Status    int    `json:"status"`    //状态
	Date      int    `json:"date"`      //日期
	Time      int    `json:"time"`      //时间
}

func (t *Goods) QueryByKey() error {
	conn, err := Mysqlpool.Get()
	if err != nil {
		return err
	}
	err = conn.(*sql.DB).QueryRow("SELECT * FROM goods where goodsid = ?", t.Goodsid).Scan(&t.Goodsid,
		&t.Goodsname,
		&t.Type,
		&t.Source,
		&t.Url,
		&t.Imgurl,
		&t.Brand,
		&t.Status,
		&t.Date,
		&t.Time)
	Mysqlpool.Put(conn)
	if err != nil {
		return errors.New("Get record from mysql failed!")
	}
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		// 使用该方式可以打印出运行时的错误信息, 该种错误是编译时无法确定的
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	return nil
}

func (t *Goods) Insert() error {
	conn, err := Mysqlpool.Get()
	if err != nil {
		return err
	}
	rs, err := conn.(*sql.DB).Exec("INSERT INTO userinfo(goodsid,goodsname,type,source,url,imgurl,brand,status,date,time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		t.Goodsid, t.Goodsname, t.Type, t.Source, t.Url, t.Imgurl, t.Brand, t.Status, t.Date, t.Time)

	Mysqlpool.Put(conn)
	if err != nil {
		return err
	}
	rowCount, err := rs.RowsAffected()
	if err != nil {
		return err
	}
	log.Printf("inserted %d rows", rowCount)
	return nil
}

func (t *Goods) QueryByStatus() ([]Goods, error) {
	goodlist := []Goods{}
	conn, err := Mysqlpool.Get()
	if err != nil {
		return nil, err
	}
	stmt, err := conn.(*sql.DB).Prepare("SELECT * FROM goods where status = ?")
	if err != nil {
		return nil, errors.New("Connection to mysql failed!")
	}
	defer stmt.Close()
	rows, err := stmt.Query(t.Status)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	Mysqlpool.Put(conn)

	for rows.Next() {
		m := Goods{}
		err = rows.Scan(&m.Goodsid,
			&m.Goodsname,
			&m.Type,
			&m.Source,
			&m.Url,
			&m.Imgurl,
			&m.Brand,
			&m.Status,
			&m.Date,
			&m.Time)

		if err != nil {
			return nil, err
		}
		goodlist = append(goodlist, m)
	}
	fmt.Println(goodlist)
	return goodlist, nil
}
