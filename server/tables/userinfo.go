package tables

import (
	. "dante/server/util/mysqlpool"
	"database/sql"
	"fmt"
	"log"
)

type Userinfo struct {
	Userid       string
	Username     string
	Passwd       string
	Sex          string
	Phone        int
	Email        string
	Status       string
	Registerdate int
}

func (t *Userinfo) QueryByKey() {
	conn, _ := Mysqlpool.Get()
	err := conn.(*sql.DB).QueryRow("SELECT * FROM userinfo where userid = ?", t.Userid).Scan(&t.Userid,
		&t.Username,
		&t.Passwd,
		&t.Sex,
		&t.Phone,
		&t.Email,
		&t.Status,
		&t.Registerdate)
	Mysqlpool.Put(conn)
	if err != nil {
		log.Fatal(err)
	}
	return
}
func (t *Userinfo) Query() {
	conn, _ := Mysqlpool.Get()
	rs, err := conn.(*sql.DB).Query("SELECT * FROM userinfo ")
	Mysqlpool.Put(conn)
	if err != nil {
		log.Fatalln(err)
	}
	//字段
	cols, _ := rs.Columns()
	for i := range cols {
		fmt.Print(cols[i])
		fmt.Print("\t")
	}
	fmt.Println("")
	fmt.Println("=================================")
	values := make([]sql.RawBytes, len(cols))
	scans := make([]interface{}, len(cols))

	for i := range values {

		scans[i] = &values[i]

	}

	results := make(map[int]map[string]string)

	i := 0

	for rs.Next() {

		if err := rs.Scan(scans...); err != nil {

			fmt.Println("Error")

			return

		}

		row := make(map[string]string)

		for j, v := range values {

			key := cols[j]

			row[key] = string(v)

		}

		results[i] = row

		i++

	}

	// 打印结果

	for i, m := range results {

		fmt.Println(i)

		for k, v := range m {

			fmt.Println(k, " : ", v)

		}

		fmt.Println("========================")

	}

	rs.Close()
}

func (t *Userinfo) Insert() {
	conn, _ := Mysqlpool.Get()
	rs, err := conn.(*sql.DB).Exec("INSERT INTO userinfo(userid,username,passwd,sex,phone,email,status,registerdate) VALUES ('1001', \"wanyuan\", 'ys6303618',\"1\", 18664324256, \"446968454@qq.com\", '0',20200628)")
	Mysqlpool.Put(conn)
	if err != nil {
		log.Fatalln(err)
	}
	rowCount, err := rs.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("inserted %d rows", rowCount)
}

// 校验用户是否存在
func (t *Userinfo) CheckUseridExist() bool {
	conn, _ := Mysqlpool.Get()
	err := conn.(*sql.DB).QueryRow("SELECT * FROM userinfo where userid = ?", t.Userid).Scan(&t.Userid,
		&t.Username,
		&t.Passwd,
		&t.Sex,
		&t.Phone,
		&t.Email,
		&t.Status,
		&t.Registerdate)
	Mysqlpool.Put(conn)
	if err != nil {
		log.Fatal(err)
	}
	return
}
