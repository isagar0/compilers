package main

import "fmt"

// Stack: Last In, Frist Out
// Almacena una lista de enteros
type Stack struct {
	items []int
}

// Push: Agrega elemento arriba de la pila
func (s *Stack) Push(item int) {
	s.items = append(s.items, item)
}

// Pop: Elimina y regresa último elemento
func (s *Stack) Pop() (int, error) {
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
func (s *Stack) Peek() (int, error) {
	// Si la pila esta vacía, retorna error
	if len(s.items) == 0 {
		return 0, fmt.Errorf("stack is empty")
	}

	return s.items[len(s.items)-1], nil
}

// IsEmpty: Verifica si la pila esta vacía
func (s *Stack) IsEmpty() bool {
	// Regresa true si esta vacía
	return len(s.items) == 0
}

// Size: Regresa cantidad de elementos en la pila
func (s *Stack) Size() int {
	return len(s.items)
}
