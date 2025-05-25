package semantics

import (
	"fmt"
)

// --------------------------------------- Inicialización ---------------------------------------
var VarTable = NewDictionary()          // Tabla global de variables
var FunctionDirectory = NewDictionary() // Directorio de funciones
var memory = NewMemoryManager()         // Memoria de direcciones
var AddressToName = map[int]string{}    // Traducir direcciones a nombre

// ------------------------------------------ Limpiar ------------------------------------------
// ResetSemanticState: Limpia todo para un programa nuevo
func ResetSemanticState() {
	// Limpia diccionarios
	Reset()
	FunctionDirectory = NewDictionary()

	// Limpia pilas y cuadruplos
	CleanStacks()

	// Limpia direcciones
	// Reset all memory segments
	memory.Global.Ints.Reset()
	memory.Global.Floats.Reset()
	memory.Local.Ints.Reset()
	memory.Local.Floats.Reset()
	memory.Temp.Ints.Reset()
	memory.Temp.Floats.Reset()
	memory.Constant.Ints.Reset()
	memory.Constant.Floats.Reset()
	memory.Constant.Strings.Reset()
	AddressToName = make(map[int]string)
}

// -------------------------------------------- Vars --------------------------------------------
// Reset: Crea una nuvea tabla global y sus scopes vacios
func Reset() {
	scopes = &ScopeManager{
		global: NewDictionary(), // Scope global
		stack:  []*Dictionary{}, // Pila soces locales vacio
	}
}

// RegisterMainProgram: Crea el scope global y registra el programa principal
func RegisterMainProgram(programName string) error {
	// Verifica si ya existe una entrada con mismo nombre
	if _, exists := FunctionDirectory.Get(programName); exists {
		return fmt.Errorf("error: el programa '%s' ya ha sido declarado", programName)
	}

	// Registra el programa principal como una función sin parámetros
	FunctionDirectory.Put(programName, FunctionStructure{
		Name:     programName,
		VarTable: Current(),
	})

	//fmt.Printf("Programa principal '%s' registrado exitosamente.\n", programName)
	return nil
}

// VarDeclaration: Procesa la declaración de variables en el scope actual
func VarDeclaration(ids []string, tipo string) error {
	// Usa tabla activa del scope
	tabla := Current()

	// Recorre cada identificador en la lista de variables a declarar
	for _, id := range ids {
		var dir int
		var err error

		// Check current scope
		if _, exists := tabla.Get(id); exists {
			return fmt.Errorf("error: variable '%s' ya declarada en este scope", id)
		}

		// Check parent scopes if in local scope
		if tabla != scopes.global {
			if _, exists := scopes.global.Get(id); exists {
				return fmt.Errorf("error: variable '%s' ya existe en scope global", id)
			}
		}

		// Determina si es global o local
		var segmento *SegmentGroup
		if Current() == scopes.global {
			segmento = &memory.Global
			// fmt.Printf("Global var %s at %d\n", id, dir)
		} else {
			segmento = &memory.Local
			//fmt.Printf("Local var %s at %d\n", id, dir)
		}

		switch tipo {
		case "int":
			dir, err = segmento.Ints.GetNext()
		case "float":
			dir, err = segmento.Floats.GetNext()
		default:
			return fmt.Errorf("tipo no soportado: %s", tipo)
		}
		if err != nil {
			return err
		}

		// Agrega la variable a la tabla con dirección virtual
		tabla.Put(id, VariableStructure{
			Name:    id,
			Type:    tipo,
			Address: dir,
		})
		AddressToName[dir] = id
		//fmt.Printf("Declared %s at address %d (type %s)\n", id, dir, tipo)
	}

	/*
		// Imprimir el contenido del scope actual
		fmt.Println(">>> Contenido del scope actual antes de declarar:")
		tabla.PrintOrdered()
		fmt.Println(">>> Fin del scope actual")
	*/

	return nil
}

// RegisterFunction: Crea la entrada de la función con nombre, retorno void
func RegisterFunction(name string) error {
	// Verifica si ya existe una función con el mismo nombre, marca error
	if _, exists := FunctionDirectory.Get(name); exists {
		return fmt.Errorf("error: función '%s' ya declarada", name)
	}

	// Crea una nueva tabla de variables locales para esta función
	localTable := NewDictionary()

	/*fmt.Printf("[DEBUG] RegisterFunction %s → local scope %p\n",
	name, localTable)*/

	// Registra la función en el directorio
	FunctionDirectory.Put(name, FunctionStructure{
		Name:       name,                  // Nombre
		Parameters: []VariableStructure{}, // Parametros (vacios)
		VarTable:   localTable,            // Tabla local de variables
	})

	return nil
}

// ValidateParams: Verifica que los parámetros de una función no estén duplicados
func ValidateParams(params []VariableStructure) error {
	// Diccionario temporal para llevar el control de nombres ya vistos
	paramSet := NewDictionary()

	// Recorre cada parametro recibido
	for _, param := range params {
		// Verifica si el parametro ya fue declarada, marca error
		if _, exists := paramSet.Get(param.Name); exists {
			return fmt.Errorf("error: parámetro '%s' duplicado en la función", param.Name)
		}

		// Si no existe, se agrega para futuras comparaciones
		paramSet.Put(param.Name, param)
	}
	return nil
}

// FuncDeclaration: Actualiza la entrada creada por RegisterFunction
func FuncDeclaration(name string, params []VariableStructure) error {
	// Verifica que no haya parámetros duplicados
	if err := ValidateParams(params); err != nil {
		return err
	}

	// Busca la función en el directorio, marca error
	raw, exists := FunctionDirectory.Get(name)
	if !exists {
		return fmt.Errorf("error interno: función '%s' no registrada previamente", name)
	}

	// Convierte la entrada a una estructura de función
	fs := raw.(FunctionStructure)

	// Asigna los parámetros recibidos a la función
	fs.Parameters = params

	for _, param := range params {
		// Get the actual address from the local scope
		if raw, exists := Current().Get(param.Name); exists {
			vs := raw.(VariableStructure)
			AddressToName[vs.Address] = param.Name
			// return fmt.Errorf("parameter %s already exists", param.Name)
			//fmt.Printf("Registered param %s → %d (actual address)\n", param.Name, vs.Address)
		}
	}

	// Asocia la tabla local de variables (scope actual donde se declararon los params)
	fs.VarTable = Current()

	// Actualiza el directorio
	FunctionDirectory.Put(name, fs)

	return nil
}
