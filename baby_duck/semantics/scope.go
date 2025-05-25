package semantics

// scopes: Instancia global
var Scopes = &ScopeManager{
	global:  NewDictionary(),
	current: nil,
}

// NewScopeManager: Inicializa con el scope global
func NewScopeManager() *ScopeManager {
	return &ScopeManager{
		global:  NewDictionary(),
		current: nil,
	}
}

// EnterScope: Crea un nuevo scope y lo establece como actual
func (s *ScopeManager) EnterScope() {
	newScope := NewDictionary()
	newScope.parent = s.current // Enlaza al scope anterior
	s.current = newScope
}

// ExitScope: Regresa al scope padre
func (s *ScopeManager) ExitScope() {
	if s.current != nil {
		s.current = s.current.parent
	}
}

// Current: Devuelve el scope actual (o global si no hay locales)
func (s *ScopeManager) Current() *Dictionary {
	if s.current == nil {
		return s.global
	}
	return s.current
}
