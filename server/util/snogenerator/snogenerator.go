package snogenerator

import (
	"dante/server/util/pool"
	"database/sql"
	"log"
	"sync"
)

var useridRW sync.RWMutex

func NewUserid() int {
	var userid int
	conn, err := pool.Mysqlpool.Get()
	useridRW.Lock()
	err = conn.(*sql.DB).QueryRow("SELECT max(userid) + 1  FROM userinfo").Scan(&userid)
	useridRW.Unlock()
	pool.Mysqlpool.Put(conn)
	if err != nil {
		log.Fatal(err)
	}
	return userid
}
