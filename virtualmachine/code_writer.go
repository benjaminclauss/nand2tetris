package virtualmachine

import (
	"io"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// A CodeWriter translates VM commands into Hack assembly code.
// The VM represents true. and false. as -1 (minus one, 0xFFFF) and 0 (zero, 0x0000), respectively.
type CodeWriter struct {
	output              io.WriteCloser
	boolean             int
	returnLabelCount    int
	currentFunctionName string
	currentFilename     string
}

// NewCodeWriter opens the output stream and gets ready to write to it.
func NewCodeWriter(output io.WriteCloser) *CodeWriter {
	return &CodeWriter{output: output, returnLabelCount: 1}
}

// SetFilename informs the CodeWriter that the translation of a new VM file is started.
func (cw *CodeWriter) SetFilename(filename string) {
	cw.currentFilename = filename
}

// WriteArithmetic writes the assembly code that is the translation of the given arithmetic command.
func (cw *CodeWriter) WriteArithmetic(command string) {
	switch command {
	case "add", "sub", "and", "or":
		cw.writeBinaryOperation(command)
	case "neg", "not":
		cw.writeUnaryOperation(command)
	case "eq", "gt", "lt":
		cw.writeComparison(command)
	}
}

func (cw *CodeWriter) writeBinaryOperation(command string) {
	operation := map[string]string{"add": "D+M", "sub": "M-D", "and": "D&M", "or": "D|M"}[command]
	// Decrement stack pointer.
	io.WriteString(cw.output, "@SP\nM=M-1\n")
	// Store first operand.
	io.WriteString(cw.output, "A=M\nD=M\n")
	// Operate.
	io.WriteString(cw.output, "A=A-1\n")
	io.WriteString(cw.output, "M="+operation+"\n")
}

func (cw *CodeWriter) writeUnaryOperation(command string) {
	result := map[string]string{"neg": "-M", "not": "!M"}[command]
	io.WriteString(cw.output, "@SP\nA=M\n")
	io.WriteString(cw.output, "A=A-1\n")
	io.WriteString(cw.output, "M="+result+"\n")
}

// TODO: Clean.
func (cw *CodeWriter) writeComparison(command string) {
	cw.boolean += 1
	comp := strings.ToUpper(command)
	count := strconv.Itoa(cw.boolean)
	io.WriteString(cw.output, "@SP\nAM=M-1\nD=M\nA=A-1\nD=M-D\n"+
		"@"+comp+".true."+count+"\nD;J"+comp+"\n"+
		"@SP\nA=M-1\nM=0\n@"+comp+".after."+count+"\n"+
		"0;JMP\n("+comp+".true."+count+")\n@SP\nA=M-1\n"+
		"M=-1\n("+comp+".after."+count+")\n")
}

// WritePushPop writes the assembly code that is the translation of the given command where command is either CPush or CPop.
func (cw *CodeWriter) WritePushPop(command CommandType, segment string, index int) {
	// TODO: Dry.
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
		if slices.Contains([]string{"local", "argument", "this", "that", "temp", "pointer"}, segment) {
			name := map[string]string{"local": "LCL", "argument": "ARG", "this": "THIS", "that": "THAT", "temp": "5", "pointer": "3"}[segment]
			// Get address of segment.
			io.WriteString(cw.output, "@"+name+"\n")
			which := "M"
			if segment == "temp" || segment == "pointer" {
				which = "A"
			}
			io.WriteString(cw.output, "D="+which+"\n")
			// Add offset.
			io.WriteString(cw.output, "@"+strconv.Itoa(index)+"\n")
			io.WriteString(cw.output, "D=D+A\n")
			// Store value in register.
			io.WriteString(cw.output, "A=D\n")
			io.WriteString(cw.output, "D=M\n")
			// Push value to stack.
			io.WriteString(cw.output, "@SP\n")
			io.WriteString(cw.output, "A=M\n")
			io.WriteString(cw.output, "M=D\n")
			// Increment stack pointer.
			io.WriteString(cw.output, "@SP\n")
			io.WriteString(cw.output, "M=M+1\n")
		}
		if segment == "static" {
			// Get address of segment.
			// Each static variable j in a VM file Xxx. vm is translated into the assembly symbol Xxx.j.
			// In the subsequent assembly process, these symbolic variables will be allocated RAM space by the Hack assembler.
			filename := strings.TrimRight(filepath.Base(cw.currentFilename), ".vm")
			io.WriteString(cw.output, "@"+filename+"."+strconv.Itoa(index)+"\n")
			io.WriteString(cw.output, "D=M\n")
			// Push value to stack.
			io.WriteString(cw.output, "@SP\n")
			io.WriteString(cw.output, "A=M\n")
			io.WriteString(cw.output, "M=D\n")
			// Increment stack pointer.
			io.WriteString(cw.output, "@SP\n")
			io.WriteString(cw.output, "M=M+1\n")
		}
	}
	if command == CPop {
		name := map[string]string{"local": "LCL", "argument": "ARG", "this": "THIS", "that": "THAT", "temp": "5", "pointer": "3"}[segment]
		if slices.Contains([]string{"local", "argument", "this", "that", "temp", "pointer"}, segment) {
			io.WriteString(cw.output, "@"+name+"\n")
			which := "M"
			if segment == "temp" || segment == "pointer" {
				which = "A"
			}
			io.WriteString(cw.output, "D="+which+"\n")
			// Add offset.
			io.WriteString(cw.output, "@"+strconv.Itoa(index)+"\n")
			io.WriteString(cw.output, "D=D+A\n")
			// Save for later.
			io.WriteString(cw.output, "@13\n")
			io.WriteString(cw.output, "M=D\n")
			// Decrement stack pointer.
			io.WriteString(cw.output, "@SP\nM=M-1\n")
			// Store value.
			io.WriteString(cw.output, "A=M\nD=M\n")
			// Point at segment.
			io.WriteString(cw.output, "@13\n")
			io.WriteString(cw.output, "A=M\n")
			// Save value.
			io.WriteString(cw.output, "M=D\n")
		}
		if segment == "static" {
			// Decrement stack pointer.
			io.WriteString(cw.output, "@SP\nM=M-1\n")
			// Store value.
			io.WriteString(cw.output, "A=M\nD=M\n")
			// Point at segment.
			// Each static variable j in a VM file Xxx. vm is translated into the assembly symbol Xxx.j.
			// In the subsequent assembly process, these symbolic variables will be allocated RAM space by the Hack assembler.
			filename := strings.TrimRight(filepath.Base(cw.currentFilename), ".vm")
			io.WriteString(cw.output, "@"+filename+"."+strconv.Itoa(index)+"\n")
			// Save value.
			io.WriteString(cw.output, "M=D\n")
		}
	}
}

