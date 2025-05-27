package semantics

import (
	"fmt"
	"strconv"
	"strings"
)

// VirtualMachine representa el estado de ejecución
type VirtualMachine struct {
	Quads  []QuadStructure     // Lista de cuádruplos a ejecutar
	Memory map[int]interface{} // Memoria virtual (dirección -> valor)
	IP     int                 // Instruction Pointer (índice actual en Quads)
}

// NewVirtualMachine Crea una nueva máquina virtual con los cuádruplos generados
func NewVirtualMachine(quads []QuadStructure) *VirtualMachine {
	return &VirtualMachine{
		Quads:  quads,
		Memory: make(map[int]interface{}),
		IP:     0,
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
				vm.Memory[addr] = intVal
				continue
			}
			// Intentar convertir a float
			if floatVal, err := strconv.ParseFloat(valueStr, 64); err == nil {
				vm.Memory[addr] = floatVal
			}
		}
	}

	// Inicializar variables globales (según VarTable)
	for _, v := range VarTable.Items {
		vs := v.(VariableStructure)
		switch vs.Type {
		case "int":
			vm.Memory[vs.Address] = 0 // Valor por defecto para int
		case "float":
			vm.Memory[vs.Address] = 0.0 // Valor por defecto para float
		}
	}
}

// ExecuteNext Ejecuta el siguiente cuádruplo y devuelve false si terminó
func (vm *VirtualMachine) ExecuteNext() bool {
	if vm.IP >= len(vm.Quads) {
		return false
	}

	quad := vm.Quads[vm.IP]
	vm.IP++ // Avanzar al siguiente cuádruplo

	switch quad.Oper {
	case "+", "-", "*", "/":
		left := vm.Memory[quad.Left.(int)]
		right := vm.Memory[quad.Right.(int)]
		resultAddr := quad.Result.(int)

		switch quad.Oper {
		case "+":
			vm.Memory[resultAddr] = add(left, right)
		case "-":
			vm.Memory[resultAddr] = sub(left, right)
		case "*":
			vm.Memory[resultAddr] = mul(left, right)
		case "/":
			vm.Memory[resultAddr] = div(left, right)
		}

	case "=": // Asignación
		source := vm.Memory[quad.Left.(int)]
		destAddr := quad.Result.(int)
		vm.Memory[destAddr] = source

	case "PRINT": // Imprimir valor
		value := vm.Memory[quad.Left.(int)]
		fmt.Printf("%v\n", value) // Funciona para int y float

	default:
		panic("Operación no soportada: " + quad.Oper)
	}

	return true
}

// Run Ejecuta todos los cuádruplos hasta terminar
func (vm *VirtualMachine) Run() {
	vm.InitializeMemory()
	for vm.ExecuteNext() {
		// Continuar ejecución
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
