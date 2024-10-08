package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/benjaminclauss/nand2tetris/assembler"
)

var assemblerCmd = &cobra.Command{
	Use:   "assembler [.asm file]",
	Short: "Assembler for Hack programs",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFilename := args[0]
		inputFile, err := os.Open(inputFilename)
		if err != nil {
			return err
		}
		defer inputFile.Close()
		outputFilename := fmt.Sprintf("%s.hack", strings.TrimSuffix(inputFilename, ".asm"))
		outputFile, err := os.Create(outputFilename)
		if err != nil {
			return err
		}
		defer outputFile.Close()
		err = Assemble(inputFile, outputFile)
		if err != nil {
			return err
		}
		return outputFile.Sync()
	},
}

func Assemble(input io.Reader, output io.Writer) error {
	var buf bytes.Buffer
	tee := io.TeeReader(input, &buf)

	st := assembler.NewSymbolTable()
	initializeSymbolTable(st)
	firstPassParser := assembler.NewParser(tee)
	currentROMAddress := 0
	for firstPassParser.HasMoreCommands() {
		switch firstPassParser.CommandType() {
		case assembler.L_COMMAND:
			symbol := firstPassParser.Symbol()
			st.AddEntry(symbol, currentROMAddress)
		case assembler.C_COMMAND, assembler.A_COMMAND:
			currentROMAddress++
		}
		firstPassParser.Advance()
	}

	secondPassParser := assembler.NewParser(&buf)
	nextAvailableRAMAddress := 16

	for secondPassParser.HasMoreCommands() {
		switch secondPassParser.CommandType() {
		case assembler.A_COMMAND:
			symbol := secondPassParser.Symbol()
			if number, err := strconv.Atoi(symbol); err == nil {
				io.WriteString(output, fmt.Sprintf("0%015b\n", number))
			} else {
				if st.Contains(symbol) {
					address := st.GetAddress(symbol)
					io.WriteString(output, fmt.Sprintf("0%015b\n", address))
				} else {
					st.AddEntry(symbol, nextAvailableRAMAddress)
					io.WriteString(output, fmt.Sprintf("0%015b\n", nextAvailableRAMAddress))
					nextAvailableRAMAddress++
				}
			}
		case assembler.C_COMMAND:
			comp := assembler.Comp(secondPassParser.Comp())
			dest := assembler.Dest(secondPassParser.Dest())
			jump := assembler.Jump(secondPassParser.Jump())
			io.WriteString(output, fmt.Sprintf("111%s%s%s\n", comp, dest, jump))
		}
		secondPassParser.Advance()
	}
	return nil
}

// Initialize the symbol table with all the predefined symbols and their pre-allocated RAM addresses.
func initializeSymbolTable(st *assembler.SymbolTable) {
	predefinedSymbols := map[string]int{
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
		"R0":     0,
		"R1":     1,
		"R2":     2,
		"R3":     3,
		"R4":     4,
		"R5":     5,
		"R6":     6,
		"R7":     7,
		"R8":     8,
		"R9":     9,
		"R10":    10,
		"R11":    11,
		"R12":    12,
		"R13":    13,
		"R14":    14,
		"R15":    15,
		"SCREEN": 16384,
		"KBD":    24576,
	}
	for symbol, address := range predefinedSymbols {
		st.AddEntry(symbol, address)
	}
}