// local, argument, this, that: Each one of these segments is mapped directly on the RAM, and its location is maintained by keeping its physical base address in a dedicated register
// (LCL, ARG, THIS, and THAT, respectively). Thus any access to the ith entry of any one of these segments should be translated to assembly code that accesses address (base + i)
// in the RAM, where base is the current value stored in the register dedicated to the respective segment.

// WriteInit writes assembly code that effects the VM initialization, also called bootstrap code.
// This code must be placed at the beginning of the output file.
func (cw *CodeWriter) WriteInit() error {
	cw.WriteLine("@256\nD=A\n@SP\nM=D\n")
	cw.WriteCall("Sys.init", 0)
	// Do I need to 0;JMP after?
	return nil
}

func (cw *CodeWriter) WriteLine(s string) error {
	_, err := io.WriteString(cw.output, s)
	return err
}

// WriteLabel writes assembly code that effects the `label` command.
func (cw *CodeWriter) WriteLabel(label string) error {
	return cw.WriteLine("(" + cw.currentFunctionName + "$" + label + ")\n")
}

// WriteGoto writes assembly code that effects the `goto` command.
func (cw *CodeWriter) WriteGoto(label string) error {
	return cw.WriteLine("@" + cw.currentFunctionName + "$" + label + "\n0;JMP\n")
}

