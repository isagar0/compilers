package semantics

// scopes: Instancia global
var scopes = &ScopeManager{
	global: NewDictionary(),
	stack:  []*Dictionary{},
}

// Current: Devuelve la tabla activa (local o global)
func Current() *Dictionary {
	// Si no hay scopes locales activos, retorna el scope global.
	if len(scopes.stack) == 0 {
		return scopes.global
	}
	return scopes.stack[len(scopes.stack)-1] // último scope en la pila
}

// EnterScope: Crea y abre un nuevo scope local
func EnterScope() {
	child := NewDictionary()                   // Crea tabla de simbolos
	child.parent = scopes.global               //Asigna el scope global como el padre
	scopes.stack = append(scopes.stack, child) // Lo agrega a la pila de scopes

	/*fmt.Printf("[DEBUG] EnterScope → new local scope %p (parent %p), depth=%d\n",
	child, child.parent, len(scopes.stack))*/
}

// ExitScope: Cierra el scope local actual
func ExitScope() {
	// No hay scopes locales que cerrar
	if len(scopes.stack) == 0 {
		// nada que cerrar
		return
	}

	scopes.stack = scopes.stack[:len(scopes.stack)-1]
}
