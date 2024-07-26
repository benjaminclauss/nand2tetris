package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
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

	st := NewSymbolTable()
	initializeSymbolTable(st)
	firstPassParser := NewParser(tee)
	currentROMAddress := 0
	for firstPassParser.hasMoreCommands() {
		switch firstPassParser.commandType() {
		case L_COMMAND:
			symbol := firstPassParser.symbol()
			st.addEntry(symbol, currentROMAddress)
		case C_COMMAND, A_COMMAND:
			currentROMAddress++
		}
		firstPassParser.advance()
	}

	secondPassParser := NewParser(&buf)
	nextAvailableRAMAddress := 16

	for secondPassParser.hasMoreCommands() {
		switch secondPassParser.commandType() {
		case A_COMMAND:
			symbol := secondPassParser.symbol()
			if number, err := strconv.Atoi(symbol); err == nil {
				io.WriteString(output, fmt.Sprintf("0%015b\n", number))
			} else {
				if st.contains(symbol) {
					address := st.GetAddress(symbol)
					io.WriteString(output, fmt.Sprintf("0%015b\n", address))
				} else {
					st.addEntry(symbol, nextAvailableRAMAddress)
					io.WriteString(output, fmt.Sprintf("0%015b\n", nextAvailableRAMAddress))
					nextAvailableRAMAddress++
				}
			}
		case C_COMMAND:
			comp := comp(secondPassParser.comp())
			dest := dest(secondPassParser.dest())
			jump := jump(secondPassParser.jump())
			io.WriteString(output, fmt.Sprintf("111%s%s%s\n", comp, dest, jump))
		}
		secondPassParser.advance()
	}
	return nil
}

// Initialize the symbol table with all the predefined symbols and their pre-allocated RAM addresses.
func initializeSymbolTable(st *SymbolTable) {
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
		st.addEntry(symbol, address)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
