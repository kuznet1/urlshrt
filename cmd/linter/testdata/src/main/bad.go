package main

import (
	"log"
	"os"
)

func foo() {
	log.Fatal("foo") // want "should not use log.Fatal"
	os.Exit(666)     // want "should not use os.Exit"
	panic("bar")     // want "should not use panic"
}
