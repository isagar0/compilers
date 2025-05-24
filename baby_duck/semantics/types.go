package semantics

// ------------------------------------------ VARS ------------------------------------------

// Dictionary: Almacena pares clave-valor
type Dictionary struct {
	items  map[string]interface{} // Mapa para cualquier tipo de valor
	keys   []string               // Slice para ordenar las claves
	parent *Dictionary            // Referencia scope padre
}

// ScopeManager: Maneja la tabla global y una pila de scopes locales
type ScopeManager struct {
	global *Dictionary   // Scope global
	stack  []*Dictionary // Pila scopes locales
}

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

// ------------------------------------------ QUADS ------------------------------------------
// Stack: Last In, Frist Out
type Stack struct {
	items []interface{}
}

// SemanticCube: Define el tipo de resultado
// [tipoIzq][tipoDer][operador] → tipoResultado
type SemanticCube map[string]map[string]map[string]string

// QuadStructure: Estructura de un Quad
type QuadStructure struct {
	Oper   string      // Operador
	Left   interface{} // Operando izquierdo
	Right  interface{} // Operando derecho
	Result interface{} // Resultado
}

// ----------------------------------------- MEMORIA -----------------------------------------

// MemorySegment: Administra un rango de direcciones
type MemorySegment struct {
	start int // Direccion inicial
	end   int // Direccion final
	next  int // Proxima disponible
}

// SegmentGroup: Agrupa segmentos de memoria por catgoría
type SegmentGroup struct {
	Ints    MemorySegment
	Floats  MemorySegment
	Strings MemorySegment
}

// MemoryManager: Administrador principal de memoria
type MemoryManager struct {
	Global   SegmentGroup
	Local    SegmentGroup
	Temp     SegmentGroup
	Constant SegmentGroup
}

// ConstTable: Tabla de constantes para evitar duplicados.
type ConstTable struct {
	ints    []string
	floats  []string
	strings []string
}
