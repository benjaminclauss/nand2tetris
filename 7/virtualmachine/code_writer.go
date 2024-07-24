package virtualmachine

import (
	"io"
	"strconv"
)

// A CodeWriter translates VM commands into Hack assembly code.
type CodeWriter struct {
	output io.WriteCloser
}

// NewCodeWriter opens the output stream and gets ready to write to it.
func NewCodeWriter(output io.WriteCloser) *CodeWriter {
	message := "@256\nD=A\n@SP\nM=D\n"
	io.WriteString(output, message)
	return &CodeWriter{output: output}
}

// SetFileName informs the CodeWriter that the translation of a new VM file is started.
func (cw *CodeWriter) SetFileName(filename string) {}

// WriteArithmetic writes the assembly code that is the translation of the given arithmetic command.
func (cw *CodeWriter) WriteArithmetic(command string) {
	switch command {
	case "add":
		// Decrement stack pointer.
		io.WriteString(cw.output, "@SP\n")
		io.WriteString(cw.output, "M=M-1\n")
		// Store first operand.
		io.WriteString(cw.output, "A=M\n")
		io.WriteString(cw.output, "D=M\n")
		// Operate.
		io.WriteString(cw.output, "A=A-1\n")
		io.WriteString(cw.output, "M=D+M\n")
	case "sub":
	case "neg":
	case "eq":
	case "gt":
	case "lt":
	case "and":
	case "or":
	case "not":
	}
}

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
