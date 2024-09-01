package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	vm "github.com/benjaminclauss/nand2tetris/virtualmachine"
)

func NewVMTranslatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "vmtranslator <source>",
		Long: `
The VM translator accepts a single command line parameter, as follows:

prompt> VMtranslator source

Where source is either a file name of the form Xxx.vm (the extension is mandatory)
or a directory name containing one or more .vm files (in which case there is no extension).

The result of the translation is always a single assembly language file named Xxx.asm, 
created in the same directory as the input Xxx. 
	`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := os.Stat(args[0])
			if err != nil {
				return fmt.Errorf("error getting FileInfo: %w", err)
			}
			if info.IsDir() {
				entries, err := os.ReadDir(args[0])
				if err != nil {
					return err
				}
				var vmFiles []string
				for _, entry := range entries {
					if strings.HasSuffix(entry.Name(), ".vm") {
						vmFiles = append(vmFiles, filepath.Join(args[0]+"/", entry.Name()))
					}
				}
				return translate(info.Name()+".asm", vmFiles...)
			}
			return translate(strings.TrimRight(filepath.Base(args[0]), ".vm")+".asm", args[0])
		},
	}

	return cmd
}

func translate(outputFilename string, files ...string) error {
	output, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	writer := vm.NewCodeWriter(output)

	init := false
	for _, filename := range files {
		if strings.HasSuffix(filename, "Sys.vm") {
			init = true
		}
	}
	if init {
		writer.WriteInit()
	}

	for _, file := range files {
		writer.SetFilename(file)
		vmf, err := os.Open(file)
		if err != nil {
			return err
		}
		parser := vm.NewParser(vmf)
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
			case vm.CLabel:
				writer.WriteLabel(parser.Arg1())
			case vm.CIf:
				writer.WriteIf(parser.Arg1())
			case vm.CGoTo:
				writer.WriteGoto(parser.Arg1())
			case vm.CFuntion:
				numLocals, _ := strconv.Atoi(parser.Arg2())
				writer.WriteFunction(parser.Arg1(), numLocals)
			case vm.CReturn:
				writer.WriteReturn()
			case vm.CCall:
				nArgs, _ := strconv.Atoi(parser.Arg2())
				writer.WriteCall(parser.Arg1(), nArgs)
			}
		}
	}
	return nil
}
