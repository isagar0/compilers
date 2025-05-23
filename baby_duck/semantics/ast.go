package semantics

import (
	"fmt"
	"strconv"
	"strings"
)

// --------------------------------------- Inicialización ---------------------------------------
var VarTable = NewDictionary()          // Tabla global de variables
var FunctionDirectory = NewDictionary() // Directorio de funciones
var memory = NewMemoryManager()         // Memoria de direcciones
var AddressToName = map[int]string{}    // Traducir direcciones a nombre
var PJumps = NewStack()                 // Stack para saltos pendientes

var (
	PilaO   = NewStack()    // Operandos
	PTypes  = NewStack()    // Tipos de operadores
	POper   = NewStack()    // Operadores
	Quads   []QuadStructure // Cuádruplos generados
	tempVar int             // Contador para nombres de variables temporales
)

// ------------------------------------------ Limpiar ------------------------------------------
// ResetSemanticState: Limpia todo para un programa nuevo
func ResetSemanticState() {
	// Limpia diccionarios
	Reset()
	FunctionDirectory = NewDictionary()

	// Limpia pilas y cuadruplos
	PilaO = NewStack()
	PTypes = NewStack()
	POper = NewStack()
	Quads = []QuadStructure{}
	tempVar = 0

	// Limpia direcciones
	// Reset all memory segments
	memory.Global.Ints.Reset()
	memory.Global.Floats.Reset()
	memory.Local.Ints.Reset()
	memory.Local.Floats.Reset()
	memory.Temp.Ints.Reset()
	memory.Temp.Floats.Reset()
	memory.Constant.Ints.Reset()
	memory.Constant.Floats.Reset()
	memory.Constant.Strings.Reset()
	AddressToName = make(map[int]string)
}

// Tabla de direcciones
func PrintAddressTable() {
	fmt.Println("\n==== Tabla de direcciones virtuales ====")

	// First print variables from scopes
	fmt.Println("\n---- Variabls Globales ----")
	printScopeVariables(scopes.global)
	for _, scope := range scopes.stack {
		printScopeVariables(scope)
	}

	// Then print constants
	fmt.Println("\n---- Constantes ----")
	for addr, name := range AddressToName {
		if strings.HasPrefix(name, "const_") && (addr < 1000 || addr >= 7000) {
			fmt.Printf("%-10s → %d\n", strings.TrimPrefix(name, "const_"), addr)
		}
	}

	// Then print temporaries
	fmt.Println("\n---- Temporales ----")
	for addr, name := range AddressToName {
		if strings.HasPrefix(name, "temp_") {
			fmt.Printf("%-10s → %d\n", name, addr)
		}
	}
}

func printScopeVariables(scope *Dictionary) {
	for _, key := range scope.Keys() {
		if val, exists := scope.Get(key); exists {
			vs := val.(VariableStructure)
			fmt.Printf("%-10s → %-6d (%-6s)\n", vs.Name, vs.Address, vs.Type)
		}
	}
}

// -------------------------------------------- Quads --------------------------------------------
// PushOperandDebug: En lugar de haer push directo lo hace desde acá para debuggear
func PushOperandDebug(value interface{}, tipo string) {
	var address int

	// ¿Es una constante?
	if tipo == "int" || tipo == "float" || tipo == "bool" || tipo == "string" {
		constID := fmt.Sprintf("%v", value)
		// Skip if this looks like a variable address
		if num, err := strconv.Atoi(constID); err == nil && num >= 1000 && num <= 6999 {
			address = num // Use as direct address
		} else {
			address = GetConstAddress(constID, tipo)
		}
		// Only register as const if not in variable range
		if address < 1000 || address >= 7000 {
			AddressToName[address] = fmt.Sprintf("const_%s", constID)
		}
	} else {
		// Si es una variable, buscamos su dirección
		raw, _ := Current().Get(fmt.Sprintf("%v", value))
		v := raw.(VariableStructure)
		address = v.Address
	}

	PilaO.Push(address)
	PTypes.Push(tipo)
}

// PushOp: En lugar de haer push directo lo hace desde acá para debuggear
func PushOp(op string) {
	//fmt.Printf("→ PUSH OPERADOR: %s\n", op)
	POper.Push(op)
}

// NewTemp: Genera el nombre de una variable temporal
func NewTemp() string {
	tempVar++
	return fmt.Sprintf("t%d", tempVar)
}

// PushQuad: Agrega un cuádruplo a la lista global
func PushQuad(oper string, left, right, res interface{}) {
	Quads = append(Quads, QuadStructure{oper, left, right, res})
}

