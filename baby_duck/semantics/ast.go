package semantics

import (
	"baby_duck/token"
	"fmt"
)

// ------------------------------------------ LIMPIAR ------------------------------------------

// ResetSemanticState: Limpia todo para un programa nuevo
func ResetSemanticState() {
	// Limpia diccionarios
	ResetVars()
	FunctionDirectory = NewDictionary()

	// Limpia pilas y cuadruplos
	ResetStacks()

	// Limpia direcciones
	ResetMemory()
}

// -------------------------------------------- PROGRAM --------------------------------------------

// RegisterMainProgram: Verificación semántica del programa principal
func RegisterMainProgram(programName string) error {
	// Verifica si ya existe una entrada con mismo nombre
	if _, exists := FunctionDirectory.Get(programName); exists {
		return fmt.Errorf("error: el programa '%s' ya ha sido declarado", programName)
	}

	return nil
}

// HandlePHeader: Registra programa y genera cuadruplo inicial
func HandlePHeader(idToken interface{}) (int, error) {
	token := idToken.(*token.Token)
	name := string(token.Lit)

	// Hace la verificación
	if err := RegisterMainProgram(name); err != nil {
		return 0, err
	}

	// Genera cuádruplo y guarda posición para actualizarlo
	PushQuad(GOTO, MAIN, "_", -1)
	gotoMainQuad := len(Quads) - 1

	return gotoMainQuad, nil
}

// HandlePBody: Maneja cuerpo programa (backpatching)
func HandlePBody(gotoMainQuad interface{}) error {
	// Registra la función main
	if err := RegisterFunction("main"); err != nil {
		return err
	}

	// Converte indice a entero
	quadIndex, ok := gotoMainQuad.(int)
	if !ok {
		return fmt.Errorf("error: se esperaba un índice de cuadruplo (int) para el GOTO a main")
	}

	// Obtiene la entrada
	raw, exists := FunctionDirectory.Get("main")
	if !exists {
		return fmt.Errorf("error: función 'main' no encontrada en el directorio de funciones")
	}
	fsMain := raw.(FunctionStructure)
	startMain := fsMain.StartQuad

	// Actualiza el cuadruplo
	if quadIndex < 0 || quadIndex >= len(Quads) {
		return fmt.Errorf("error: índice de cuadruplo GOTO a main (%d) fuera de rango", quadIndex)
	}
	Quads[quadIndex].Result = startMain

	return nil
}

// -------------------------------------------- VARS --------------------------------------------

// Reset: Crea una nuvea tabla global y sus scopes vacios
func ResetVars() {
	global := NewDictionary()
	Scopes = &ScopeManager{
		global:  global,
		current: global,
	}
}

// VarDeclaration: Hace el trabajo de memoria y tablas
func VarDeclaration(ids []string, tipo string) error {
	// Usa tabla activa del scope
	tabla := Scopes.Current()

	// Recorre cada identificador en la lista de variables a declarar
	for _, id := range ids {
		var dir int
		var err error

		// Determina si es global o local
		var segmento *SegmentGroup // Prepara el segmento para asignar la direccion
		if Scopes.Current() == Scopes.global {
			segmento = &memory.Global
			// fmt.Printf("Global var %s at %d\n", id, dir)
		} else {
			segmento = &memory.Local
			//fmt.Printf("Local var %s at %d\n", id, dir)
		}

		// Elegir el subsegmento de variable
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
	}

	/*
		// Imprimir el contenido del scope actual
		fmt.Println(">>> Contenido del scope actual antes de declarar:")
		tabla.PrintOrdered()
		fmt.Println(">>> Fin del scope actual")
	*/

	return nil
}

// CountVars: Cuenta variables en el scope actual
func (d *Dictionary) CountVars() int {
	count := 0
	for _, v := range d.Items {
		if _, ok := v.(VariableStructure); ok {
			count++
		}
	}
	return count
}

// HandleVarDecl: Valida y prepara variables
func HandleVarDecl(ids interface{}, typeToken interface{}) error {
	// Lista variables sea slice string ( ["x", "y"...] )
	idList, ok := ids.([]string)
	if !ok {
		return fmt.Errorf("tipo inválido para lista de IDs")
	}

	// Tipo sea valido (int o float)
	tipoToken, ok := typeToken.(*token.Token)
	if !ok {
		return fmt.Errorf("tipo inválido para tipo de variable")
	}
	tipo := string(tipoToken.Lit)

	// Si es correcto, pasa datos
	return VarDeclaration(idList, tipo)
}

