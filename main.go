package main

import (
	"fmt"
	"os"

	"github.com/benjaminclauss/nand2tetris/command"
)

func main() {
	if err := command.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
