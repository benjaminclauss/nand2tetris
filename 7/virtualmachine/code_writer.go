package virtualmachine

import (
	"io"
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
		io.WriteString(cw.output, "@SP\n")
		io.WriteString(cw.output, "D=M\n")
		io.WriteString(cw.output, "A=A-\n")
		// @address
		// D=M+1
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
func (cw *CodeWriter) WritePushPop(command CommandType, segment string, index int) {}

// Close closes the output stream.
func (cw *CodeWriter) Close() error {
	return cw.output.Close()
}