// WriteIf writes assembly code that effects the if-goto command.
//
// This command effects a conditional goto operation.
// The stack’s topmost value is popped.
// If the value is not zero, execution continues from the location marked by the label.
// Otherwise, execution continues from the next command in the program.
// The jump destination must be located in the same function.
func (cw *CodeWriter) WriteIf(label string) error {
	// Decrement stack pointer.
	io.WriteString(cw.output, "@SP\nM=M-1\n")
	// // Store value.
	io.WriteString(cw.output, "A=M\nD=M\n")
	// io.WriteString(cw.output, "@"+label+"\n")
	// io.WriteString(cw.output, "D;JGT"+"\n")
	// return "@" + cw.funcName + "$" + label + "\n0;JMP\n"
	return cw.WriteLine("@" + cw.currentFunctionName + "$" + label + "\nD;JNE\n")
}

// WriteCall writes assembly code that effects the call command.
//
// This command effects an unconditional goto operation, causing execution to continue from the location marked by the label.
// The jump destination must be located in the same function.
func (cw *CodeWriter) WriteCall(functionName string, numArgs int) error {
	count := strconv.Itoa(cw.returnLabelCount)
	cw.returnLabelCount += 1
	io.WriteString(cw.output, "@SP\nD=M\n@R13\nM=D\n")
	io.WriteString(cw.output, "@ret."+count+"\nD=A\n@SP\nA=M\nM=D\n")
	io.WriteString(cw.output, "@SP\nM=M+1\n")
	io.WriteString(cw.output, "@"+"LCL"+"\nD=M\n@SP\nA=M\nM=D\n")
	io.WriteString(cw.output, "@SP\nM=M+1\n")
	io.WriteString(cw.output, "@"+"ARG"+"\nD=M\n@SP\nA=M\nM=D\n")
	io.WriteString(cw.output, "@SP\nM=M+1\n")
	io.WriteString(cw.output, "@"+"THIS"+"\nD=M\n@SP\nA=M\nM=D\n")
	io.WriteString(cw.output, "@SP\nM=M+1\n")
	io.WriteString(cw.output, "@"+"THAT"+"\nD=M\n@SP\nA=M\nM=D\n")
	io.WriteString(cw.output, "@SP\nM=M+1\n")
	io.WriteString(cw.output, "@R13\nD=M\n@"+strconv.Itoa(numArgs)+"\nD=D-A\n@ARG\nM=D\n"+
		"@SP\nD=M\n@LCL\nM=D\n@"+functionName+"\n"+
		"0;JMP\n(ret."+count+")\n")
	return nil
}

