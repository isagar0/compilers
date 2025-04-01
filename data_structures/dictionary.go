package main

import "fmt"

// Dictionary: Almacena pares clave-valor (string-int)
// Utiliza keys para matener el orden de inserción de las claves
type Dictionary struct {
	items map[string]int // Mapa
	keys  []string       // Slice para ordenar
}

// NewDictionary: Constructor del diccionario
func NewDictionary() *Dictionary {
	return &Dictionary{
		items: make(map[string]int),
		keys:  []string{},
	}
}

// Put: Inserta una clave-valor manteniendo el orden de inserción
func (d *Dictionary) Put(key string, value int) {
	// Si la clave no existe, se agrega a la lista
	if _, exists := d.items[key]; !exists {
		d.keys = append(d.keys, key) // Agrega la clave al final
	}
	d.items[key] = value // Asigna el valor en el mapa
}

// Get: Obtiene el valor asociado a una clave
func (d *Dictionary) Get(key string) (int, bool) {
	value, exists := d.items[key]
	return value, exists
}

// Remove: Elimina una clave-valor y actualiza la lista de claves
func (d *Dictionary) Remove(key string) {
	// Si la clave existe en el diccionario, se elimina
	if _, exists := d.items[key]; exists {
		delete(d.items, key) // Elimina la clave del mapa

		// Busca y elimina la clave del slice
		for i, k := range d.keys {
			if k == key {
				d.keys = append(d.keys[:i], d.keys[i+1:]...) // Quita la clave del slice
				break                                        // Termina cuando encuentra la clave
			}
		}
	}
}

// IsEmpty: Verifica si el diccionario está vacío
func (d *Dictionary) IsEmpty() bool {
	// Regresa true si no hay elementos en el mapa
	return len(d.items) == 0
}

// Size: Regresa la cantidad de elementos en el diccionario
func (d *Dictionary) Size() int {
	return len(d.items)
}

// Keys: Regresa las claves en orden de inserción
func (d *Dictionary) Keys() []string {
	return d.keys
}

// PrintOrdered: Imprime el diccionario en orden de inserción
func (d *Dictionary) PrintOrdered() {
	// Itera sobre las claves en orden y muestra su valor
	for _, key := range d.keys {
		fmt.Printf("%s: %d\n", key, d.items[key])
	}
}
