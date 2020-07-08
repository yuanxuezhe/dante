package tables

import (
	. "dante/server/util/pool"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"runtime"
)

type Userinfo struct {
	Userid       int
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
	//if err != nil {
	//	fmt.Println("err5")
	//	log.Fatal(err)
	//}
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		// 使用该方式可以打印出运行时的错误信息, 该种错误是编译时无法确定的
		if _, file, line, ok := runtime.Caller(0); ok {
			fmt.Println(err, file, line)
		}
	}
	return
}
func (t *Userinfo) Insert() {
	conn, _ := Mysqlpool.Get()
	rs, err := conn.(*sql.DB).Exec("INSERT INTO userinfo(userid,username,passwd,sex,phone,email,status,registerdate) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		t.Userid, t.Username, t.Passwd, t.Sex, t.Phone, t.Email, t.Status, t.Registerdate)

	Mysqlpool.Put(conn)
	if err != nil {
		log.Println("err1")
		log.Fatalln(err)
	}
	rowCount, err := rs.RowsAffected()
	if err != nil {
		log.Println("err2")
		log.Fatalln(err)
	}
	log.Printf("inserted %d rows", rowCount)
}

//// 校验用户是否存在,若存在返回用户信息
//func (t *Userinfo) CheckAccountExist() (*Userinfo, error) {
//	conn, _ := Mysqlpool.Get()
//	err := conn.(*sql.DB).QueryRow("SELECT * FROM userinfo where (userid = ? or phone = ? or email = ?) and passwd = ?",
//		t.Userid,
//		t.Phone,
//		t.Email,
//		t.Passwd).Scan(&t.Userid,
//		&t.Username,
//		&t.Passwd,
//		&t.Sex,
//		&t.Phone,
//		&t.Email,
//		&t.Status,
//		&t.Registerdate)
//	Mysqlpool.Put(conn)
//
//	switch {
//	case err == sql.ErrNoRows:
//		return nil, errors.New("Login failed! Userinfo not exists!")
//	case err != nil:
//		// 使用该方式可以打印出运行时的错误信息, 该种错误是编译时无法确定的
//		if _, file, line, ok := runtime.Caller(0); ok {
//			fmt.Println(err, file, line)
//		}
//		return nil, err
//	}
//	return t, nil
//}

// 校验用户是否存在,若存在返回用户信息
func (t *Userinfo) CheckAccountExist() (userinfo *Userinfo, err error) {

	conn, _ := Mysqlpool.Get()
	stmt, err := conn.(*sql.DB).Prepare("SELECT * FROM userinfo where (userid = ? or phone = ? or email = ?) and passwd = ?")
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(t.Userid, t.Phone, t.Email, t.Passwd)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	Mysqlpool.Put(conn)

	if rows.Next() {
		err = rows.Scan(&t.Userid,
			&t.Username,
			&t.Passwd,
			&t.Sex,
			&t.Phone,
			&t.Email,
			&t.Status,
			&t.Registerdate)

		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Login failed! Userinfo not exists or passwd is wrong!")
	}

	return t, nil
}

func (t *Userinfo) CheckAvailable_Phone() error {
	var err error
	conn, err := Mysqlpool.Get()
	if err != nil {
		return err
	}
	rows, err := conn.(*sql.DB).Query("SELECT * FROM userinfo where phone = ?", t.Phone)
	if err != nil {
		return err
	}
	Mysqlpool.Put(conn)

	if rows.Next() {
		return errors.New("phone num has been used!")
	}

	return nil
}
