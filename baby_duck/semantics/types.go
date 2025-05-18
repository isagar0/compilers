package semantics

// VariableStructure: Estructura de una variable
type VariableStructure struct {
	Name    string // Nombre
	Type    string // Tipo
	Address int    // Direccion
}

// FunctionStructure: Estructura de una función
type FunctionStructure struct {
	Name       string              // Nombre
	Parameters []VariableStructure // Lista de parametros
	VarTable   *Dictionary         // Variables locales (scope local)
}

// SemanticCube: Estructura de 3 niveles que define el tipo de resultado
// Estructura: [tipoIzq][tipoDer][operador] → tipoResultado
type SemanticCube map[string]map[string]map[string]string

// QuadStructure: Estructura de un Quad
type QuadStructure struct {
	Oper   string      // Operador
	Left   interface{} // Operando izquierdo
	Right  interface{} // Operando derecho
	Result interface{} // Resultado
}

// MemorySegment: Administra un rango de direcciones
type MemorySegment struct {
	start int
	end   int
	next  int
}
