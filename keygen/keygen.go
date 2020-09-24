package main

import (
	"fmt"
	"time"
)

func main() {
	dtNow := time.Now().Format("2006-01-02")
	dtKeyExpaire := "2019-01-01"

	fmt.Println(dtNow)
	fmt.Println(dtKeyExpaire)
	if dtNow < dtKeyExpaire {
		fmt.Println("ok")
	} else {
		fmt.Println("bad key")
	}
}