// -------------------------------------------- FUNCS --------------------------------------------

// RegisterFunction: Crea la entrada de la función con nombre, retorno void
func RegisterFunction(name string) error {
	// Verifica si ya existe una función con el mismo nombre, marca error
	if _, exists := FunctionDirectory.Get(name); exists {
		return fmt.Errorf("error: función '%s' ya declarada", name)
	}

	// Crea una nueva tabla de variables locales para esta función
	localTable := NewDictionary()

	// Registra la función en el directorio
	FunctionDirectory.Put(name, FunctionStructure{
		Name:       name,                  // Nombre
		Parameters: []VariableStructure{}, // Parametros (vacios)
		VarTable:   localTable,            // Tabla local de variables
		ParamCount: 0,                     // Numero param
		TempCount:  0,                     // Numero temporales
		StartQuad:  len(Quads),            // Donde empieza
	})

	//fmt.Printf("[DEBUG] RegisterFunction: %d\n", len(Quads))

	return nil
}

// FuncDeclaration: Actualiza la entrada creada por RegisterFunction
func FuncDeclaration(name string, params []VariableStructure, localVarCount, startQuad, tempCount int) error {
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
	for i, param := range params {
		// paramName := fmt.Sprintf("param_%d", i)
		if raw, exists := Scopes.Current().Get(param.Name); exists {
			vs := raw.(VariableStructure)
			AddressToName[vs.Address] = fmt.Sprintf("%s_param_%d", name, i+1)
		}
	}

	// Asigna los nuevos campos
	fs.ParamCount = len(params)      // Número de parámetros
	fs.LocalVarCount = localVarCount // Variables locales
	fs.StartQuad = startQuad         // Cuadruplo inicial
	fs.TempCount = tempCount         // Numero temporales

	// Asocia la tabla local de variables (scope actual donde se declararon los params)
	fs.VarTable = Scopes.Current()

	// Actualiza el directorio
	FunctionDirectory.Put(name, fs)

	return nil
}

// HandleFunctionHeader: Procesa encabezado (nombre y parametros)
func HandleFunctionHeader(idToken, paramsToken interface{}) (FuncInfo, error) {
	token := idToken.(*token.Token)
	name := string(token.Lit)
	params := paramsToken.([]VariableStructure)

	// Registra nombre función
	if err := RegisterFunction(name); err != nil {
		return FuncInfo{}, err
	}

	// Crea scope local y declara los parametros
	Scopes.EnterScope()
	for _, p := range params {
		if err := VarDeclaration([]string{p.Name}, p.Type); err != nil {
			Scopes.ExitScope()
			return FuncInfo{}, err
		}
	}

	// Retorna nombre función con sus parametros
	return FuncInfo{Name: name, Params: params}, nil
}

// HandleFunctionHeaderTwo: Obtiene información importante
func HandleFunctionHeaderTwo(funcInfo interface{}) (FuncInfo, error) {
	info := funcInfo.(FuncInfo)
	localVarCount := Scopes.Current().CountVars() - len(info.Params)
	startQuad := GetCurrentQuad() + 1

	// Registra nombre, parametros, num local vars, donde empieza, y temporales usados
	if err := FuncDeclaration(info.Name, info.Params, localVarCount, startQuad, 0); err != nil {
		Scopes.ExitScope()
		return FuncInfo{}, err
	}
	return info, nil
}

// HandleFEra: Genera ERA
func HandleFEra(idToken interface{}) (interface{}, error) {
	// Extrae nombre de la función
	fnTok, ok := idToken.(*token.Token)
	if !ok {
		return nil, fmt.Errorf("esperaba identificador de función, pero fue %T", idToken)
	}
	name := string(fnTok.Lit)

	// Comprueba que la función exista
	raw, exists := FunctionDirectory.Get(name)
	if !exists {
		return nil, fmt.Errorf("error: función '%s' no declarada", name)
	}

	// Calcula tamaño y genera ERA
	fs := raw.(FunctionStructure)
	size := fs.LocalVarCount + fs.TempCount + fs.ParamCount
	PushQuad(ERA, "_", "_", size)

	return fnTok, nil
}

