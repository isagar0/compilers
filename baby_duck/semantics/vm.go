package semantics

import (
	"fmt"
	"strconv"
	"strings"
)

type ActivationRecord struct {
	ReturnIP int
	LocalMem map[int]interface{}
}

type VirtualMachine struct {
	Quads        []QuadStructure
	GlobalMemory map[int]interface{} // Memoria global (variables + constantes)
	LocalMemory  map[int]interface{} // Memoria local (función actual)
	IP           int
	CallStack    []ActivationRecord           // Pila de llamadas
	FuncDir      map[string]FunctionStructure // Directorio de funciones
	PendingAR    map[int]interface{}          // Registro de activación pendiente (para ERA)
}

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

// InitializeMemory Inicializa la memoria con constantes y variables globales
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

func (vm *VirtualMachine) readMem(addr int) interface{} {
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

func (vm *VirtualMachine) writeMem(addr int, value interface{}) {
	if (addr >= 1000 && addr <= 2999) || (addr >= 8000) {
		vm.GlobalMemory[addr] = value
	} else {
		vm.LocalMemory[addr] = value
	}
}

// ExecuteNext Ejecuta el siguiente cuádruplo y devuelve false si terminó
func (vm *VirtualMachine) ExecuteNext() bool {
	if vm.IP >= len(vm.Quads) {
		return false
	}

	quad := vm.Quads[vm.IP]
	vm.IP++

	switch quad.Oper {
	case "+", "-", "*", "/":
		left := vm.readMem(quad.Left.(int))
		right := vm.readMem(quad.Right.(int))
		resultAddr := quad.Result.(int)

		switch quad.Oper {
		case "+":
			vm.writeMem(resultAddr, add(left, right))
		case "-":
			vm.writeMem(resultAddr, sub(left, right))
		case "*":
			vm.writeMem(resultAddr, mul(left, right))
		case "/":
			vm.writeMem(resultAddr, div(left, right))
		}

	case "<", ">", "!=", "==":
		left := vm.readMem(quad.Left.(int))
		right := vm.readMem(quad.Right.(int))
		resultAddr := quad.Result.(int)

		leftVal := toFloat(left)
		rightVal := toFloat(right)

		var result bool
		switch quad.Oper {
		case "<":
			result = leftVal < rightVal
		case ">":
			result = leftVal > rightVal
		case "!=":
			result = leftVal != rightVal
		case "==":
			result = leftVal == rightVal
		}

		vm.writeMem(resultAddr, result)

	case "=":
		source := vm.readMem(quad.Left.(int))
		destAddr := quad.Result.(int)
		vm.writeMem(destAddr, source)

	case "PRINT":
		addr := quad.Left.(int)
		if addr >= 10000 && addr <= 10999 {
			name, exists := AddressToName[addr]
			if !exists {
				panic("String no encontrada en AddressToName")
			}
			value := strings.TrimPrefix(name, "const_")
			fmt.Println(value)
		} else {
			value := vm.readMem(addr)
			fmt.Println(value)
		}

	case "GOTOF":
		conditionAddr := quad.Left.(int)
		condition := vm.readMem(conditionAddr).(bool)
		if !condition {
			target := quad.Result.(int)
			vm.IP = target
		}

	case "GOTO":
		target := quad.Result.(int)
		vm.IP = target

	case "ERA":
		vm.PendingAR = make(map[int]interface{})

	case "PARAMETER":
		srcAddr := quad.Left.(int)
		paramIndex := quad.Result.(int)
		value := vm.readMem(srcAddr)
		vm.PendingAR[paramIndex] = value

	case "GOSUB":
		funcName := quad.Left.(string)
		funcData := vm.FuncDir[funcName]

		vm.CallStack = append(vm.CallStack, ActivationRecord{
			ReturnIP: vm.IP,
			LocalMem: vm.LocalMemory,
		})

		vm.LocalMemory = vm.PendingAR
		vm.IP = funcData.StartQuad
		vm.PendingAR = nil

	case "ENDFUNC":
		if len(vm.CallStack) == 0 {
			panic("ENDFUNC sin llamada activa")
		}
		frame := vm.CallStack[len(vm.CallStack)-1]
		vm.CallStack = vm.CallStack[:len(vm.CallStack)-1]
		vm.LocalMemory = frame.LocalMem
		vm.IP = frame.ReturnIP

	case "END": // Manejar fin del programa
		return false

	default:
		panic("Operación no soportada: " + quad.Oper)
	}

	return true
}

// Run Ejecuta todos los cuádruplos hasta terminar
func (vm *VirtualMachine) Run() {
	vm.InitializeMemory()
	vm.LocalMemory = make(map[int]interface{}) // Memoria para main
	for vm.ExecuteNext() {
	}
}

// ---- Funciones auxiliares para operaciones aritméticas ----
func toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	default:
		panic("Tipo no soportado para conversión a float")
	}
}

func add(a, b interface{}) interface{} {
	// Si al menos uno es float, convertir ambos
	if _, isFloatA := a.(float64); isFloatA {
		return toFloat(a) + toFloat(b)
	}
	if _, isFloatB := b.(float64); isFloatB {
		return toFloat(a) + toFloat(b)
	}
	// Ambos son int
	return a.(int) + b.(int)
}

func sub(a, b interface{}) interface{} {
	if _, isFloatA := a.(float64); isFloatA {
		return toFloat(a) - toFloat(b)
	}
	if _, isFloatB := b.(float64); isFloatB {
		return toFloat(a) - toFloat(b)
	}
	return a.(int) - b.(int)
}

func mul(a, b interface{}) interface{} {
	if _, isFloatA := a.(float64); isFloatA {
		return toFloat(a) * toFloat(b)
	}
	if _, isFloatB := b.(float64); isFloatB {
		return toFloat(a) * toFloat(b)
	}
	return a.(int) * b.(int)
}

func div(a, b interface{}) interface{} {
	// La división siempre devuelve float
	return toFloat(a) / toFloat(b)
}
