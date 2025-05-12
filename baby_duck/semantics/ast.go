package semantics

import "fmt"

var VarTable = NewDictionary()
var FunctionDirectory = NewDictionary()

// Pilas para cuadruplos
var (
	PilaO   = NewStack()
	PTypes  = NewStack()
	POper   = NewStack()
	Quads   []QuadStructure
	tempVar int
)

// ------------------------------------ Expreciones
func PushOperandDebug(value interface{}, tipo string) {
	fmt.Printf("→ PUSH OPERAND: %v (type: %s)\n", value, tipo)
	PilaO.Push(value)
	PTypes.Push(tipo)
}

func PushOp(op string) {
	fmt.Printf("→ PUSH OPERADOR: %s\n", op)
	POper.Push(op)
}

// Función auxiliar para generar temporales
func newTemp() string {
	tempVar++
	return fmt.Sprintf("t%d", tempVar)
}

// PushQuad agrega un cuadruplo a la lista
func PushQuad(oper string, left, right, res interface{}) {
	Quads = append(Quads, QuadStructure{oper, left, right, res})
}

func PrintStacks() {
	fmt.Println("Operandos:", PilaO.items)
	fmt.Println("Tipos:", PTypes.items)
	fmt.Print("Operadores: [")
	for i := len(POper.items) - 1; i >= 0; i-- {
		fmt.Print(POper.items[i], " ")
	}
	fmt.Println("]")
}

// PrintQuads imprime los cuadruplos con formato limpio
func PrintQuads() {
	fmt.Println("Cuádruplos generados:")
	for i, q := range Quads {
		fmt.Printf("%d: (%s %v %v %v)\n", i+1, q.Oper, q.Left, q.Right, q.Result)
	}
}

func DoAddSub() error {
	fmt.Printf("→ ADENTRO DOADDSUB ANTES DEL FOR \n")
	for {
		fmt.Printf("→ ADENTRO DOADDSUB ADENTRO DEL FOR \n")

		top, err := POper.Peek()

		fmt.Printf("→ TOPE POper: %v (%T)\n", top, top)

		if err != nil {
			fmt.Printf("→ PILA VACIA \n")
			return nil // pila vacía
		}

		op := top.(string)
		if op != "+" && op != "-" {
			fmt.Printf("→ NO OPERADOR DE SUMA/RESTA \n")
			break // no es operador de suma/resta
		}

		// Sacar operandos y tipos
		rightOp, _ := PilaO.Pop()
		rightType, _ := PTypes.Pop()
		fmt.Printf("→ rightOp: %v, rightType: %v\n", rightOp, rightType)

		leftOp, _ := PilaO.Pop()
		leftType, _ := PTypes.Pop()
		fmt.Printf("→ leftOp: %v, leftType: %v\n", leftOp, leftType)

		ltype, ok1 := leftType.(string)
		rtype, ok2 := rightType.(string)
		if !ok1 || !ok2 {
			return fmt.Errorf("DoAddSub error: tipos no son string: left=%T, right=%T", leftType, rightType)
		}

		POper.Pop()

		resType, err := GetResultType(ltype, rtype, op)
		if err != nil {
			fmt.Printf("→ ADENTRO DE IF ERR \n")
			return err
		}

		fmt.Printf("→ ANTES DE GENERATE QUAD \n")

		temp := newTemp()
		PushQuad(op, leftOp, rightOp, temp)

		PilaO.Push(temp)
		PTypes.Push(resType)

		fmt.Printf("→ GENERATE QUAD: %s %v %v -> %v\n", op, leftOp, rightOp, temp)
	}
	return nil
}

func DoMulDiv() error {
	for {
		top, err := POper.Peek()
		if err != nil {
			return nil // pila vacía
		}

		op := top.(string)
		if op != "*" && op != "/" {
			break // no es operador de mul/div
		}

		rightOp, _ := PilaO.Pop()
		rightType, _ := PTypes.Pop()
		leftOp, _ := PilaO.Pop()
		leftType, _ := PTypes.Pop()

		ltype, ok1 := leftType.(string)
		rtype, ok2 := rightType.(string)
		if !ok1 || !ok2 {
			return fmt.Errorf("DoMulDiv error: tipos no son string: left=%T, right=%T", leftType, rightType)
		}

		POper.Pop()

		resType, err := GetResultType(ltype, rtype, op)
		if err != nil {
			return err
		}

		temp := newTemp()
		PushQuad(op, leftOp, rightOp, temp)

		PilaO.Push(temp)
		PTypes.Push(resType)

		fmt.Printf("→ GENERATE QUAD: %s %v %v -> %v\n", op, leftOp, rightOp, temp)
	}
	return nil
}