// PrintStacks: Imprime las pilas actuales
func PrintStacks() {
	fmt.Println("\nOperandos:", PilaO.items)
	fmt.Println("Tipos:", PTypes.items)
	fmt.Print("Operadores: [")
	for i := len(POper.items) - 1; i >= 0; i-- {
		fmt.Print(POper.items[i], " ")
	}
	fmt.Println("]")
}

// PrintQuads: Imprime los cuádruplos generados
func PrintQuads() {
	fmt.Println("\nCuádruplos generados:")
	for i, q := range Quads {
		fmt.Printf("%d: (%s %v %v %v)\n", i, q.Oper, q.Left, q.Right, q.Result)
	}
}

// DoAddSub: Genera cuádruplos para operaciones + y -
func DoAddSub() error {
	for {
		// Imprime estado de las pilas
		// PrintStacks()

		// Verifica el tope, si esta vacía termina el ciclo
		top, err := POper.Peek()
		if err != nil {
			return nil
		}

		// Convierte el tope a string
		op := top.(string)

		// Si no es un operador + o -, salimos del ciclo
		if op != "+" && op != "-" {
			break
		}

		// Operandos y tipos
		// Derecho
		rightOp, _ := PilaO.Pop()
		rightType, _ := PTypes.Pop()
		//fmt.Printf("→ rightOp: %v, rightType: %v\n", rightOp, rightType)

		// Izquierdo
		leftOp, _ := PilaO.Pop()
		leftType, _ := PTypes.Pop()
		//fmt.Printf("→ leftOp: %v, leftType: %v\n", leftOp, leftType)

		// Convertir a string y mandar error si no son de ese tipo
		ltype, ok1 := leftType.(string)
		rtype, ok2 := rightType.(string)
		if !ok1 || !ok2 {
			return fmt.Errorf("\nDoAddSub error: tipos no son string: left=%T, right=%T", leftType, rightType)
		}

		// Quitar el operador de la pila
		POper.Pop()

		// Llama al cubo semántico para validar los tipos, si no es valido regresa error
		resType, err := GetResultType(ltype, rtype, op)
		if err != nil {
			return err
		}

		// Genera variable temporal
		var tempAddr int
		switch resType {
		case "int":
			tempAddr, _ = memory.Temp.Ints.GetNext()
		case "float":
			tempAddr, _ = memory.Temp.Floats.GetNext()
		}
		AddressToName[tempAddr] = fmt.Sprintf("temp_%d", tempAddr)

		// Genera el cuádruplo y lo agregamos a la lista global
		PushQuad(op, leftOp, rightOp, tempAddr)

		// Mete la variable temporal y su tipo en las pilas
		PilaO.Push(tempAddr)
		PTypes.Push(resType)

		// Imprime el cuádruplo
		// fmt.Printf("\n→ GENERATE QUAD: %s %v %v -> %v\n", op, leftOp, rightOp, temp)
	}
	return nil
}

// DoMulDiv: Genera cuádruplos para operaciones * y /
func DoMulDiv() error {
	for {
		// Imprime estado de las pilas
		// PrintStacks()

		// Verifica el tope, si esta vacía termina el ciclo
		top, err := POper.Peek()
		if err != nil {
			return nil
		}

		// Convierte el tope a string
		op := top.(string)

		// Si no es un operador * o /, salimos del ciclo
		if op != "*" && op != "/" {
			break
		}

		// Operandos y tipos
		// Derecho
		rightOp, _ := PilaO.Pop()
		rightType, _ := PTypes.Pop()

		// Izquierdo
		leftOp, _ := PilaO.Pop()
		leftType, _ := PTypes.Pop()

		//fmt.Printf("→ DEBUG DoMulDiv: leftType=%T(%v), rightType=%T(%v)\n", leftType, leftType, rightType, rightType)

		// Convertir a string y mandar error si no son de ese tipo
		ltype, ok1 := leftType.(string)
		rtype, ok2 := rightType.(string)
		if !ok1 || !ok2 {
			return fmt.Errorf("DoMulDiv error: tipos no son string: left=%T, right=%T", leftType, rightType)
		}

		// Quitar el operador de la pila
		POper.Pop()

		// Llama al cubo semántico para validar los tipos, si no es valido regresa error
		resType, err := GetResultType(ltype, rtype, op)
		if err != nil {
			return err
		}

		// Genera variable temporal
		var tempAddr int
		switch resType {
		case "int":
			tempAddr, _ = memory.Temp.Ints.GetNext()
		case "float":
			tempAddr, _ = memory.Temp.Floats.GetNext()
		}
		AddressToName[tempAddr] = fmt.Sprintf("temp_%d", tempAddr)

		// Genera el cuádruplo y lo agregamos a la lista global
		PushQuad(op, leftOp, rightOp, tempAddr)

		// Mete la variable temporal y su tipo en las pilas
		PilaO.Push(tempAddr)
		PTypes.Push(resType)

		// Imprime el cuádruplo
		// fmt.Printf("\n→ GENERATE QUAD: %s %v %v -> %v\n", op, leftOp, rightOp, temp)
	}
	return nil
}

