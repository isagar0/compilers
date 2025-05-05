package semantics

import "fmt"

// Estructura de una variable
type VariableStructure struct {
	Name string
	Type string
}

// Estructura de una función
type FunctionStructure struct {
	Name       string
	Parameters []VariableStructure
	VarTable   *Dictionary
}

var VarTable = NewDictionary()
var FunctionDirectory = NewDictionary()

// Función para registrar el programa principal
func RegisterMainProgram(programName string) error {
	// Verificamos si ya existe el nombre del programa en el diccionario de funciones
	if _, exists := FunctionDirectory.Get(programName); exists {
		return fmt.Errorf("error: el programa '%s' ya ha sido declarado", programName)
	}

	// Inicializamos la tabla de variables global para el programa
	VarTable = NewDictionary()

	// Registramos el programa en el FunctionDirectory
	FunctionDirectory.Put(programName, FunctionStructure{
		Name:     programName,
		VarTable: VarTable, // Asignamos la tabla de variables globales
	})

	fmt.Printf("Programa principal '%s' registrado exitosamente.\n", programName)
	return nil
}

// Función para procesar la declaración de variables
func VarDeclaration(ids []string, tipo string, tabla *Dictionary) (*Dictionary, error) {
	if tabla == nil {
		tabla = NewDictionary()
	}

	// Iteramos sobre los identificadores (nombres de variables)
	for _, id := range ids {
		// Verificamos si la variable ya está declarada en la tabla
		if _, exists := tabla.Get(id); exists {
			// Si ya existe, lanzamos un error
			return nil, fmt.Errorf("error: variable '%s' ya declarada", id)
		}
		// Si no existe, la agregamos a la tabla
		tabla.Put(id, VariableStructure{Name: id, Type: tipo})
	}

	return tabla, nil
}

// Verifica que los parámetros de una función no estén duplicados
func ValidateParams(params []VariableStructure) error {
	paramSet := NewDictionary()
	for _, param := range params {
		if _, exists := paramSet.Get(param.Name); exists {
			return fmt.Errorf("error: parámetro '%s' duplicado en la función", param.Name)
		}
		paramSet.Put(param.Name, param)
	}
	return nil
}

// Procesa la declaración de una función
func FuncDeclaration(name string, params []VariableStructure, localVars *Dictionary) error {
	if _, exists := FunctionDirectory.Get(name); exists {
		return fmt.Errorf("error: funcion '%s' ya declarada", name)
	}

	// Valida que no haya parámetros duplicados
	if err := ValidateParams(params); err != nil {
		return err
	}

	// Se agrega la función a la tabla de funciones
	FunctionDirectory.Put(name, FunctionStructure{
		Name:       name,
		Parameters: params,
		VarTable:   localVars,
	})

	return nil
}

// Reinicia el estado semántico (limpia las tablas globales)
func ResetSemanticState() {
	VarTable = NewDictionary()
	FunctionDirectory = NewDictionary()
}
