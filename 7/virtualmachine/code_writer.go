package virtualmachine

import (
	"io"
	"strconv"
	"strings"
)

// A CodeWriter translates VM commands into Hack assembly code.
type CodeWriter struct {
	output   io.WriteCloser
	filename string
	boolean  int
}

// NewCodeWriter opens the output stream and gets ready to write to it.
func NewCodeWriter(output io.WriteCloser) *CodeWriter {
	message := "@256\nD=A\n@SP\nM=D\n"
	io.WriteString(output, message)
	return &CodeWriter{output: output}
}

// SetFileName informs the CodeWriter that the translation of a new VM file is started.
func (cw *CodeWriter) SetFileName(filename string) {
	cw.filename = filename
}

// The VM represents true. and false. as -1 (minus one, 0xFFFF) and 0 (zero, 0x0000), respectively.

// WriteArithmetic writes the assembly code that is the translation of the given arithmetic command.
func (cw *CodeWriter) WriteArithmetic(command string) {
	switch command {
	case "add":
		io.WriteString(cw.output, writeBinaryOperation("+"))
	case "sub":
		io.WriteString(cw.output, writeBinaryOperation("-"))
	case "neg":
		io.WriteString(cw.output, writeUnaryOperation("-"))
	case "eq":
		cw.boolean += 1
		comp := "EQ"
		count := strconv.Itoa(cw.boolean)
		io.WriteString(cw.output, "@SP\nAM=M-1\nD=M\nA=A-1\nD=M-D\n"+
			"@"+comp+".true."+count+"\nD;J"+comp+"\n"+
			"@SP\nA=M-1\nM=0\n@"+comp+".after."+count+"\n"+
			"0;JMP\n("+comp+".true."+count+")\n@SP\nA=M-1\n"+
			"M=-1\n("+comp+".after."+count+")\n")
	case "gt":
		cw.boolean += 1
		comp := "GT"
		count := strconv.Itoa(cw.boolean)
		io.WriteString(cw.output, "@SP\nAM=M-1\nD=M\nA=A-1\nD=M-D\n"+
			"@"+comp+".true."+count+"\nD;J"+comp+"\n"+
			"@SP\nA=M-1\nM=0\n@"+comp+".after."+count+"\n"+
			"0;JMP\n("+comp+".true."+count+")\n@SP\nA=M-1\n"+
			"M=-1\n("+comp+".after."+count+")\n")
	case "lt":
		cw.boolean += 1
		comp := "LT"
		count := strconv.Itoa(cw.boolean)
		io.WriteString(cw.output, "@SP\nAM=M-1\nD=M\nA=A-1\nD=M-D\n"+
			"@"+comp+".true."+count+"\nD;J"+comp+"\n"+
			"@SP\nA=M-1\nM=0\n@"+comp+".after."+count+"\n"+
			"0;JMP\n("+comp+".true."+count+")\n@SP\nA=M-1\n"+
			"M=-1\n("+comp+".after."+count+")\n")
	case "and":
		io.WriteString(cw.output, writeBinaryOperation("&"))
	case "or":
		io.WriteString(cw.output, writeBinaryOperation("|"))
	case "not":
		io.WriteString(cw.output, writeUnaryOperation("!"))
	}
}

func writeBinaryOperation(operator string) string {
	var sb strings.Builder
	// Decrement stack pointer.
	sb.WriteString("@SP\n")
	sb.WriteString("M=M-1")
	// Store first operand.
	sb.WriteString("A=M\n")
	sb.WriteString("D=M\n")
	// Operate.
	sb.WriteString("A=A-1\n")
	sb.WriteString("M=D" + operator + "M\n")
	return sb.String()
}

func writeUnaryOperation(operator string) string {
	var sb strings.Builder
	sb.WriteString("@SP\n")
	sb.WriteString("A=M\n")
	sb.WriteString("A=A-1\n")
	sb.WriteString("M=" + operator + "M\n")
	return sb.String()
}

// 1. Next, handle the segments local, argument, this, and that.
// 2. Next, handle the pointer and temp segments, in particular allowing modification of the bases of the this and that segments.
// 3. Finally, handle the static segment.

// WritePushPop writes the assembly code that is the translation of the given command where command is either CPush or CPop.
func (cw *CodeWriter) WritePushPop(command CommandType, segment string, index int) {
	if command == CPush {
		if segment == "constant" {
			// Store constant in register.
			io.WriteString(cw.output, "@"+strconv.Itoa(index)+"\n")
			io.WriteString(cw.output, "D=A\n")
			// Push value to stack.
			io.WriteString(cw.output, "@SP\n")
			io.WriteString(cw.output, "A=M\n")
			io.WriteString(cw.output, "M=D\n")
			// Increment stack pointer.
			io.WriteString(cw.output, "@SP\n")
			io.WriteString(cw.output, "M=M+1\n")
		}
	}
}

// Close closes the output stream.
func (cw *CodeWriter) Close() error {
	return cw.output.Close()
}
