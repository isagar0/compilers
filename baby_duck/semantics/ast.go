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
// RegisterMainProgram: Crea el scope global y registra el programa principal
func RegisterMainProgram(programName string) error {
	// Verifica si ya existe una entrada con mismo nombre
	if _, exists := FunctionDirectory.Get(programName); exists {
		return fmt.Errorf("error: el programa '%s' ya ha sido declarado", programName)
	}

	return nil
}

func HandlePHeader(idToken interface{}) (int, error) {
	token := idToken.(*token.Token)
	name := string(token.Lit)

	if err := RegisterMainProgram(name); err != nil {
		return 0, err
	}

	PushQuad(GOTO, MAIN, "_", -1)
	gotoMainQuad := len(Quads) - 1
	return gotoMainQuad, nil
}

func HandlePBody(gotoMainQuad interface{}) error {
	// Registrar la función 'main' (si no está registrada)
	if err := RegisterFunction("main"); err != nil {
		return err
	}

	// Convertir el atributo (que es un int) a entero
	quadIndex, ok := gotoMainQuad.(int)
	if !ok {
		return fmt.Errorf("error: se esperaba un índice de cuadruplo (int) para el GOTO a main")
	}

	// Obtener la entrada de la función main
	raw, exists := FunctionDirectory.Get("main")
	if !exists {
		return fmt.Errorf("error: función 'main' no encontrada en el directorio de funciones")
	}
	fsMain := raw.(FunctionStructure)
	startMain := fsMain.StartQuad

	// Actualizar el cuadruplo GOTO con la dirección de inicio de main
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
		current: global, // Ahora current apunta a global inicialmente
	}
}

