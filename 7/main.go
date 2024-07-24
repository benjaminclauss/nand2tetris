package main

import (
	"fmt"
	"os"
	"strconv"

	vm "github.com/benjaminclauss/nand2tetris/7/virtualmachine"
)

func main() {
	input, err := os.Open(os.Args[1])
	fmt.Println(os.Args[1])
	fmt.Println(err)
	parser := vm.NewParser(input)
	output, err := os.Create("output.asm")
	fmt.Println(err)
	writer := vm.NewCodeWriter(output)
	for parser.HasMoreCommands() {
		parser.Advance()
		fmt.Println("yo")
		fmt.Println(parser.CommandType())
		switch parser.CommandType() {
		case vm.CArithmetic:
			writer.WriteArithmetic(parser.Arg1())
		case vm.CPush:
			index, _ := strconv.Atoi(parser.Arg2())
			writer.WritePushPop(vm.CPush, parser.Arg1(), index)
		case vm.CPop:
			index, _ := strconv.Atoi(parser.Arg2())
			writer.WritePushPop(vm.CPop, parser.Arg1(), index)
		}
	}
}
