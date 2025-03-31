package main

import "fmt"

// Queue: Firt in, First Out
type Queue struct {
	items []int
}

// Enqueue: Agrega elemento al final de la cola
func (q *Queue) Enqueue(item int) {
	q.items = append(q.items, item)
}

// Dequeue: Elimina y regresa el primer elemento
func (q *Queue) Dequeue() (int, error) {
	// Si la pila esta vacía, retorna error
	if len(q.items) == 0 {
		return 0, fmt.Errorf("queue is empty")
	}

	item := q.items[0]    // Guarda elemento
	q.items = q.items[1:] // Elimina elemento
	return item, nil
}

// Peek: Regresa primer elemento
func (q *Queue) Peek() (int, error) {
	// Si la pila esta vacía, retorna error
	if len(q.items) == 0 {
		return 0, fmt.Errorf("queue is empty")
	}

	return q.items[0], nil
}

// IsEmpty: Verifica si la cola esta vacía
func (q *Queue) IsEmpty() bool {
	return len(q.items) == 0
}

// Size: Regresa cantidad de elementos en la cola
func (q *Queue) Size() int {
	return len(q.items)
}

// Print: Imprime los elementos de la cola en orden
func (q *Queue) Print() {
	for _, item := range q.items {
		fmt.Print(item, " ")
	}
	fmt.Println()
}
