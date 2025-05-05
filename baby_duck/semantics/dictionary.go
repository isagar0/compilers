package semantics

import "fmt"

// Dictionary: Almacena pares clave-valor (string-valor de cualquier tipo)
// Utiliza keys para mantener el orden de inserción de las claves
type Dictionary struct {
	items map[string]interface{} // Mapa para almacenar cualquier tipo de valor
	keys  []string               // Slice para ordenar las claves
}

// NewDictionary: Constructor del diccionario
func NewDictionary() *Dictionary {
	return &Dictionary{
		items: make(map[string]interface{}),
		keys:  []string{},
	}
}

// Put: Inserta una clave-valor manteniendo el orden de inserción
func (d *Dictionary) Put(key string, value interface{}) {
	// Si la clave no existe, se agrega a la lista
	if _, exists := d.items[key]; !exists {
		d.keys = append(d.keys, key) // Agrega la clave al final
	}
	d.items[key] = value // Asigna el valor en el mapa
}

// Get: Obtiene el valor asociado a una clave
func (d *Dictionary) Get(key string) (interface{}, bool) {
	value, exists := d.items[key]
	return value, exists
}

// Remove: Elimina una clave-valor y actualiza la lista de claves
func (d *Dictionary) Remove(key string) {
	if _, exists := d.items[key]; exists {
		delete(d.items, key)

		// Elimina la clave de la lista de claves
		for i, k := range d.keys {
			if k == key {
				d.keys = append(d.keys[:i], d.keys[i+1:]...)
				break
			}
		}
	}
}

// IsEmpty: Verifica si el diccionario está vacío
func (d *Dictionary) IsEmpty() bool {
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
	for _, key := range d.keys {
		fmt.Printf("%s: %v\n", key, d.items[key])
	}
}
