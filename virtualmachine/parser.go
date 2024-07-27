package virtualmachine

import (
	"bufio"
	"io"
	"strings"
)

type CommandType string

const (
	CArithmetic CommandType = "C_ARITHMETIC"
	CPush       CommandType = "C_PUSH"
	CPop        CommandType = "C_POP"
	CLabel      CommandType = "C_LABEL"
	CGoTo       CommandType = "C_GOTO"
	CIf         CommandType = "C_IF"
	CFuntion    CommandType = "C_FUNCTION"
	CReturn     CommandType = "C_RETURN"
	CCall       CommandType = "C_CALL"
)

// A Parser handles the parsing of a single .vm file, and encapsulates access to the input code.
// It reads VM commands, parses them, and provides convenient access to their components.
// In addition, it removes all white space and comments.
type Parser struct {
	input *bufio.Scanner

	currentCommand string
	nextCommand    string
}

// NewParser opens the input stream and gets ready to parse it.
func NewParser(input io.Reader) *Parser {
	parser := &Parser{input: bufio.NewScanner(input)}
	parser.advance()
	return parser
}

// HasMoreCommands reutrns true if there are more commands in the input.
func (p *Parser) HasMoreCommands() bool {
	return p.nextCommand != ""
}

// Advance reads the next command from the input and makes it the current command.
// This routine should only be called only if HasMoreCommands is true.
// Initially, there is no current command.
func (p *Parser) Advance() {
	p.currentCommand = p.nextCommand
	p.nextCommand = ""
	p.advance()
}

func (p *Parser) advance() {
	/*
	   Within a .vm file, each VM command appears in a separate line, and in one of the following formats:
	     - command (e.g., add)
	     - command arg (e.g., goto loop)
	     - command arg1 arg2 (e.g., push local 3).

	   The arguments are separated from each other and from the command part by an arbitrary number of spaces.
	*/
	for p.input.Scan() {
		text := p.input.Text()
		// “//” comments can appear at the end of any line and are ignored. Blank lines are permitted and ignored.
		if strings.HasPrefix(text, "//") || len(strings.TrimSpace(text)) == 0 {
			continue
		} else {
			p.nextCommand = text
			break
		}
	}
}

func (p *Parser) CommandType() CommandType {
	parts := strings.Fields(p.currentCommand)
	return map[string]CommandType{
		"add":      CArithmetic,
		"sub":      CArithmetic,
		"neg":      CArithmetic,
		"eq":       CArithmetic,
		"gt":       CArithmetic,
		"lt":       CArithmetic,
		"and":      CArithmetic,
		"or":       CArithmetic,
		"not":      CArithmetic,
		"push":     CPush,
		"pop":      CPop,
		"label":    CLabel,
		"goto":     CGoTo,
		"if-goto":  CIf,
		"function": CFuntion,
		"call":     CCall,
		"return":   CReturn,
	}[parts[0]]
}

// Arg1 returns the first argument of the current command.
// In the case of CArithmetic, the command itself (add, sub, etc.) is returned
// It should not be called if the current command is CReturn.
func (p *Parser) Arg1() string {
	parts := strings.Fields(p.currentCommand)
	if p.CommandType() == CArithmetic {
		return parts[0]
	}
	return parts[1]
}

// Arg2 returns the second argument of the current command.
// It should be called only if the current command is CPush, CPop, CFunction or CCall.
func (p *Parser) Arg2() string {
	parts := strings.Fields(p.currentCommand)
	return parts[2]
}
