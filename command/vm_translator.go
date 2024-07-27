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
		// TODO: Validate input.
		// The VM translator should accept a single command line parameter, as follows:
		// prompt> VMtranslator source
		// Where source is either a file name of the form Xxx.vm (the extension is mandatory)
		// or a directory name containing one or more .vm files (in which case there is no extension).
		input, err := os.Open(args[0])
		if err != nil {
			return err
		}
		parser := vm.NewParser(input)
		// TODO: Name output appropriately.
		// The result of the translation is always a single assembly language file named Xxx.asm, created in the same directory as
		// the input Xxx. The translated code must conform to the standard VM mapping on the Hack platform.
		output, err := os.Create("output.asm")
		if err != nil {
			return err
		}
		// TODO: Handle multiple files.
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
