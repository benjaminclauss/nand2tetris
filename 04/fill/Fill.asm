// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

    @8192
    D=A
    @blocks
    M=D
(INPUT)
    @KBD
    D=M
    @FILL
    D;JGT
    @CLEAR
    0;JMP
    @INPUT
    0;JMP
(FILL)
    // Return to input if already filled.
    @SCREEN
    D=M
    @INPUT
    D;JLT
    // Begin filling.
    @fill
    M=-1
    @START
    0;JMP
(CLEAR)
    // Return to input if already clear.
    @SCREEN
    D=M
    @INPUT
    D;JEQ
    // Begin clearing.
    @fill
    M=0
    @START
    0;JMP
(START)
    @blocks
    D=M
    @remaining
    M=D
    @SCREEN
    D=A
    @address
    M=D
(LOOP)
    // Go back to input if filling/clearing complete.
    @remaining
    D=M
    @INPUT
    D;JEQ
    // Process (fill or clear) block.
    @fill
    D=M
    @address
    A=M
    M=D
    @address
    D=M+1
    @address
    M=D
    // Decrement remaining blocks.
    @remaining
    M=M-1
    @LOOP
    0;JMP

