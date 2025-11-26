package main

import (
	"log"
	"os"
)

func main() {
	log.Fatal("foo")
	os.Exit(666)
}
