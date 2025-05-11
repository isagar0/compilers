package semantics

// Estructura de una variable
type VariableStructure struct {
	Name string
	Type string
}

// Estructura de una función
type FunctionStructure struct {
	Name       string
	Parameters []VariableStructure
	VarTable   *Dictionary
}