// DoRelational: Genera cuádruplos para operadores relacionales <, >, !=
func DoRelational() error {
	// Imprime estado de las pilas
	// PrintStacks()

	// Verifica el tope, si esta vacía termina el ciclo
	top, err := POper.Peek()
	if err != nil {
		return nil // pila vacía
	}

	// Convierte el tope a string
	op := top.(string)
	// Si no es un operador relacional, salimos del ciclo
	if op != "<" && op != ">" && op != "!=" {
		return nil // no es operador relacional
	}

	// Operandos y tipos
	// Derecho
	rightOp, _ := PilaO.Pop()
	rightType, _ := PTypes.Pop()

	// Izquierdo
	leftOp, _ := PilaO.Pop()
	leftType, _ := PTypes.Pop()

	// Convertir a string y mandar error si no son de ese tipo
	ltype, ok1 := leftType.(string)
	rtype, ok2 := rightType.(string)
	if !ok1 || !ok2 {
		return fmt.Errorf("DoRelational error: tipos no son string: left=%T, right=%T", leftType, rightType)
	}

	// Quitar el operador de la pila
	POper.Pop()

	// Llama al cubo semántico para validar los tipos, si no es valido regresa error
	resType, err := GetResultType(ltype, rtype, op)
	if err != nil {
		return err
	}

	// Genera variable temporal
	var tempAddr int
	switch resType {
	case "int":
		tempAddr, _ = memory.Temp.Ints.GetNext()
	case "float":
		tempAddr, _ = memory.Temp.Floats.GetNext()
	}
	AddressToName[tempAddr] = fmt.Sprintf("temp_%d", tempAddr)

	// Genera el cuádruplo y lo agregamos a la lista global
	PushQuad(op, leftOp, rightOp, tempAddr)

	// Mete la variable temporal y su tipo en las pilas
	PilaO.Push(tempAddr)
	PTypes.Push(resType)

	// Imprime el cuádruplo
	// fmt.Printf("→ GENERATE RELATIONAL: %s %v %v -> %v\n", op, leftOp, rightOp, tempAddr)

	return nil
}

// PopUntilFakeBottom: Procesa operadores hasta encontrar el fondo falso (⏊) para ()
func PopUntilFakeBottom() error {
	for {
		// Verifica el tope, si esta vacía termina el ciclo
		top, err := POper.Peek()
		if err != nil {
			break
		}

		// Convertir a string y mandar error si no son de ese tipo
		op := top.(string)

		// Su encuentra fondo falso ⏊, termina procesamiento
		if op == "⏊" {
			POper.Pop() // quitamos la marca
			// fmt.Println("→ POP OPERADOR: ⏊ (fin de paréntesis)")
			break
		}

		// Si es suma o resta, genera su cuádruplo
		if op == "+" || op == "-" {
			err := DoAddSub()
			if err != nil {
				return err
			}

			// Si es multiplicación o división, genera su cuádruplo
		} else if op == "*" || op == "/" {
			err := DoMulDiv()
			if err != nil {
				return err
			}

			// Si encuentra otro operador, error
		} else {
			return fmt.Errorf("operador inesperado en paréntesis: %v", op)
		}
	}
	return nil
}

// -------------------------------------------- Vars --------------------------------------------
// Reset: Crea una nuvea tabla global y sus scopes vacios
func Reset() {
	scopes = &ScopeManager{
		global: NewDictionary(), // Scope global
		stack:  []*Dictionary{}, // Pila soces locales vacio
	}
}

// RegisterMainProgram: Crea el scope global y registra el programa principal
func RegisterMainProgram(programName string) error {
	// Verifica si ya existe una entrada con mismo nombre
	if _, exists := FunctionDirectory.Get(programName); exists {
		return fmt.Errorf("error: el programa '%s' ya ha sido declarado", programName)
	}

	// Registra el programa principal como una función sin parámetros
	FunctionDirectory.Put(programName, FunctionStructure{
		Name:     programName,
		VarTable: Current(),
	})

	//fmt.Printf("Programa principal '%s' registrado exitosamente.\n", programName)
	return nil
}

