package cpool

import (
	"errors"
	"fmt"
	"time"
)

// CpoolMannger
type CPoolMannger struct {
	servMap map[string]*ServPool
}

// 注册服务
func (cpm *CPoolMannger) Register(servId string, options ServOptions) (err error) {
	if cpm.servMap == nil {
		cpm.servMap = make(map[string]*ServPool)
	}
	var pool *ServPool
	if _, ok := cpm.servMap[servId]; !ok {
		pool = new(ServPool)
		if err = pool.init(options); err == nil {
			go cpm.heartbeat(pool)
			cpm.servMap[servId] = pool
		} else {
			// Write error log
		}
	}
	return err
}

// 租借链接
func (cpm *CPoolMannger) Get(servId string) (c interface{}, err error) {
	var pool *ServPool
	var ok bool
	if pool, ok = cpm.servMap[servId]; !ok {
		return c, errors.New(`CPoolMannger get connection failed. not found service by servId`)
	}
	return pool.get()
}

// 归还链接
func (cpm *CPoolMannger) Release(servId string, c interface{}) {
	if pool, ok := cpm.servMap[servId]; ok {
		pool.release(c)
	}
}

// 保活
func (cpm *CPoolMannger) heartbeat(servPool *ServPool) {
	options := servPool.serv.options()
	ticker := time.NewTicker(time.Second * time.Duration(options.Heartbeat))
	for range ticker.C {
		servPool.r.Do(func(c interface{}) {
			if c != nil {
				if err := servPool.serv.heartbeat(c); err != nil {
					fmt.Printf("CPoolMannger heartbeat service happend error. message:%s\n", err)
				} else {
					//fmt.Printf("CPoolMannger heartbeat service is alived info:%s:%d\n", options.Host, options.Port)
				}
			}
		})
	}
}
