package semantics

import "fmt"

// NewDictionary: Constructor del diccionario
func NewDictionary() *Dictionary {
	return &Dictionary{
		items:  make(map[string]interface{}),
		keys:   []string{},
		parent: nil,
	}
}

// NewChildDictionary: Constructor anidado
func NewChildDictionary(parent *Dictionary) *Dictionary {
	return &Dictionary{
		items:  make(map[string]interface{}),
		keys:   []string{},
		parent: parent,
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
	if val, ok := d.items[key]; ok {
		return val, true
	}
	if d.parent != nil {
		return d.parent.Get(key)
	}
	return nil, false
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
