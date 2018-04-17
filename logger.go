package main

import (
	"fmt"
	"os"
)

func print(params ...interface{}) {
	if os.Getenv("ENV") != "prod" {
		fmt.Println(params)
	}
}

func err(e error) bool {
	if e != nil {
		if os.Getenv("ENV") != "prod" {
			fmt.Println(e.Error())
		}
		return true
	}
	return false
}

func crash(e error) {
	if e != nil {
		panic(e)
	}
}
