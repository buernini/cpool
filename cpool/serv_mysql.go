package cpool

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type ServMysql struct {
	etc ServOptions
}

// Init
func (msq *ServMysql) init(options ServOptions) (err error) {
	msq.etc = options
	return err
}

// 返回配置
func (msq *ServMysql) options() ServOptions {
	return msq.etc
}

// 建立一个链接
func (msq *ServMysql) connect() (conn interface{}, err error) {
	var address string = fmt.Sprintf("%s/%s@tcp(%s:%d)/%s",
		msq.etc.Username,
		msq.etc.Password,
		msq.etc.Host,
		msq.etc.Port,
		msq.etc.DbName)
	return sql.Open(`mysql`, address)
}

// 保活
func (msq *ServMysql) heartbeat(c interface{}) error {
	time.Sleep(time.Millisecond * time.Duration(1))
	var value string
	if rows, err := c.(*sql.DB).Query(`SELECT version() AS ver`); err != nil {
		return err
	} else {
		for rows.Next() {
			if err = rows.Scan(&value); err != nil {
				return err
			}
		}
	}
	if len(value) > 0 {
		return nil
	}
	return errors.New(`heartbeat call response null.`)
}

// 销毁链接
func (msq *ServMysql) destroy(c interface{}) error {
	return c.(*sql.DB).Close()
}
