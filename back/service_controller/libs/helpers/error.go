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
