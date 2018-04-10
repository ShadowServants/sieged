package main

import (
	"log"
	"os"
	"time"
	"fmt"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Print("Kek");
	var c int
	a := make([]int, 10)
	append(a, )
	t := time.NewTicker(500 * time.Millisecond);
	for {
		k := <- t.C
		log.Println("Heh ",k)
		fmt.Scan(&c);
	}
}