// VarDeclaration: Procesa la declaración de variables en el scope actual
func VarDeclaration(ids []string, tipo string) error {
	// Usa tabla activa del scope
	tabla := Current()

	// Recorre cada identificador en la lista de variables a declarar
	for _, id := range ids {
		var dir int
		var err error

		// Check current scope
		if _, exists := tabla.Get(id); exists {
			return fmt.Errorf("error: variable '%s' ya declarada en este scope", id)
		}

		// Check parent scopes if in local scope
		if tabla != scopes.global {
			if _, exists := scopes.global.Get(id); exists {
				return fmt.Errorf("error: variable '%s' ya existe en scope global", id)
			}
		}

		// Determina si es global o local
		var segmento *SegmentGroup
		if Current() == scopes.global {
			segmento = &memory.Global
			// fmt.Printf("Global var %s at %d\n", id, dir)
		} else {
			segmento = &memory.Local
			//fmt.Printf("Local var %s at %d\n", id, dir)
		}

		switch tipo {
		case "int":
			dir, err = segmento.Ints.GetNext()
		case "float":
			dir, err = segmento.Floats.GetNext()
		default:
			return fmt.Errorf("tipo no soportado: %s", tipo)
		}
		if err != nil {
			return err
		}

		// Agrega la variable a la tabla con dirección virtual
		tabla.Put(id, VariableStructure{
			Name:    id,
			Type:    tipo,
			Address: dir,
		})
		AddressToName[dir] = id
		//fmt.Printf("Declared %s at address %d (type %s)\n", id, dir, tipo)
	}

	/*
		// Imprimir el contenido del scope actual
		fmt.Println(">>> Contenido del scope actual antes de declarar:")
		tabla.PrintOrdered()
		fmt.Println(">>> Fin del scope actual")
	*/

	return nil
}

// RegisterFunction: Crea la entrada de la función con nombre, retorno void
func RegisterFunction(name string) error {
	// Verifica si ya existe una función con el mismo nombre, marca error
	if _, exists := FunctionDirectory.Get(name); exists {
		return fmt.Errorf("error: función '%s' ya declarada", name)
	}

	// Crea una nueva tabla de variables locales para esta función
	localTable := NewDictionary()

	/*fmt.Printf("[DEBUG] RegisterFunction %s → local scope %p\n",
	name, localTable)*/

	// Registra la función en el directorio
	FunctionDirectory.Put(name, FunctionStructure{
		Name:       name,                  // Nombre
		Parameters: []VariableStructure{}, // Parametros (vacios)
		VarTable:   localTable,            // Tabla local de variables
	})

	return nil
}

// ValidateParams: Verifica que los parámetros de una función no estén duplicados
func ValidateParams(params []VariableStructure) error {
	// Diccionario temporal para llevar el control de nombres ya vistos
	paramSet := NewDictionary()

	// Recorre cada parametro recibido
	for _, param := range params {
		// Verifica si el parametro ya fue declarada, marca error
		if _, exists := paramSet.Get(param.Name); exists {
			return fmt.Errorf("error: parámetro '%s' duplicado en la función", param.Name)
		}

		// Si no existe, se agrega para futuras comparaciones
		paramSet.Put(param.Name, param)
	}
	return nil
}

// FuncDeclaration: Actualiza la entrada creada por RegisterFunction
func FuncDeclaration(name string, params []VariableStructure) error {
	// Verifica que no haya parámetros duplicados
	if err := ValidateParams(params); err != nil {
		return err
	}

	// Busca la función en el directorio, marca error
	raw, exists := FunctionDirectory.Get(name)
	if !exists {
		return fmt.Errorf("error interno: función '%s' no registrada previamente", name)
	}

	// Convierte la entrada a una estructura de función
	fs := raw.(FunctionStructure)

	// Asigna los parámetros recibidos a la función
	fs.Parameters = params

	for _, param := range params {
		// Get the actual address from the local scope
		if raw, exists := Current().Get(param.Name); exists {
			vs := raw.(VariableStructure)
			AddressToName[vs.Address] = param.Name
			// return fmt.Errorf("parameter %s already exists", param.Name)
			//fmt.Printf("Registered param %s → %d (actual address)\n", param.Name, vs.Address)
		}
	}

	// Asocia la tabla local de variables (scope actual donde se declararon los params)
	fs.VarTable = Current()

	// Actualiza el directorio
	FunctionDirectory.Put(name, fs)

	return nil
}
