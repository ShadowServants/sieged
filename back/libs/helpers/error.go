package helpers

import (
	"log"
	"fmt"
	"bytes"
)

func FailOnError(err error, msg string) {
        if err != nil {
			message := fmt.Sprintf("%s: %s", msg, err)
			log.Fatal(message)
			panic(message)
        }
}

func MinInt(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func MaxInt(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0 // return correctly abs(-0)
	}
	return x
}


func FromBytesToString(buf []byte) string {
	n := bytes.IndexByte(buf, 0)
	e := bytes.IndexByte(buf,10)
	if (n - e == 1){
		return string(buf[:e])
	}
	return string(buf[:n])
}
