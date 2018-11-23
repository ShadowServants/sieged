package team

import (
	"sync"
)

type Data struct {
	Id    int
	Score Score
	sync.Mutex
}
