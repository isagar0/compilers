package semantics

// ------------------------------------------ VARS ------------------------------------------

// Dictionary: Almacena pares clave-valor
type Dictionary struct {
	Items  map[string]interface{} // Mapa para cualquier tipo de valor
	parent *Dictionary            // Referencia scope padre
}

// ScopeManager: Maneja la tabla global y una pila de scopes locales
type ScopeManager struct {
	global  *Dictionary // Scope global
	current *Dictionary // Pila scopes locales
}

// VariableStructure: Estructura de una variable
type VariableStructure struct {
	Name    string // Nombre
	Type    string // Tipo
	Address int    // Direccion
}

// FunctionStructure: Estructura de una función
type FunctionStructure struct {
	Name          string              // Nombre
	Parameters    []VariableStructure // Lista de parametros
	VarTable      *Dictionary         // Variables locales (scope local)
	ParamCount    int                 // Número de parámetros
	LocalVarCount int                 // Variables locales
	StartQuad     int                 // Cuadruplo inicial
	TempCount     int                 // Temporales Usados
}

// FuncInfo: Helper para pasar nombre+params
type FuncInfo struct {
	Name   string
	Params []VariableStructure
}

// ------------------------------------------ QUADS ------------------------------------------

// Stack: Last In, Firsst Out
type Stack struct {
	items []interface{}
}

// SemanticCube: Define el tipo de resultado
// [tipoIzq][tipoDer][operador] → tipoResultado
type SemanticCube map[string]map[string]map[string]string

// QuadStructure: Estructura de un Quad
type QuadStructure struct {
	Oper   int         // Operador
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
	Bools   MemorySegment
}

// MemoryManager: Administrador principal de memoria
type MemoryManager struct {
	Global   SegmentGroup
	Local    SegmentGroup
	Temp     SegmentGroup
	Constant SegmentGroup
}

// ------------------------------------------ VM ------------------------------------------

// ActivationRecord: Contexto de una función llamadora
type ActivationRecord struct {
	ReturnIP int                 // Direccion de instruccion (IP) a la que debe regresar
	LocalMem map[int]interface{} // Memoria local funcion llamadora
}

// VirtualMachine: Nucelo de la VM
type VirtualMachine struct {
	Quads        []QuadStructure              // Lista de cuadruplos a ejecutar
	GlobalMemory map[int]interface{}          // Memoria global (variables + constantes)
	LocalMemory  map[int]interface{}          // Memoria local (función actual)
	IP           int                          // Pointer (índice actual)
	CallStack    []ActivationRecord           // Pila de llamadas (para las funciones)
	FuncDir      map[string]FunctionStructure // Directorio de funciones
	PendingAR    map[int]interface{}          // Registro de activación pendiente (para ERA)
}