// WriteReturn writes assembly code that effects the return command.
func (cw *CodeWriter) WriteReturn() error {
	cw.WriteLine("@LCL\nD=M\n@5\nA=D-A\nD=M\n@R13\nM=D\n" +
		"@SP\nA=M-1\nD=M\n@ARG\nA=M\nM=D\n" +
		"D=A+1\n@SP\nM=D\n" +
		"@LCL\nAM=M-1\nD=M\n@THAT\nM=D\n" +
		"@LCL\nAM=M-1\nD=M\n@THIS\nM=D\n" +
		"@LCL\nAM=M-1\nD=M\n@ARG\nM=D\n" +
		"@LCL\nA=M-1\nD=M\n@LCL\nM=D\n" +
		"@R13\nA=M\n0;JMP\n")
	return nil
	// // FRAME = LCL; FRAME is a temporary variable. (R14 is FRAME)
	// io.WriteString(cw.output, "@LCL\n")
	// io.WriteString(cw.output, "D=M\n")
	// io.WriteString(cw.output, "@14\n")
	// io.WriteString(cw.output, "M=D\n")
	// // RET = *(FRAME - 5); Put the return address in a temporary variable.
	// io.WriteString(cw.output, "@5\n")
	// io.WriteString(cw.output, "D=D-A\n")
	// io.WriteString(cw.output, "A=D\n")
	// io.WriteString(cw.output, "D=M\n")
	// io.WriteString(cw.output, "@15\n")
	// io.WriteString(cw.output, "M=D\n")

	// // *ARG = pop
	// // Decrement stack pointer.
	// io.WriteString(cw.output, "@SP\nM=M-1\n")
	// // Store value.
	// io.WriteString(cw.output, "A=M\nD=M\n")
	// io.WriteString(cw.output, "@ARG\n")
	// io.WriteString(cw.output, "A=M\n")
	// io.WriteString(cw.output, "M=D\n")

	// // SP = ARG + 1
	// io.WriteString(cw.output, "@ARG\n")
	// // TODO: Can I consolidate this to M+1?
	// io.WriteString(cw.output, "D=M\n")
	// io.WriteString(cw.output, "D=D+1\n")
	// io.WriteString(cw.output, "@SP\n")
	// io.WriteString(cw.output, "M=D\n")

	// // THAT = *(FRAME-1)
	// io.WriteString(cw.output, "@14\n")
	// io.WriteString(cw.output, "D=M\n")
	// io.WriteString(cw.output, "@1\n")
	// io.WriteString(cw.output, "D=D-A\n")
	// io.WriteString(cw.output, "A=D\n")
	// io.WriteString(cw.output, "D=M\n")
	// io.WriteString(cw.output, "@THAT\n")
	// io.WriteString(cw.output, "M=D\n")

	// // THIS = *(FRAME-2)
	// io.WriteString(cw.output, "@14\n")
	// io.WriteString(cw.output, "D=M\n")
	// io.WriteString(cw.output, "@2\n")
	// io.WriteString(cw.output, "D=D-A\n")
	// io.WriteString(cw.output, "A=D\n")
	// io.WriteString(cw.output, "D=M\n")
	// io.WriteString(cw.output, "@THIS\n")
	// io.WriteString(cw.output, "M=D\n")

	// // ARG = *(FRAME-3)
	// io.WriteString(cw.output, "@14\n")
	// io.WriteString(cw.output, "D=M\n")
	// io.WriteString(cw.output, "@3\n")
	// io.WriteString(cw.output, "D=D-A\n")
	// io.WriteString(cw.output, "A=D\n")
	// io.WriteString(cw.output, "D=M\n")
	// io.WriteString(cw.output, "@ARG\n")
	// io.WriteString(cw.output, "M=D\n")

	// // LCL = *(FRAME-4)
	// io.WriteString(cw.output, "@14\n")
	// io.WriteString(cw.output, "D=M\n")
	// io.WriteString(cw.output, "@4\n")
	// io.WriteString(cw.output, "D=D-A\n")
	// io.WriteString(cw.output, "A=D\n")
	// io.WriteString(cw.output, "D=M\n")
	// io.WriteString(cw.output, "@LCL\n")
	// io.WriteString(cw.output, "M=D\n")

	// // goto RET
	// io.WriteString(cw.output, "@15\n")
	// io.WriteString(cw.output, "D=M\n")
	// // “we assume that the A register has been previously set to the address to which we have to jump.”
	// io.WriteString(cw.output, "A=D\n")
	// io.WriteString(cw.output, "0;JMP\n")
	// return nil
}

// WriteFunction writes assembly code that effects the function command.
func (cw *CodeWriter) WriteFunction(functionName string, numLocals int) error {
	cw.currentFunctionName = functionName
	s := "(" + functionName + ")\n@SP\nA=M\n"
	for i := 0; i < numLocals; i++ {
		s += "M=0\nA=A+1\n"
	}
	cw.WriteLine(s + "D=A\n@SP\nM=D\n")
	// io.WriteString(cw.output, fmt.Sprintf("(%s)\n", functionName))
	// // for i := 0; i < numLocals; i++ {
	// // 	cw.WritePushPop(CPush, "constant", 0)
	// // }
	// return nil
	return nil
}

// Close closes the output stream.
func (cw *CodeWriter) Close() error {
	return cw.output.Close()
}
