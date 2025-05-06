package semantics

// scope.go
import "fmt"

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
	fmt.Printf("[DEBUG] EnterScope → new local scope %p (parent %p), depth=%d\n",
		child, child.parent, len(scopes.stack))
}

// ExitScope cierra el scope local actual
func ExitScope() {
	if len(scopes.stack) == 0 {
		// nada que cerrar
		return
	}
	// 1) capturamos la tabla que vamos a sacar
	popped := scopes.stack[len(scopes.stack)-1]
	// 2) la quitamos de la pila
	scopes.stack = scopes.stack[:len(scopes.stack)-1]
	// 3) imprimimos el debug con el nuevo depth y el Current()
	fmt.Printf("[DEBUG] ExitScope  → popped scope %p, new depth=%d, current=%p\n",
		popped, len(scopes.stack), Current())
}
