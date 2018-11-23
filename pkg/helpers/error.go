package helpers

import (
	"fmt"
	"log"
)

func FailOnError(err error, msg string) {
        if err != nil {
			message := fmt.Sprintf("%s: %s", msg, err)
			log.Fatal(message)
			panic(message)
        }
}

