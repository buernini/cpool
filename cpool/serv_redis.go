package cpool

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type ServRedis struct {
	etc ServOptions
}

// Init
func (rds *ServRedis) init(options ServOptions) error {
	rds.etc = options
	return nil
}

// 返回配置
func (rds *ServRedis) options() ServOptions {
	return rds.etc
}

// 建立一个链接
func (rds *ServRedis) connect() (conn interface{}, err error) {
	var addr string = fmt.Sprintf("%s:%d", rds.etc.Host, rds.etc.Port)
	if conn, err = redis.Dial("tcp", addr); err != nil {
		return conn, err
	}
	return conn, err
}

// 保活
func (rds *ServRedis) heartbeat(c interface{}) error {
	time.Sleep(time.Millisecond * time.Duration(1))
	if data, err := c.(redis.Conn).Do(`PING`); err != nil {
		return err
	} else if data.(string) != `PONG` {
		return errors.New(`heartbeat call response is not PONG.`)
	} else {
		return nil
	}
}

// 销毁链接
func (rds *ServRedis) destroy(c interface{}) error {
	return c.(redis.Conn).Close()
}
