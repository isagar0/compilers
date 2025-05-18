package semantics

import "strconv"

type ConstTable struct {
	ints    []string
	floats  []string
	strings []string // Add string support
}

var Consts = ConstTable{
	ints:    []string{},
	floats:  []string{},
	strings: []string{}, // Initialize string storage
}

func GetConstAddress(literal string, tipo string) int {
	var list *[]string
	var segment *MemorySegment
	var addr int
	var err error

	switch tipo {
	case "int":
		list = &Consts.ints
		segment = &memory.Constant.Ints
	case "float":
		list = &Consts.floats
		segment = &memory.Constant.Floats
	case "string": // Add string case
		list = &Consts.strings
		segment = &memory.Constant.Strings // You'll need to add this to MemoryManager
	default:
		panic("Tipo de constante no soportado: " + tipo)
	}

	// Buscar si ya existe
	for i, val := range *list {
		if val == literal {
			return segment.start + i
		}
	}

	// Si no existe, asigna nueva dirección
	addr, err = segment.GetNext()

	if err != nil {
		panic(err)
	}

	// Evita registrar direcciones virtuales como constantes
	if num, err := strconv.Atoi(literal); err == nil {
		if num >= 1000 && num <= 6999 { // dentro de rangos de variables o temporales
			return num // ya es una dirección, no registrar
		}
	}

	*list = append(*list, literal)
	return addr
}
