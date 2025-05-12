package semantics

// VariableStructure: Estructura de una variable
type VariableStructure struct {
	Name string // Nombre
	Type string // Tipo
}

// FunctionStructure: Estructura de una función
type FunctionStructure struct {
	Name       string              // Nombre
	Parameters []VariableStructure // Lista de parametros
	VarTable   *Dictionary         // Variables locales (scope local)
}

// QuadStructure: Estructura de un Quad
type QuadStructure struct {
	Oper   string      // Operador
	Left   interface{} // Operando izquierdo
	Right  interface{} // Operando derecho
	Result interface{} // Resultado
}
