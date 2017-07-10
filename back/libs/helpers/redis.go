package helpers

import (
	"github.com/garyburd/redigo/redis"
)

func NewPool(conn redis.Conn,maxIdle int) *redis.Pool {
  return &redis.Pool{
    MaxIdle: 1000,
	MaxActive:0,
  	Wait:true,
    Dial: func () (redis.Conn, error) { return conn,nil },
  }
}

//func t() {
//	c,_ := redis.Dial("asd","lol")
//}
