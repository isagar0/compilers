package semantics

import (
	"fmt"
	"strconv"
)

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
		},
		Constant: SegmentGroup{
			Ints:    NewMemorySegment(7000, 7999),
			Floats:  NewMemorySegment(8000, 8999),
			Strings: NewMemorySegment(9000, 9999), // Add string constant range
		},
	}
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
