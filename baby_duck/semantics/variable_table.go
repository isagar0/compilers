package semantics

import (
	"fmt"
	"sort"
)

// VariableInfo representa la información de una variable
type VariableInfo struct {
	Name    string
	Type    string
	Scope   string // global, local o temporal
	Address int    // para futura generación de código
	// Puedes agregar más campos aquí según lo necesites
}

// VariableTable gestiona múltiples tablas de variables, una por función
type VariableTable struct {
	Tables map[string]map[string]VariableInfo // función -> (variable -> info)
}

// NewVariableTable crea una nueva instancia de VariableTable
func NewVariableTable() *VariableTable {
	return &VariableTable{
		Tables: make(map[string]map[string]VariableInfo),
	}
}

// Put agrega una nueva variable a una función específica
func (vt *VariableTable) Put(functionName string, varInfo VariableInfo) error {
	if vt.Tables[functionName] == nil {
		vt.Tables[functionName] = make(map[string]VariableInfo)
	}

	if _, exists := vt.Tables[functionName][varInfo.Name]; exists {
		return fmt.Errorf("la variable '%s' ya está declarada en la función '%s'", varInfo.Name, functionName)
	}

	vt.Tables[functionName][varInfo.Name] = varInfo
	return nil
}

// Get obtiene la información de una variable en una función específica
func (vt *VariableTable) Get(functionName string, varName string) (VariableInfo, error) {
	if functionVars, exists := vt.Tables[functionName]; exists {
		if info, found := functionVars[varName]; found {
			return info, nil
		}
		return VariableInfo{}, fmt.Errorf("la variable '%s' no está declarada en la función '%s'", varName, functionName)
	}
	return VariableInfo{}, fmt.Errorf("la función '%s' no tiene variables registradas", functionName)
}

// Remove elimina una variable de una función
func (vt *VariableTable) Remove(functionName string, varName string) error {
	if functionVars, exists := vt.Tables[functionName]; exists {
		if _, found := functionVars[varName]; found {
			delete(functionVars, varName)
			return nil
		}
		return fmt.Errorf("la variable '%s' no existe en la función '%s'", varName, functionName)
	}
	return fmt.Errorf("la función '%s' no tiene variables registradas", functionName)
}

// IsEmpty verifica si la tabla de variables de una función está vacía
func (vt *VariableTable) IsEmpty(functionName string) (bool, error) {
	if functionVars, exists := vt.Tables[functionName]; exists {
		return len(functionVars) == 0, nil
	}
	return false, fmt.Errorf("la función '%s' no existe en la tabla de variables", functionName)
}

// Size devuelve la cantidad de variables en una función
func (vt *VariableTable) Size(functionName string) (int, error) {
	if functionVars, exists := vt.Tables[functionName]; exists {
		return len(functionVars), nil
	}
	return 0, fmt.Errorf("la función '%s' no existe en la tabla de variables", functionName)
}

// Keys devuelve una lista ordenada alfabéticamente de los nombres de variables en una función
func (vt *VariableTable) Keys(functionName string) ([]string, error) {
	if functionVars, exists := vt.Tables[functionName]; exists {
		keys := make([]string, 0, len(functionVars))
		for k := range functionVars {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		return keys, nil
	}
	return nil, fmt.Errorf("la función '%s' no tiene variables registradas", functionName)
}

// PrintOrdered imprime las variables de una función ordenadas alfabéticamente
func (vt *VariableTable) PrintOrdered(functionName string) error {
	keys, err := vt.Keys(functionName)
	if err != nil {
		return err
	}
	for _, k := range keys {
		v := vt.Tables[functionName][k]
		fmt.Printf("Variable: %s, Tipo: %s, Scope: %s, Address: %d\n", v.Name, v.Type, v.Scope, v.Address)
	}
	return nil
}
