package main

import (
	"github.com/benjaminclauss/nand2tetris/command"
	"log"
)

func main() {
	if err := command.RootCmd.Execute(); err != nil {
		log.Fatalf(err.Error())
	}
}
