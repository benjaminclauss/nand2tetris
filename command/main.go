package command

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"

	vm "github.com/benjaminclauss/nand2tetris/virtualmachine"
)

var vmTranslatorCommand = &cobra.Command{
	Use:  "vmtranslator [.vm file]",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input, err := os.Open(os.Args[1])
		if err != nil {
			return err
		}
		parser := vm.NewParser(input)
		output, err := os.Create("output.asm")
		if err != nil {
			return err
		}
		writer := vm.NewCodeWriter(output)
		for parser.HasMoreCommands() {
			parser.Advance()
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
		return nil
	},
}
