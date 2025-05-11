package semantics

import "fmt"

var VarTable = NewDictionary()
var FunctionDirectory = NewDictionary()

// Pilas para cuadruplos
var (
	PilaO   = NewStack()
	PTypes  = NewStack()
	POper   = NewStack()
	Quads   []QuadStructure
	tempVar int
)

// ------------------------------------ Expreciones
// Función auxiliar para generar temporales
func newTemp() string {
	tempVar++
	return fmt.Sprintf("t%d", tempVar)
}

// PushQuad agrega un cuadruplo a la lista
func PushQuad(oper string, left, right, res interface{}) {
	Quads = append(Quads, QuadStructure{oper, left, right, res})
}

// ------------------------------------ Vars

// Reset reinicia el scope global y limpia la pila de scopes locales.
func Reset() {
	scopes = &ScopeManager{
		global: NewDictionary(),
		stack:  []*Dictionary{},
	}
}

// ResetSemanticState limpia el directorio de funciones y el scope global.
func ResetSemanticState() {
	// 1) reinicia el scope manager
	Reset() // de scope.go
	// 2) reinicia el FunctionDirectory
	FunctionDirectory = NewDictionary()
}

// RegisterMainProgram crea el scope global y registra el programa principal.
func RegisterMainProgram(programName string) error {
	// 1) volver a un estado limpio
	// ResetSemanticState()

	// 2) error si ya existe el programa
	if _, exists := FunctionDirectory.Get(programName); exists {
		return fmt.Errorf("error: el programa '%s' ya ha sido declarado", programName)
	}

	// 3) registra el programa como función void, usando la tabla global actual
	FunctionDirectory.Put(programName, FunctionStructure{
		Name:     programName,
		VarTable: Current(), // aquí guardamos el scope global
	})

	fmt.Printf("Programa principal '%s' registrado exitosamente.\n", programName)

	/*fmt.Printf("[DEBUG] RegisterMainProgram %s → global scope %p\n",
	programName, Current())*/
	return nil
}

// Función para procesar la declaración de variables en el scope actual
func VarDeclaration(ids []string, tipo string) error {
	tabla := Current() // usa la tabla activa de scope.go
	/*fmt.Printf("[DEBUG] VarDeclaration %v in scope %p, parent=%p, before=%v\n",
	ids, tabla, tabla.parent, tabla.Keys())*/

	for _, id := range ids {
		if _, exists := tabla.Get(id); exists {
			return fmt.Errorf("error: variable '%s' ya declarada en este ámbito", id)
		}
		tabla.Put(id, VariableStructure{Name: id, Type: tipo})
	}

	/*
		// ——— Imprimimos el contenido del scope actual ———
		fmt.Println(">>> Contenido del scope actual antes de declarar:")
		tabla.PrintOrdered()
		fmt.Println(">>> Fin del scope actual")
	*/

	/*fmt.Printf("[DEBUG] VarDeclaration %v done in scope %p, after=%v\n",
	ids, tabla, tabla.Keys())*/

	return nil
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

// FuncDeclaration actualiza la entrada creada por RegisterFunction
func FuncDeclaration(name string, params []VariableStructure) error {
	if err := ValidateParams(params); err != nil {
		return err
	}
	raw, exists := FunctionDirectory.Get(name)
	if !exists {
		return fmt.Errorf("error interno: función '%s' no registrada previamente", name)
	}
	fs := raw.(FunctionStructure)
	fs.Parameters = params
	fs.VarTable = Current() // la tabla local donde metimos los params
	FunctionDirectory.Put(name, fs)
	return nil
}

// RegisterFunction crea la entrada de la función con nombre, retorno void,
// sin parámetros todavía, y le asigna una tabla de variables local vacía.
func RegisterFunction(name string) error {
	if _, exists := FunctionDirectory.Get(name); exists {
		return fmt.Errorf("error: función '%s' ya declarada", name)
	}

	// 1) creamos la tabla local (de momento vacía)
	localTable := NewDictionary()

	/*fmt.Printf("[DEBUG] RegisterFunction %s → local scope %p\n",
	name, localTable)*/

	// 2) registramos la función en el directorio
	FunctionDirectory.Put(name, FunctionStructure{
		Name:       name,
		Parameters: []VariableStructure{}, // params vacía por ahora
		VarTable:   localTable,
	})

	return nil
}
