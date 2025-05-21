package semantics

import "fmt"

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