// HandleFunction: Maneja el cierre de una función
func HandleFunction(funcInfo interface{}) error {
	info, ok := funcInfo.(FuncInfo)
	if !ok {
		return fmt.Errorf("error interno: información de función inválida")
	}

	// Obtiene conteo de temporales (solo es contador en quads)
	tempCount := TempVar

	// Obtiene la entrada de la función
	raw, exists := FunctionDirectory.Get(info.Name)
	if !exists {
		return fmt.Errorf("error: función '%s' no encontrada en el directorio", info.Name)
	}

	// Verifica y convierte el tipo
	fs, ok := raw.(FunctionStructure)
	if !ok {
		return fmt.Errorf("error interno: entrada para '%s' no es FunctionStructure", info.Name)
	}

	// Actualiza el conteo de temporales
	fs.TempCount = tempCount
	FunctionDirectory.Put(info.Name, fs)

	// Reinicia contador para la siguiente función
	TempVar = 0

	// Genera ENDFUNC
	PushQuad(ENDFUNC, "_", "_", "_")

	// Sale del scope local (borra las tablas solo existen dentro de ese scope)
	Scopes.ExitScope()

	return nil
}

// -------------------------------------------- PARAMS --------------------------------------------

// ValidateParams: Verifica que los parámetros de una función no estén duplicados
func ValidateParams(params []VariableStructure) error {
	// Diccionario temporal para llevar el control de nombres ya vistos
	paramSet := NewDictionary()

	// Recorre cada parametro recibido
	for _, param := range params {
		// Si no existe, se agrega para futuras comparaciones
		paramSet.Put(param.Name, param)
	}
	return nil
}

// AssignAddressToParam: Asigna dirección de memoria
func AssignAddressToParam(tipo string) (int, error) {
	switch tipo {
	case "int":
		return memory.Local.Ints.GetNext()
	case "float":
		return memory.Local.Floats.GetNext()
	default:
		return 0, fmt.Errorf("tipo no soportado: %s", tipo)
	}
}

// DeclareInCurrentScope: Decalra parametro en su scope
func DeclareInCurrentScope(name, tipo string, address int) error {
	scope := Scopes.Current()

	// Crea estructura de variable y la agrega
	scope.Put(name, VariableStructure{
		Name:    name,
		Type:    tipo,
		Address: address,
	})

	// Mapea la direccion
	AddressToName[address] = name
	return nil
}

// HandleParam: Procesa un parametro
func HandleParam(idToken, typeToken interface{}) (VariableStructure, error) {
	// Convierte a *token.Token (type, lit, pos) -> Nombre
	nameTok, ok := idToken.(*token.Token)
	if !ok {
		return VariableStructure{}, fmt.Errorf("esperaba token para identificador")
	}

	// Convierte a *token.Token (type, lit, pos) -> Tipo
	tipoTok, ok := typeToken.(*token.Token)
	if !ok {
		return VariableStructure{}, fmt.Errorf("esperaba token para tipo")
	}

	// Cambia a string
	name := string(nameTok.Lit)
	tipo := string(tipoTok.Lit)

	// Asigna dirección según tipo
	dir, err := AssignAddressToParam(tipo)
	if err != nil {
		return VariableStructure{}, err
	}

	// Declara en scope actual
	if err := DeclareInCurrentScope(name, tipo, dir); err != nil {
		return VariableStructure{}, err
	}

	return VariableStructure{Name: name, Type: tipo, Address: dir}, nil
}

// -------------------------------------------- ASSIGN --------------------------------------------

// HandleAssign: Maneja asignaciones
func HandleAssign(idToken interface{}) error {
	name := string(idToken.(*token.Token).Lit)

	// Verifica que variable fue declarada
	if _, exists := Scopes.Current().Get(name); !exists {
		return fmt.Errorf("error: variable '%s' no declarada", name)
	}

	rightOp, _ := PilaO.Pop()            // Operando derecho
	raw, _ := Scopes.Current().Get(name) // Valor guardado en la tabla para nombre
	vs := raw.(VariableStructure)        // Obtienes su info

	// Haces cuadruplo
	PushQuad(ASSIGN, rightOp, "_", vs.Address)
	return nil
}

// -------------------------------------------- CONDITIONAL --------------------------------------------

