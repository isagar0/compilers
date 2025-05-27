// vm.go
package semantics

import (
	"strconv"
	"strings"
)

type VirtualMemory struct {
	Storage map[int]interface{} // Almacena valores por dirección
}

func NewVirtualMemory() *VirtualMemory {
	return &VirtualMemory{
		Storage: make(map[int]interface{}),
	}
}

// Inicializa memoria con constantes y variables globales
func (vm *VirtualMemory) Initialize() {
	// Cargar constantes
	for addr, name := range AddressToName {
		if strings.HasPrefix(name, "const_") {
			valueStr := strings.TrimPrefix(name, "const_")
			if intVal, err := strconv.Atoi(valueStr); err == nil {
				vm.Storage[addr] = intVal
			}
		}
	}

	// Inicializar variables globales a 0
	for _, v := range VarTable.Items {
		vs := v.(VariableStructure)
		vm.Storage[vs.Address] = 0
	}
}

func (vm *VirtualMemory) Get(addr int) interface{} {
	return vm.Storage[addr]
}

func (vm *VirtualMemory) Set(addr int, value interface{}) {
	vm.Storage[addr] = value
}

// vm.go (continuación)
func ExecuteQuads(quads []QuadStructure) {
	vm := NewVirtualMemory()
	vm.Initialize()

	for ip := 0; ip < len(quads); ip++ {
		quad := quads[ip]
		switch quad.Oper {
		case "+":
			left := vm.Get(quad.Left.(int)).(int)
			right := vm.Get(quad.Right.(int)).(int)
			resultAddr := quad.Result.(int)
			vm.Set(resultAddr, left+right)

		case "=":
			source := vm.Get(quad.Left.(int)).(int)
			destAddr := quad.Result.(int)
			vm.Set(destAddr, source)

		case "PRINT":
			value := vm.Get(quad.Left.(int)).(int)
			println(value) // O usa fmt.Print para mejor formato
		}
	}
}
