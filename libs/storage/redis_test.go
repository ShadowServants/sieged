package storage

import (
	"github.com/garyburd/redigo/redis"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"hackforces/libs/helpers"
)

func GetSimpleRedisExecutor() *RedisPoolExecutor{
	conn,err := redis.Dial("tcp",":6379")
	helpers.FailOnError(err,"Redis is down")
	pool := helpers.NewPool(conn,1)
	pool_exec := RedisPoolExecutor{pool,sync.Mutex{}}
	return &pool_exec

}


func TestRedisPoolExecutor_Exec(t *testing.T) {
	executor := GetSimpleRedisExecutor()
	Convey("Test that redis executor works",t,func(){
		res,err := executor.Exec("SET","key","value")
		So(err,ShouldEqual,nil)
		So(res,ShouldEqual,"OK")
		res,err = executor.Exec("GET","key")
		So(err,ShouldEqual,nil)
		result, ok := res.([]byte)
		So(ok,ShouldEqual,true)

		So(string(result),ShouldEqual,"value")
	})
}



func TestBaseRedisStorage_Get(t *testing.T) {
	executor := GetSimpleRedisExecutor()
	storage := BaseRedisStorage{executor}
	Convey("Check base redis storage",t,func(){
		storage.Set("SET","key","value")
		res, err := storage.Get("GET","key")
		So(err,ShouldEqual,nil)
		So(res,ShouldEqual,"value")

	})
}

func TestSimpleRedisStorage_Get(t *testing.T) {
	executor := GetSimpleRedisExecutor()
	storage := SimpleRedisStorage{BaseRedisStorage{executor}}
	Convey("Check base redis storage",t,func(){
		storage.Set("key","value")
		res, err := storage.Get("key")
		So(err,ShouldEqual,nil)
		So(res,ShouldEqual,"value")

	})
}

func TestHsetRedisStorage_Get(t *testing.T) {
	executor := GetSimpleRedisExecutor()
	storage := HsetRedisStorage{BaseRedisStorage{executor},"TESTSET"}

	Convey("Check hset redis storage",t,func(){
		storage.Set("key","value")
		res, err := storage.Get("key")
		So(err,ShouldEqual,nil)
		So(res,ShouldEqual,"value")

	})
}


func TestRedisKeySet_Add(t *testing.T) {
	executor := GetSimpleRedisExecutor()
	keyset := RedisKeySet{executor}
	Convey("Check redis keyset works",t,func(){
		keyset.Add("first_key","1")
		keyset.Add("first_key","2")
	})
	Convey("Check add was good",t,func(){
		exist := keyset.Check("first_key","2")
		exist1 := keyset.Check("first_key","1")

		So(exist,ShouldEqual,true)
		So(exist1,ShouldEqual,true)

	})
	Convey("Check bad data",t,func(){
		exist := keyset.Check("first","1")
		So(exist,ShouldEqual,false)
		exist = keyset.Check("first_key","3")
		So(exist,ShouldEqual,false)
	})

}