
# README

## demo
```
package main
import (
    "cpool"
    "fmt"
)

func main() {
    mannger := &cpool.CPoolMannger{}
    options := cpool.ServOptions{
        Host:        `127.0.0.1`,
        Port:        3306,
        MaxConnNum:  8,  
        IdleConnNum: 3,
        Heartbeat:   3,  
        ServType:    cpool.SERV_MYSQL,
    }   
    mannger.Register(`mysql`, options)
    options = cpool.ServOptions{
        Host:        `127.0.0.1`,
        Port:        6379,
        MaxConnNum:  8,  
        IdleConnNum: 3,
        Heartbeat:   3,  
        ServType:    cpool.SERV_REDIS,
    }   
    mannger.Register(`redis`, options)
    var xx chan int 
    <-xx
}
```

