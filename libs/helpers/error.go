package helpers

import (
	"log"
	"fmt"
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


func FromBytesToString(buf []byte, index int) string {
	return string(buf[:index])
}