// ------------------------------------ Vars

// Reset reinicia el scope global y limpia la pila de scopes locales.
func Reset() {
	scopes = &ScopeManager{
		global: NewDictionary(),
		stack:  []*Dictionary{},
	}
}

// ResetSemanticState limpia el directorio de funciones y el scope global.
func ResetSemanticState() {
	// 1) reinicia el scope manager
	Reset() // de scope.go
	// 2) reinicia el FunctionDirectory
	FunctionDirectory = NewDictionary()
}

// RegisterMainProgram crea el scope global y registra el programa principal.
func RegisterMainProgram(programName string) error {
	// 1) volver a un estado limpio
	// ResetSemanticState()

	// 2) error si ya existe el programa
	if _, exists := FunctionDirectory.Get(programName); exists {
		return fmt.Errorf("error: el programa '%s' ya ha sido declarado", programName)
	}

	// 3) registra el programa como función void, usando la tabla global actual
	FunctionDirectory.Put(programName, FunctionStructure{
		Name:     programName,
		VarTable: Current(), // aquí guardamos el scope global
	})

	fmt.Printf("Programa principal '%s' registrado exitosamente.\n", programName)

	/*fmt.Printf("[DEBUG] RegisterMainProgram %s → global scope %p\n",
	programName, Current())*/
	return nil
}

// Función para procesar la declaración de variables en el scope actual
func VarDeclaration(ids []string, tipo string) error {
	tabla := Current() // usa la tabla activa de scope.go
	/*fmt.Printf("[DEBUG] VarDeclaration %v in scope %p, parent=%p, before=%v\n",
	ids, tabla, tabla.parent, tabla.Keys())*/

	for _, id := range ids {
		if _, exists := tabla.Get(id); exists {
			return fmt.Errorf("error: variable '%s' ya declarada en este ámbito", id)
		}
		tabla.Put(id, VariableStructure{Name: id, Type: tipo})
	}

	/*
		// ——— Imprimimos el contenido del scope actual ———
		fmt.Println(">>> Contenido del scope actual antes de declarar:")
		tabla.PrintOrdered()
		fmt.Println(">>> Fin del scope actual")
	*/

	/*fmt.Printf("[DEBUG] VarDeclaration %v done in scope %p, after=%v\n",
	ids, tabla, tabla.Keys())*/

	return nil
}

// Verifica que los parámetros de una función no estén duplicados
func ValidateParams(params []VariableStructure) error {
	paramSet := NewDictionary()
	for _, param := range params {
		if _, exists := paramSet.Get(param.Name); exists {
			return fmt.Errorf("error: parámetro '%s' duplicado en la función", param.Name)
		}
		paramSet.Put(param.Name, param)
	}
	return nil
}

// FuncDeclaration actualiza la entrada creada por RegisterFunction
func FuncDeclaration(name string, params []VariableStructure) error {
	if err := ValidateParams(params); err != nil {
		return err
	}
	raw, exists := FunctionDirectory.Get(name)
	if !exists {
		return fmt.Errorf("error interno: función '%s' no registrada previamente", name)
	}
	fs := raw.(FunctionStructure)
	fs.Parameters = params
	fs.VarTable = Current() // la tabla local donde metimos los params
	FunctionDirectory.Put(name, fs)
	return nil
}

// RegisterFunction crea la entrada de la función con nombre, retorno void,
// sin parámetros todavía, y le asigna una tabla de variables local vacía.
func RegisterFunction(name string) error {
	if _, exists := FunctionDirectory.Get(name); exists {
		return fmt.Errorf("error: función '%s' ya declarada", name)
	}

	// 1) creamos la tabla local (de momento vacía)
	localTable := NewDictionary()

	/*fmt.Printf("[DEBUG] RegisterFunction %s → local scope %p\n",
	name, localTable)*/

	// 2) registramos la función en el directorio
	FunctionDirectory.Put(name, FunctionStructure{
		Name:       name,
		Parameters: []VariableStructure{}, // params vacía por ahora
		VarTable:   localTable,
	})

	return nil
}
