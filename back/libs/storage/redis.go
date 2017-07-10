package storage

import (
	"github.com/garyburd/redigo/redis"
	"errors"
	"github.com/jnovikov/hackforces/back/libs/helpers"
	"fmt"
	"sync"
)


func GetRedisExecutor(port string,pool_size int) *RedisPoolExecutor{
	conn,err := redis.Dial("tcp",":"+port)
	helpers.FailOnError(err,"Redis is down")
	pool := helpers.NewPool(conn,pool_size)
	pool_exec := RedisPoolExecutor{pool,sync.Mutex{}}
	return &pool_exec

}

type RedisPoolExecutor struct {
	Pool *redis.Pool
	mu sync.Mutex
}


func (re *RedisPoolExecutor) Exec(command string,args...string) (interface {},error) {
	re.mu.Lock()
	defer re.mu.Unlock()
	real_args := make([]interface{},len(args))
	for k,v := range args {
		real_args[k] = v
	}
	red := re.Pool.Get()

	defer red.Close()
	return red.Do(command,real_args...)
}


func (re *RedisPoolExecutor) Close() {
	re.Pool.Close()
}



type BaseRedisStorage struct {
	Redis *RedisPoolExecutor
}

func (rs *BaseRedisStorage) Get(command string,args...string) (string,error) {
	data , err := rs.Redis.Exec(command,args...)
	if err != nil {
		return "",err
	}
	if res, ok := data.([]byte); ok {
		return string(res),nil
	}
	return "", errors.New("Not string")
}


func (rs *BaseRedisStorage) Set(command string,args...string) () {
	_, err := rs.Redis.Exec(command,args...)
	if err != nil {
		fmt.Println("BAD",command,args, err.Error())
		helpers.FailOnError(err,"Redis error") //?
		//Write to logfile or panic ?
		return
	}
}

type SimpleRedisStorage struct {
	BaseRedisStorage
}

func (sr *SimpleRedisStorage) Get(key string) (string,error){
	return sr.BaseRedisStorage.Get("GET",key)
}

func (sr *SimpleRedisStorage) Set(key string,value string) {
	sr.BaseRedisStorage.Set("SET",key,value)
}

type HsetRedisStorage struct {
	BaseRedisStorage
	SetName string
}

func (rs *HsetRedisStorage) Get(key string) (string,error) {
	return rs.BaseRedisStorage.Get("HGET",rs.SetName,key)
}


func (rs *HsetRedisStorage) Set(key string,value string) () {
	rs.BaseRedisStorage.Set("HSET",rs.SetName,key,value)
}

type RedisKeySet struct {
	Redis   *RedisPoolExecutor
}

func (ks *RedisKeySet) Add(key string, value string) {
	ks.Redis.Exec("SADD",key,value)
}

func (ks *RedisKeySet) Check(key string, value string) bool {
	exist, err := ks.Redis.Exec("SISMEMBER",key,value)
	if err != nil {
		//Panic ?
		return false
	}

	if res, ok := exist.(int64); ok {
		return res != 0
	}
	return false

}