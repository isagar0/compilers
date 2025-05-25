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

// ProcessOperation: Agregar quads
func ProcessOperation(validOps []string, stopOnFakeBottom bool) error {
	for {
		top, err := POper.Peek()
		if err != nil {
			return nil
		}

		op := top.(string)

		// Caso especial: fake bottom
		if stopOnFakeBottom && op == "⏊" {
			POper.Pop()
			break
		}

		// Verificar operadores válidos
		valid := false
		for _, validOp := range validOps {
			if op == validOp {
				valid = true
				break
			}
		}
		if !valid {
			if stopOnFakeBottom {
				return fmt.Errorf("operador inesperado en paréntesis: %v", op)
			}
			break
		}

		// Lógica común de procesamiento
		POper.Pop()
		rightOp, _ := PilaO.Pop()
		rightType, _ := PTypes.Pop()
		leftOp, _ := PilaO.Pop()
		leftType, _ := PTypes.Pop()

		ltype, ok1 := leftType.(string)
		rtype, ok2 := rightType.(string)
		if !ok1 || !ok2 {
			return fmt.Errorf("error: tipos no son string: left=%T, right=%T", leftType, rightType)
		}

		resType, err := GetResultType(ltype, rtype, op)
		if err != nil {
			return err
		}

		var tempAddr int
		switch resType {
		case "int":
			tempAddr, _ = memory.Temp.Ints.GetNext()
		case "float":
			tempAddr, _ = memory.Temp.Floats.GetNext()
		case "bool":
			tempAddr, _ = memory.Temp.Bools.GetNext()
		}

		AddressToName[tempAddr] = fmt.Sprintf("temp_%d", tempAddr)

		PushQuad(op, leftOp, rightOp, tempAddr)
		PilaO.Push(tempAddr)
		PTypes.Push(resType)
	}
	return nil
}

// DoAddSub: Agregar quad para suma o resta
func DoAddSub() error {
	return ProcessOperation([]string{"+", "-"}, false)
}

// DoMulDiv: Agregar quad para multiplicación o divición
func DoMulDiv() error {
	return ProcessOperation([]string{"*", "/"}, false)
}

// DoRelational: Agregar quad para operadores relacionales
func DoRelational() error {
	if top, err := POper.Peek(); err != nil || !(top.(string) == "<" || top.(string) == ">" || top.(string) == "!=") {
		return nil
	}
	return ProcessOperation([]string{"<", ">", "!="}, false)
}

// PopUntilFakeBottom: Agregar simbolo para parentesis
func PopUntilFakeBottom() error {
	return ProcessOperation([]string{"+", "-", "*", "/"}, true)
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

// ResetStacks: Limpia todas las variables de los quads
func ResetStacks() {
	PilaO = NewStack()
	PTypes = NewStack()
	POper = NewStack()
	Quads = []QuadStructure{}
	tempVar = 0
}
