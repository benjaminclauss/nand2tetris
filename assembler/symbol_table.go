package main

// Keeps a correspondence between symbolic labels and numeric addresses.
type SymbolTable struct {
	symbols map[string]int
}

// Creates a new empty symbol table.
func NewSymbolTable() *SymbolTable {
	return &SymbolTable{make(map[string]int)}
}

// Adds the pair (symbol, address) to the table.
func (st *SymbolTable) addEntry(symbol string, address int) {
	st.symbols[symbol] = address
}

// Does the symbol table contain the given symbol?
func (st *SymbolTable) contains(symbol string) bool {
	_, containsSymbol := st.symbols[symbol]
	return containsSymbol
}

// Returns the address associated with the symbol.
func (st *SymbolTable) GetAddress(symbol string) int {
	return st.symbols[symbol]
}
