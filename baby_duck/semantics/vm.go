package semantics

import (
	"fmt"
	"strconv"
	"strings"
)

// -------------------------------------------- VM --------------------------------------------

// NewVirtualMachine: Inicializa la VM
func NewVirtualMachine(quads []QuadStructure, funcDir *Dictionary) *VirtualMachine {
	// Convertir Dictionary a map[string]FunctionStructure
	funcDirMap := make(map[string]FunctionStructure)
	for key, value := range funcDir.Items {
		if fs, ok := value.(FunctionStructure); ok {
			funcDirMap[key] = fs
		}
	}

	return &VirtualMachine{
		Quads:        quads,
		GlobalMemory: make(map[int]interface{}),
		LocalMemory:  make(map[int]interface{}),
		IP:           0,
		CallStack:    make([]ActivationRecord, 0),
		FuncDir:      funcDirMap,
		PendingAR:    nil,
	}
}

// InitializeMemory: Inicializa la memoria con constantes y variables globales
func (vm *VirtualMachine) InitializeMemory() {
	// Cargar constantes desde AddressToName
	for addr, name := range AddressToName {
		if strings.HasPrefix(name, "const_") {
			valueStr := strings.TrimPrefix(name, "const_")
			// Intentar convertir a int
			if intVal, err := strconv.Atoi(valueStr); err == nil {
				vm.GlobalMemory[addr] = intVal
				continue
			}

			// Intentar convertir a float
			if floatVal, err := strconv.ParseFloat(valueStr, 64); err == nil {
				vm.GlobalMemory[addr] = floatVal
			}
		}
	}

	// Inicializar variables globales (según VarTable)
	for _, v := range VarTable.Items {
		vs := v.(VariableStructure)
		if vs.Address < 3000 { // Solo globales
			switch vs.Type {
			case "int":
				vm.GlobalMemory[vs.Address] = 0
			case "float":
				vm.GlobalMemory[vs.Address] = 0.0
			case "bool":
				vm.GlobalMemory[vs.Address] = false
			}
		}
	}
}

// ReadMem: Retorna valor almacenado en una dirección de memoria
func (vm *VirtualMachine) ReadMem(addr int) interface{} {
	// Rangos de memoria global/constantes
	if (addr >= 1000 && addr <= 2999) || (addr >= 8000) {
		val, ok := vm.GlobalMemory[addr]
		if !ok {
			// Inicializar según tipo de dirección
			if addr >= 1000 && addr <= 1999 {
				val = 0 // int
			} else if addr >= 2000 && addr <= 2999 {
				val = 0.0 // float
			} else if addr >= 8000 && addr <= 8999 {
				val = 0 // int
			} else if addr >= 9000 && addr <= 9999 {
				val = 0.0 // float
			} else if addr >= 10000 && addr <= 10999 {
				val = "" // string
			}
			vm.GlobalMemory[addr] = val
		}
		return val
	}

	// Memoria local/temporal
	val, ok := vm.LocalMemory[addr]
	if !ok {
		// Inicializar según tipo de dirección
		if addr >= 3000 && addr <= 3999 {
			val = 0 // int
		} else if addr >= 4000 && addr <= 4999 {
			val = 0.0 // float
		} else if addr >= 5000 && addr <= 5999 {
			val = 0 // int
		} else if addr >= 6000 && addr <= 6999 {
			val = 0.0 // float
		} else if addr >= 7000 && addr <= 7999 {
			val = false // bool
		}
		vm.LocalMemory[addr] = val
	}
	return val
}

// WriteMem: Escribe valor en memorial global o local
func (vm *VirtualMachine) WriteMem(addr int, value interface{}) {
	if (addr >= 1000 && addr <= 2999) || (addr >= 8000) {
		vm.GlobalMemory[addr] = value
	} else {
		vm.LocalMemory[addr] = value
	}
}

// LogMemory: Imprime estado actual de la memorias
func (vm *VirtualMachine) LogMemory() {
	/*
		fmt.Println("--- Memory State ---")
		for addr, val := range vm.GlobalMemory {
			fmt.Printf("Global [%d]: %v\n", addr, val)
		}
		for addr, val := range vm.LocalMemory {
			fmt.Printf("Local [%d]: %v\n", addr, val)
		}
	*/
}

// -------------------------------------------- EXECUTE --------------------------------------------