// VarDeclaration: Procesa la declaración de variables en el scope actual
func VarDeclaration(ids []string, tipo string) error {
	// Usa tabla activa del scope
	tabla := Scopes.Current()

	// Recorre cada identificador en la lista de variables a declarar
	for _, id := range ids {
		var dir int
		var err error

		// Check parent scopes if in local scope
		/*
			if tabla != Scopes.global {
				if _, exists := Scopes.global.Get(id); exists {
					return fmt.Errorf("error: variable '%s' ya existe en scope global", id)
				}
			}
		*/

		// Determina si es global o local
		var segmento *SegmentGroup
		if Scopes.Current() == Scopes.global {
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

func HandleVarDecl(ids interface{}, typeToken interface{}) error {
	idList, ok := ids.([]string)
	if !ok {
		return fmt.Errorf("tipo inválido para lista de IDs")
	}

	tipoToken, ok := typeToken.(*token.Token)
	if !ok {
		return fmt.Errorf("tipo inválido para tipo de variable")
	}
	tipo := string(tipoToken.Lit)

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
		ParamCount: 0,
		TempCount:  0,
		StartQuad:  len(Quads),
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

func HandleFunctionHeader(idToken, paramsToken interface{}) (FuncInfo, error) {
	token := idToken.(*token.Token)
	name := string(token.Lit)
	params := paramsToken.([]VariableStructure)

	if err := RegisterFunction(name); err != nil {
		return FuncInfo{}, err
	}

	Scopes.EnterScope()
	for _, p := range params {
		if err := VarDeclaration([]string{p.Name}, p.Type); err != nil {
			Scopes.ExitScope()
			return FuncInfo{}, err
		}
	}
	return FuncInfo{Name: name, Params: params}, nil
}

func HandleFunctionHeaderTwo(funcInfo interface{}) (FuncInfo, error) {
	info := funcInfo.(FuncInfo)
	localVarCount := Scopes.Current().CountVars() - len(info.Params)
	startQuad := GetCurrentQuad() + 1

	if err := FuncDeclaration(info.Name, info.Params, localVarCount, startQuad, 0); err != nil {
		Scopes.ExitScope()
		return FuncInfo{}, err
	}
	return info, nil
}

func HandleFEra(idToken interface{}) (interface{}, error) {
	// 1) Extraer nombre de la función
	fnTok, ok := idToken.(*token.Token)
	if !ok {
		return nil, fmt.Errorf("esperaba identificador de función, pero fue %T", idToken)
	}
	name := string(fnTok.Lit)

	// 2) Comprobar que la función exista
	raw, exists := FunctionDirectory.Get(name)
	if !exists {
		return nil, fmt.Errorf("error: función '%s' no declarada", name)
	}

	// 3) Calcular tamaño y generar ERA
	fs := raw.(FunctionStructure)
	size := fs.LocalVarCount + fs.TempCount + fs.ParamCount
	PushQuad("ERA", "_", "_", size)

	return fnTok, nil
}

func HandleFunction(funcInfo interface{}) error {
	info, ok := funcInfo.(FuncInfo)
	if !ok {
		return fmt.Errorf("error interno: información de función inválida")
	}

	// Obtener el conteo actual de temporales
	tempCount := TempVar

	// Obtener la entrada de la función
	raw, exists := FunctionDirectory.Get(info.Name)
	if !exists {
		return fmt.Errorf("error: función '%s' no encontrada en el directorio", info.Name)
	}

	// Verificar y convertir el tipo
	fs, ok := raw.(FunctionStructure)
	if !ok {
		return fmt.Errorf("error interno: entrada para '%s' no es FunctionStructure", info.Name)
	}

	// Actualizar el conteo de temporales
	fs.TempCount = tempCount
	FunctionDirectory.Put(info.Name, fs)

	// Reiniciar contador para la siguiente función
	TempVar = 0

	// Generar ENDFUNC
	PushQuad("ENDFUNC", "_", "_", "_")

	// Salir del scope local
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

func DeclareInCurrentScope(name, tipo string, address int) error {
	scope := Scopes.Current()

	scope.Put(name, VariableStructure{
		Name:    name,
		Type:    tipo,
		Address: address,
	})
	AddressToName[address] = name
	return nil
}

func HandleParam(idToken, typeToken interface{}) (VariableStructure, error) {
	nameTok, ok := idToken.(*token.Token)
	if !ok {
		return VariableStructure{}, fmt.Errorf("esperaba token para identificador")
	}
	tipoTok, ok := typeToken.(*token.Token)
	if !ok {
		return VariableStructure{}, fmt.Errorf("esperaba token para tipo")
	}

	name := string(nameTok.Lit)
	tipo := string(tipoTok.Lit)

	dir, err := AssignAddressToParam(tipo)
	if err != nil {
		return VariableStructure{}, err
	}

	if err := DeclareInCurrentScope(name, tipo, dir); err != nil {
		return VariableStructure{}, err
	}

	return VariableStructure{Name: name, Type: tipo, Address: dir}, nil
}

// -------------------------------------------- ASSIGN --------------------------------------------

func HandleAssign(idToken interface{}) error {
	name := string(idToken.(*token.Token).Lit)

	if _, exists := Scopes.Current().Get(name); !exists {
		return fmt.Errorf("error: variable '%s' no declarada", name)
	}

	rightOp, _ := PilaO.Pop()
	raw, _ := Scopes.Current().Get(name)
	vs := raw.(VariableStructure)

	PushQuad("=", rightOp, "_", vs.Address)
	return nil
}

// -------------------------------------------- CONDITIONAL --------------------------------------------

func HandleConditionTail() error {
	condAddr, _ := PilaO.Pop()
	condType, _ := PTypes.Pop()

	if condType != "bool" {
		return fmt.Errorf("condición debe ser booleana, recibió: %v", condType)
	}

	PushQuad("GOTOF", condAddr, "_", -1)
	PJumps.Push(len(Quads) - 1)
	return nil
}

func HandleElseTail() error {
	PushQuad("GOTO", "_", "_", -1)
	PJumps.Push(len(Quads) - 1)
	return nil
}

func HandleCycleHeader() error {
	PJumps.Push(len(Quads))
	return nil
}

func HandleCycleTail() error {
	falseJumpRaw, _ := PJumps.Pop()
	returnJumpRaw, _ := PJumps.Pop()

	falseJump := falseJumpRaw.(int)
	returnJump := returnJumpRaw.(int)

	PushQuad("GOTO", "_", "_", returnJump)
	Quads[falseJump].Result = len(Quads)
	return nil
}

func HandleCondition(hasElse bool) error {
	if hasElse {
		endJumpRaw, _ := PJumps.Pop()
		falseJumpRaw, _ := PJumps.Pop()
		endJump := endJumpRaw.(int)
		falseJump := falseJumpRaw.(int)

		Quads[falseJump].Result = endJump + 1
		Quads[endJump].Result = len(Quads)
	} else {
		falseJumpRaw, _ := PJumps.Pop()
		falseJump := falseJumpRaw.(int)
		Quads[falseJump].Result = len(Quads)
	}
	return nil
}

func HandleCycleExpression() error {
	// Extraer la condición de la pila de operandos
	condAddr, err := PilaO.Pop()
	if err != nil {
		return fmt.Errorf("error al extraer condición de la pila: %v", err)
	}

	// Extraer el tipo de la condición
	condTypeRaw, err := PTypes.Pop()
	if err != nil {
		return fmt.Errorf("error al extraer tipo de condición: %v", err)
	}

	condType, ok := condTypeRaw.(string)
	if !ok {
		return fmt.Errorf("tipo de condición no es string: %T", condTypeRaw)
	}

	// Verificar que la condición sea booleana
	if condType != "bool" {
		return fmt.Errorf("condición en while debe ser booleana, recibió: %v", condType)
	}

	// Generar cuadruplo GOTOF (salto si falso) con dirección pendiente
	PushQuad("GOTOF", condAddr, "_", -1)

	// Guardar la posición del cuadruplo para backpatching
	PJumps.Push(len(Quads) - 1)

	return nil
}

// -------------------------------------------- PRINT --------------------------------------------

func HandlePrintExpression() error {
	value, _ := PilaO.Pop()
	PushQuad("PRINT", value, "_", "_")
	return nil
}

func HandlePrintString(strToken interface{}) error {
	tok := strToken.(*token.Token)
	str := string(tok.Lit)
	addr := GetConstAddress(str, "string")
	PushQuad("PRINT", addr, "_", "_")
	return nil
}