// HandleConditionTail: Maneja final de condición
func HandleConditionTail() error {
	condAddr, _ := PilaO.Pop()  // Direccion
	condType, _ := PTypes.Pop() // Tipo

	// Condicion debe ser booleana (porque es del if)
	if condType != "bool" {
		return fmt.Errorf("condición debe ser booleana, recibió: %v", condType)
	}

	PushQuad(GOTOF, condAddr, "_", -1) // GOTOF direccion pendiente
	PJumps.Push(len(Quads) - 1)        // Guarda posicion para backpatching
	return nil
}

// HandleElseTail: Final del else
func HandleElseTail() error {
	PushQuad(GOTO, "_", "_", -1) // Salto si ya entro al if
	PJumps.Push(len(Quads) - 1)
	return nil
}

// HandleCycleHeader: Inicio de un ciclo
func HandleCycleHeader() error {
	// Guarda posicion actual asi al final sabe donde saltar volver a revisar condicion
	PJumps.Push(len(Quads))
	return nil
}

// HandleCycleTail: Finaldel ciclo
func HandleCycleTail() error {
	falseJumpRaw, _ := PJumps.Pop()  // Posicion salto condicional
	returnJumpRaw, _ := PJumps.Pop() // Posicion retorno

	falseJump := falseJumpRaw.(int)
	returnJump := returnJumpRaw.(int)

	PushQuad(GOTO, "_", "_", returnJump) // Salto inicio ciclo
	Quads[falseJump].Result = len(Quads) // Completa salto GOTOF
	return nil
}

// HandleCondition: Completa saltos if
func HandleCondition(hasElse bool) error {
	if hasElse {
		// Condiciones con else
		endJumpRaw, _ := PJumps.Pop()   // Salto al final del else
		falseJumpRaw, _ := PJumps.Pop() // Salto al else

		endJump := endJumpRaw.(int)
		falseJump := falseJumpRaw.(int)

		Quads[falseJump].Result = endJump + 1 // Completa el salto al else +1 despues del if
		Quads[endJump].Result = len(Quads)    //Compelta salto final
	} else {
		// Condiciones sin else solo salto GOTOF
		falseJumpRaw, _ := PJumps.Pop()
		falseJump := falseJumpRaw.(int)
		Quads[falseJump].Result = len(Quads)
	}
	return nil
}

// HandleCycleExpression: Expresión de un ciclo
func HandleCycleExpression() error {
	// Extrae la condición de la pila de operandos
	condAddr, err := PilaO.Pop()
	if err != nil {
		return fmt.Errorf("error al extraer condición de la pila: %v", err)
	}

	// Extrae el tipo de la condición
	condTypeRaw, err := PTypes.Pop()
	if err != nil {
		return fmt.Errorf("error al extraer tipo de condición: %v", err)
	}

	condType, ok := condTypeRaw.(string)
	if !ok {
		return fmt.Errorf("tipo de condición no es string: %T", condTypeRaw)
	}

	// Verifica que la condición sea booleana
	if condType != "bool" {
		return fmt.Errorf("condición en while debe ser booleana, recibió: %v", condType)
	}

	// Generar cuadruplo GOTOF con dirección pendiente
	PushQuad(GOTOF, condAddr, "_", -1)

	// Guardar la posición del cuadruplo para backpatching
	PJumps.Push(len(Quads) - 1)

	return nil
}

// -------------------------------------------- PRINT --------------------------------------------

var printArgs []interface{} // Variable temporal para hacerlo a la inversa

// HandlePrintExpression: Expresión dentro de print
func HandlePrintExpression() error {
	value, _ := PilaO.Pop()              // Obtiene valor imprimir
	printArgs = append(printArgs, value) // Agrega a lista temporal
	return nil
}

// HandlePrintString: Maneja string
func HandlePrintString(strToken interface{}) error {
	tok := strToken.(*token.Token)
	str := string(tok.Lit)
	addr := GetConstAddress(str, "string")
	printArgs = append(printArgs, addr)
	return nil
}

// FinalizePrint: Genera cuadruplos
func FinalizePrint() error {
	// Genera cuadruplos en orden inverso
	for i := len(printArgs) - 1; i >= 0; i-- {
		PushQuad(PRINT, printArgs[i], "_", "_")
	}
	printArgs = nil // Limpia la lista temporal
	return nil
}
