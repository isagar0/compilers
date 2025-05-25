package semantics

import (
	"fmt"
	"strings"
)

// --------------------------------------- DICTIONARY ---------------------------------------
// NewDictionary: Constructor del diccionario
func NewDictionary() *Dictionary {
	return &Dictionary{
		Items:  make(map[string]interface{}),
		parent: nil,
	}
}

// Put: Inserta una clave-valor manteniendo el orden de inserción
func (d *Dictionary) Put(key string, value interface{}) {
	d.Items[key] = value
}

// Get: Obtiene el valor asociado a una clave
func (d *Dictionary) Get(key string) (interface{}, bool) {
	if val, ok := d.Items[key]; ok {
		return val, true
	}
	if d.parent != nil {
		return d.parent.Get(key)
	}
	return nil, false
}

// ----------------------------------------- MEMORY -----------------------------------------

func NewMemorySegment(start, end int) MemorySegment {
	return MemorySegment{start: start, end: end, next: start}
}

func (m *MemorySegment) GetNext() (int, error) {
	if m.next > m.end {
		return -1, fmt.Errorf("memoria llena en el rango %d–%d", m.start, m.end)
	}
	addr := m.next
	m.next++
	return addr, nil
}

func (m *MemorySegment) Reset() {
	m.next = m.start
}

func NewMemoryManager() *MemoryManager {
	return &MemoryManager{
		Global: SegmentGroup{
			Ints:   NewMemorySegment(1000, 1999),
			Floats: NewMemorySegment(2000, 2999),
		},
		Local: SegmentGroup{
			Ints:   NewMemorySegment(3000, 3999),
			Floats: NewMemorySegment(4000, 4999),
		},
		Temp: SegmentGroup{
			Ints:   NewMemorySegment(5000, 5999),
			Floats: NewMemorySegment(6000, 6999),
			Bools:  NewMemorySegment(7000, 7999),
		},
		Constant: SegmentGroup{
			Ints:    NewMemorySegment(8000, 8999),
			Floats:  NewMemorySegment(9000, 9999),
			Strings: NewMemorySegment(10000, 10999),
		},
	}
}

func GetConstAddress(literal string, tipo string) int {
	var segment *MemorySegment

	switch tipo {
	case "int":
		segment = &memory.Constant.Ints
	case "float":
		segment = &memory.Constant.Floats
	case "string":
		segment = &memory.Constant.Strings
	default:
		panic("Tipo de constante no soportado: " + tipo)
	}

	// Buscar si ya existe la constante en el segmento
	for i := 0; i < segment.next-segment.start; i++ {
		addr := segment.start + i
		if AddressToName[addr] == "const_"+literal {
			return addr
		}
	}

	// Si no existe, asignar nueva dirección
	addr, err := segment.GetNext()
	if err != nil {
		panic(err)
	}

	// Registrar en AddressToName
	AddressToName[addr] = "const_" + literal

	return addr
}

// Tabla de direcciones
func PrintAddressTable() {
	fmt.Println("\n==== Tabla de direcciones virtuales ====")

	// First print variables from scopes
	fmt.Println("\n---- Variables Globales ----")
	printScopeVariables(scopes.global)
	for _, scope := range scopes.stack {
		printScopeVariables(scope)
	}

	// Then print constants
	fmt.Println("\n---- Constantes ----")
	for addr, name := range AddressToName {
		if strings.HasPrefix(name, "const_") && (addr < 1000 || addr >= 7000) {
			fmt.Printf("%-10s → %d\n", strings.TrimPrefix(name, "const_"), addr)
		}
	}

	// Then print temporaries
	fmt.Println("\n---- Temporales ----")
	for addr, name := range AddressToName {
		if strings.HasPrefix(name, "temp_") {
			fmt.Printf("%-10s → %d\n", name, addr)
		}
	}
}

func printScopeVariables(scope *Dictionary) {
	for _, val := range scope.Items {
		vs := val.(VariableStructure)
		fmt.Printf("%-10s → %-6d (%-6s)\n", vs.Name, vs.Address, vs.Type)
	}
}
