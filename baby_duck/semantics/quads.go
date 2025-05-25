package semantics

import (
	"fmt"
	"strconv"
)

// --------------------------------------- DECLARACION ---------------------------------------

var (
	PilaO   = NewStack()    // Operandos
	PTypes  = NewStack()    // Tipos de operadores
	POper   = NewStack()    // Operadores
	Quads   []QuadStructure // Cuádruplos generados
	tempVar int             // Contador para nombres de variables temporales
	PJumps  = NewStack()    // Stack para saltos pendientes
)

// ------------------------------------------ STACK ------------------------------------------

// NewStack: Crea una nueva pila vacía.
func NewStack() *Stack {
	return &Stack{
		items: make([]interface{}, 0, 16),
	}
}

// Push: Agrega elemento arriba de la pila
func (s *Stack) Push(v interface{}) {
	s.items = append(s.items, v)
}

// Pop: Elimina y regresa último elemento
func (s *Stack) Pop() (interface{}, error) {
	// Si la pila esta vacía, retorna error
	if len(s.items) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}

	lastIndex := len(s.items) - 1 // Obtiene elemento
	item := s.items[lastIndex]    // Guarda elemento
	s.items = s.items[:lastIndex] // Elimina elemento
	return item, nil
}

// Peek: Regresa último elemento
func (s *Stack) Peek() (interface{}, error) {
	// Si la pila esta vacía, retorna error
	if len(s.items) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}

	return s.items[len(s.items)-1], nil
}

// ------------------------------------------ QUADS ------------------------------------------
// PushOperandDebug: Push con debug
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

// CleanStacks: Limpia todas las variables de los quads
func CleanStacks() {
	PilaO = NewStack()
	PTypes = NewStack()
	POper = NewStack()
	Quads = []QuadStructure{}
	tempVar = 0
}
