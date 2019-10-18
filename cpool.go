package cpool

import (
	"container/ring"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	SERV_REDIS     int = 1
	SERV_MYSQL     int = 2
	SERV_MEMCACHED int = 3
	SERV_HBASE     int = 4
	SERV_PG        int = 5
	SERV_MONGODB   int = 6
)

// Init server connection config
type ServOptions struct {
	Host                    string
	Port                    int
	ServType                int
	Username                string
	Password                string
	DbName                  string
	MaxConnNum, IdleConnNum int
	Heartbeat               int
}

type ServPool struct {
	serv IServ
	r    *ring.Ring
	m    *sync.RWMutex
}

// Init serv pool
func (pool *ServPool) init(options ServOptions) (err error) {
	// 静默
	_ = fmt.Printf
	var serv IServ
	switch options.ServType {
	case SERV_REDIS:
		serv = new(ServRedis)
	case SERV_MYSQL:
		serv = new(ServMysql)
	}
	if serv == nil {
		return
	}
	if err = serv.init(options); err != nil {
		return err
	}
	if err = pool.initPool(serv); err != nil {
		return err
	}
	pool.m, pool.serv = new(sync.RWMutex), serv
	return err
}

// Init connections to ring queue
func (pool *ServPool) initPool(serv IServ) (err error) {
	var (
		conn    interface{}
		options ServOptions = serv.options()
	)
	pool.r = ring.New(options.MaxConnNum)
	for i := 0; i < options.IdleConnNum; i++ {
		if conn, err = serv.connect(); err != nil {
			break
		} else {
			pool.r.Value = conn
			pool.r = pool.r.Next()
		}
	}
	return err
}

// Get connection from ring queue
// 通过迭代节点的方式取链接, 减少物理申请新节点的开销
func (pool *ServPool) get() (c interface{}, err error) {
	pool.m.Lock()
	defer pool.m.Unlock()
	options := pool.serv.options()
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(options.MaxConnNum)
	for i := 0; i < options.MaxConnNum; i++ {
		if i >= n && pool.r.Value != nil {
			c, err = pool.r.Value, nil
			pool.r.Value = nil
			return c, err
		}
		pool.r = pool.r.Next()
	}
	// 如果没有链接, 直接建立链接
	return pool.serv.connect()
}

// Try release connection to ring queue
// 放回任意一个空的槽
func (pool *ServPool) release(c interface{}) {
	pool.m.Lock()
	pool.m.Unlock()
	options := pool.serv.options()
	for i := 0; i < options.MaxConnNum; i++ {
		// 放入空的节点
		if pool.r.Value == nil {
			pool.r.Value = c
			return
		}
		pool.r = pool.r.Next()
	}
	// 若无空结点则销毁
	pool.serv.destroy(c)
}
