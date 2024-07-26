package main

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type CommandType string

const (
	COMMENT_PREFIX = "//"

	// @Xxx where Xxx is either a symbol or a decimal number
	A_COMMAND = CommandType("A")
	// dest=comp;jump
	C_COMMAND = CommandType("C")
	// (actually, pseudo-command) for (Xxx) where Xxx is a symbol
	L_COMMAND = CommandType("L")

	UNRECOGNIZED_COMMAND = CommandType("unrecognized")
)

var (
	A_COMMAND_PATTERN = regexp.MustCompile(`^@.+$`)
	C_COMMAND_PATTERN = regexp.MustCompile(`^(?P<dest>D|A|M|AD|AM|MD)?=?(?P<comp>0|1|-1|D|A|!D|!A|-D|-A|D\+1|A\+1|D-1|A-1|D\+A|D-A|A-D|D&A|D\|A|M|!M|-M|M\+1|M-1|D\+M|D-M|M-D|D&M|D\|M);?(?P<jump>JGT|JEQ|JGE|JLT|JNE|JLE|JMP)?[ ]*(\/\/.*)?$`)
	L_COMMAND_PATTERN = regexp.MustCompile(`^\((?P<variable>.*)\)$`)
)

// Encapsulates access to the input code.
// Reads an assembly language command, parses it, and provides
// convenient access to the command’s components (fields and symbols).
// In addition, removes all white space and comments.
type Parser struct {
	scanner        *bufio.Scanner
	moreCommands   bool
	currentCommand string
}

// Opens the input file/stream and gets ready to parse it.
func NewParser(input io.Reader) *Parser {
	p := &Parser{
		scanner:        bufio.NewScanner(input),
		moreCommands:   true,
		currentCommand: "",
	}
	p.advance()
	return p
}

// Are there more commands in the input?
func (p *Parser) hasMoreCommands() bool {
	return p.moreCommands
}

// Reads the next command from the input and makes it the current command.
// Should be called only if hasMoreCommands() is true.
// Initially there is no current command.
func (p *Parser) advance() {
	if !p.moreCommands {
		return
	}
	moreCommands := p.scanner.Scan()
	if !moreCommands {
		p.moreCommands = false
	}
	p.currentCommand = strings.TrimSpace(p.scanner.Text())
	if len(p.currentCommand) == 0 {
		p.advance()
	} else if strings.HasPrefix(p.currentCommand, COMMENT_PREFIX) {
		p.advance()
	}
}

// Returns the type of the current command.
func (p *Parser) commandType() CommandType {
	switch {
	case A_COMMAND_PATTERN.MatchString(p.currentCommand):
		return A_COMMAND
	case C_COMMAND_PATTERN.MatchString(p.currentCommand):
		return C_COMMAND
	case L_COMMAND_PATTERN.MatchString(p.currentCommand):
		return L_COMMAND
	default:
		return UNRECOGNIZED_COMMAND
	}
}

// Returns the symbol or decimal Xxx of the current command @Xxx or (Xxx).
// Should be called only when commandType() is A_COMMAND or L_COMMAND.
func (p *Parser) symbol() string {
	if A_COMMAND_PATTERN.MatchString(p.currentCommand) {
		split := strings.Split(p.currentCommand, "@")
		addr := split[1]
		return addr
	} else if L_COMMAND_PATTERN.MatchString(p.currentCommand) {
		res := L_COMMAND_PATTERN.FindStringSubmatch(p.currentCommand)
		variable := res[1]
		return variable
	}
	return ""
}

// Returns the dest mnemonic in the current C-command (8 possibilities).
// Should be called only when commandType() is C_COMMAND.
func (p *Parser) dest() string {
	return C_COMMAND_PATTERN.FindStringSubmatch(p.currentCommand)[1]
}

// Returns the comp mnemonic in the current C-command (28 possibilities).
// Should be called only when commandType() is C_COMMAND.
func (p *Parser) comp() string {
	return C_COMMAND_PATTERN.FindStringSubmatch(p.currentCommand)[2]
}

// Returns the jump mnemonic in the current C-command (8 possibilities).
// Should be called only when commandType() is C_COMMAND.
func (p *Parser) jump() string {
	return C_COMMAND_PATTERN.FindStringSubmatch(p.currentCommand)[3]
}
