package storage

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jnovikov/hackforces/back/libs/helpers"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func GetSimpleRedisExecutor() *RedisPoolExecutor{
	conn,_ := redis.Dial("tcp",":6379")
	pool := helpers.NewPool(conn,1)
	pool_exec := RedisPoolExecutor{pool}
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