// ExecuteNext: Ejecuta el siguiente cuádruplo y devuelve false si terminó
func (vm *VirtualMachine) ExecuteNext() bool {
	// Verifica si IP se paso del numero de quadruplos
	if vm.IP >= len(vm.Quads) {
		return false
	}

	// Obtiene cuadruplo actual y avanza
	quad := vm.Quads[vm.IP]
	vm.IP++
	//fmt.Printf("Executing %d: %s %v %v %v\n", vm.IP, quad.Oper, quad.Left, quad.Right, quad.Result)

	// Determina tipo de operación a ejecutar
	switch FixedAddresses[quad.Oper] {
	case "+", "-", "*", "/":
		// Lee operando izquierdo y derecho desde memoria
		left := vm.ReadMem(quad.Left.(int))
		right := vm.ReadMem(quad.Right.(int))
		resultAddr := quad.Result.(int)

		// Ejecuta la operación
		switch FixedAddresses[quad.Oper] {
		case "+":
			vm.WriteMem(resultAddr, Add(left, right))
		case "-":
			vm.WriteMem(resultAddr, Sub(left, right))
		case "*":
			vm.WriteMem(resultAddr, Mul(left, right))
		case "/":
			vm.WriteMem(resultAddr, Div(left, right))
		}

	case "<", ">", "!=", "==":
		// Lee operando izquierdo y derecho desde memoria
		left := vm.ReadMem(quad.Left.(int))
		right := vm.ReadMem(quad.Right.(int))
		resultAddr := quad.Result.(int)

		// Convertir para comparación
		leftVal := ToFloat(left)
		rightVal := ToFloat(right)

		// Ejecuta la evaluacion
		var result bool
		switch FixedAddresses[quad.Oper] {
		case "<":
			result = leftVal < rightVal
		case ">":
			result = leftVal > rightVal
		case "!=":
			result = leftVal != rightVal
		case "==":
			result = leftVal == rightVal
		}

		vm.WriteMem(resultAddr, result)

	case "=":
		// Copia el valor de una dirección a otra
		source := vm.ReadMem(quad.Left.(int))
		destAddr := quad.Result.(int)
		vm.WriteMem(destAddr, source)

	case "PRINT":
		addr := quad.Left.(int)
		// Si es string constante, obtiene su valor
		if addr >= 10000 && addr <= 10999 {
			name, exists := AddressToName[addr]
			if !exists {
				panic("String no encontrada en AddressToName")
			}

			value := strings.TrimPrefix(name, "const_")
			value = strings.Trim(value, "\"")

			fmt.Println(value)
		} else {
			// Si no es string lee el valor de memoria
			value := vm.ReadMem(addr)
			fmt.Println(value)
		}

	case "GOTOF":
		conditionAddr := quad.Left.(int)
		condition := vm.ReadMem(conditionAddr).(bool)

		if !condition {
			// Salta al cuádruplo indicado si la conficion es falsa
			target := quad.Result.(int)
			vm.IP = target
		}

	case "GOTO":
		// Salta al cuadruplo destino imediatamente
		target := quad.Result.(int)
		vm.IP = target

	case "ERA":
		// Inicializa ActivationRecord para parametros
		vm.PendingAR = make(map[int]interface{})

	case "PARAMETER":
		// Asigna valor a un parametro en ActivationRecord temporal
		srcAddr := quad.Left.(int)
		paramIndex := quad.Result.(int)
		value := vm.ReadMem(srcAddr)
		vm.PendingAR[paramIndex] = value

	case "GOSUB":
		funcName := quad.Left.(string)
		funcData := vm.FuncDir[funcName]

		// Crear memoria local usando los índices de parámetros
		newLocal := make(map[int]interface{})
		for paramIndex, paramValue := range vm.PendingAR {
			// Asegura indice sea valido dentro del arreglo parametro
			if paramIndex-1 < len(funcData.Parameters) && paramIndex-1 >= 0 {
				paramAddr := funcData.Parameters[paramIndex-1].Address + 1 // Usa dirección del parámetro
				newLocal[paramAddr] = paramValue

			} else {
				panic(fmt.Sprintf("Índice inválido: %d (parámetros: %d)", paramIndex, len(funcData.Parameters)))
			}
		}

		// Guarda el estado actual pila llamadas
		vm.CallStack = append(vm.CallStack, ActivationRecord{
			ReturnIP: vm.IP,
			LocalMem: vm.LocalMemory,
		})

		// Cambia la memoria local y salta cuadruplo de inicio
		vm.LocalMemory = newLocal
		vm.IP = funcData.StartQuad
		vm.PendingAR = nil

	case "ENDFUNC":
		if len(vm.CallStack) == 0 {
			panic("ENDFUNC sin llamada activa")
		}

		// Finaliza función actual restaura estado anterior
		frame := vm.CallStack[len(vm.CallStack)-1]
		vm.CallStack = vm.CallStack[:len(vm.CallStack)-1]
		vm.LocalMemory = frame.LocalMem
		vm.IP = frame.ReturnIP

	case "END":
		// Marca fin del programa
		return false

	default:
		// Si es desconocida lanza error
		panic(fmt.Sprintf("Operación no soportada: %d", quad.Oper))
	}

	vm.LogMemory()

	return true
}

// Run: Ejecuta todos los cuádruplos hasta terminar
func (vm *VirtualMachine) Run() {
	vm.InitializeMemory()
	vm.LocalMemory = make(map[int]interface{}) // Memoria para main
	vm.LogMemory()
	for vm.ExecuteNext() {
	}
}

// -------------------------------------------- FUN --------------------------------------------

// ToFloat: Convierte int o float a float64
func ToFloat(val interface{}) float64 {
	switch v := val.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	default:
		panic("Tipo no soportado para conversión a float")
	}
}

// Add: Realiza operación suma
func Add(a, b interface{}) interface{} {
	// Si al menos uno es float, convertir ambos
	if _, isFloatA := a.(float64); isFloatA {
		return ToFloat(a) + ToFloat(b)
	}
	if _, isFloatB := b.(float64); isFloatB {
		return ToFloat(a) + ToFloat(b)
	}
	// Ambos son int
	return a.(int) + b.(int)
}

// Sub: Realiza operación resta
func Sub(a, b interface{}) interface{} {
	if _, isFloatA := a.(float64); isFloatA {
		return ToFloat(a) - ToFloat(b)
	}
	if _, isFloatB := b.(float64); isFloatB {
		return ToFloat(a) - ToFloat(b)
	}
	return a.(int) - b.(int)
}

// Mul: Realiza operación multiplicación
func Mul(a, b interface{}) interface{} {
	if _, isFloatA := a.(float64); isFloatA {
		return ToFloat(a) * ToFloat(b)
	}
	if _, isFloatB := b.(float64); isFloatB {
		return ToFloat(a) * ToFloat(b)
	}
	return a.(int) * b.(int)
}

// Div: Realiza operación divición
func Div(a, b interface{}) interface{} {
	// La división siempre devuelve float
	return ToFloat(a) / ToFloat(b)
}
