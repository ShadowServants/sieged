package storage

import (
	"github.com/mediocregopher/radix.v2/pool"
	"fmt"
	"github.com/jnovikov/hackforces/libs/helpers"
)



type RadixFactory struct {
	Pool *RadixPool
}

func (rf *RadixFactory) GetSimpleStorage() *SimpleRadixStorage {
	return &SimpleRadixStorage{rf.Pool}
}

func (rf *RadixFactory) GetHsetStorage(setname string) *HsetRadixStorage {
	return &HsetRadixStorage{rf.Pool,setname}
}


func (rf *RadixFactory) GetKeySet() *RadixKeySet {
	return &RadixKeySet{rf.Pool}
}

type RadixPool struct {
	pool *pool.Pool
}

func (rp *RadixPool) Build (host string, port string,size int) {
	p, err := pool.New("tcp", fmt.Sprintf("%s:%s",host,port), size)
	if err != nil {
		helpers.FailOnError(err,"Redis is down")
	}

	rp.pool = p
}

type SimpleRadixStorage struct {
	RadixP *RadixPool
}

func (bs *SimpleRadixStorage) Get(key string) (string,error) {
	conn, err := bs.RadixP.pool.Get()
	defer bs.RadixP.pool.Put(conn)
	if err != nil {
		return "",err
	}
	data ,err := conn.Cmd("GET",key).Str()
	return data,err
}

func (bs *SimpleRadixStorage) Set(key string,value string) () {
	conn, err := bs.RadixP.pool.Get()
	defer bs.RadixP.pool.Put(conn)
	if err != nil {
		return
	}
	conn.Cmd("SET",key,value)
}


type HsetRadixStorage struct {
	RadixP *RadixPool
	SetName string
}

func (hr *HsetRadixStorage) Get(key string) (string,error) {
	conn, err := hr.RadixP.pool.Get()
	defer hr.RadixP.pool.Put(conn)
	if err != nil {
		return "",err
	}
	data ,err := conn.Cmd("HGET",hr.SetName,key).Str()
	return data,err
}

func (hr *HsetRadixStorage) Set(key string,value string) () {
	conn, err := hr.RadixP.pool.Get()
	defer hr.RadixP.pool.Put(conn)
	if err != nil {
		return
	}
	conn.Cmd("HSET",hr.SetName,key,value)
}

type RadixKeySet struct {
	RadixP	*RadixPool
}

func (ks *RadixKeySet) Add(key string, value string) {
	conn, err := ks.RadixP.pool.Get()
	defer ks.RadixP.pool.Put(conn)
	if err != nil {
		return
	}
	conn.Cmd("SADD",key,value)
}

func (ks *RadixKeySet) Check(key string, value string) bool {
	conn, err := ks.RadixP.pool.Get()
	defer ks.RadixP.pool.Put(conn)
	if err != nil {
		return false
	}
	exist, err := conn.Cmd("SISMEMBER",key,value).Int()
	if err != nil {
		//Panic ?
		return false
	}
	return exist != 0

}
