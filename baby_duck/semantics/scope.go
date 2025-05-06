package semantics

// scope.go

// ScopeManager maneja la tabla global y una pila de scopes locales
type ScopeManager struct {
	global *Dictionary
	stack  []*Dictionary
}

// instancia única
var scopes = &ScopeManager{
	global: NewDictionary(),
	stack:  []*Dictionary{},
}

// Current devuelve la tabla activa (local o global)
func Current() *Dictionary {
	if len(scopes.stack) == 0 {
		return scopes.global
	}
	return scopes.stack[len(scopes.stack)-1]
}

// EnterScope abre un nuevo scope local
func EnterScope() {
	child := NewDictionary()
	child.parent = scopes.global
	scopes.stack = append(scopes.stack, child)
}

// ExitScope cierra el scope local actual
func ExitScope() {
	if len(scopes.stack) > 0 {
		scopes.stack = scopes.stack[:len(scopes.stack)-1]
	}
}
