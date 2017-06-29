package helpers

import (
	"time"
	"github.com/garyburd/redigo/redis"
)

func NewPool(conn redis.Conn,maxIdle int) *redis.Pool {
  return &redis.Pool{
    MaxIdle: maxIdle,
    IdleTimeout: 240 * time.Second,
    Dial: func () (redis.Conn, error) { return conn,nil },
  }
}

//func t() {
//	c,_ := redis.Dial("asd","lol")
//}
