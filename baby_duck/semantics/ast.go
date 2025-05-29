package semantics

import "fmt"

// ------------------------------------------ LIMPIAR ------------------------------------------
// ResetSemanticState: Limpia todo para un programa nuevo
func ResetSemanticState() {
	// Limpia diccionarios
	ResetVars()
	FunctionDirectory = NewDictionary()

	// Limpia pilas y cuadruplos
	ResetStacks()

	// Limpia direcciones
	ResetMemory()
}

// -------------------------------------------- VARS --------------------------------------------
// Reset: Crea una nuvea tabla global y sus scopes vacios
func ResetVars() {
	global := NewDictionary()
	Scopes = &ScopeManager{
		global:  global,
		current: global, // Ahora current apunta a global inicialmente
	}
}

// RegisterMainProgram: Crea el scope global y registra el programa principal
func RegisterMainProgram(programName string) error {
	// Verifica si ya existe una entrada con mismo nombre
	if _, exists := FunctionDirectory.Get(programName); exists {
		return fmt.Errorf("error: el programa '%s' ya ha sido declarado", programName)
	}

	/*
		// Registra el programa principal como una función sin parámetros
		FunctionDirectory.Put(programName, FunctionStructure{
			Name:       programName,
			Parameters: []VariableStructure{},
			VarTable:   Scopes.global,
			ParamCount: 0,
			StartQuad:  0,
		})
	*/

	//fmt.Printf("Programa principal '%s' registrado exitosamente.\n", programName)
	return nil
}

// VarDeclaration: Procesa la declaración de variables en el scope actual
func VarDeclaration(ids []string, tipo string) error {
	// Usa tabla activa del scope
	tabla := Scopes.Current()

	// Recorre cada identificador en la lista de variables a declarar
	for _, id := range ids {
		var dir int
		var err error

		// Check parent scopes if in local scope
		/*
			if tabla != Scopes.global {
				if _, exists := Scopes.global.Get(id); exists {
					return fmt.Errorf("error: variable '%s' ya existe en scope global", id)
				}
			}
		*/

		// Determina si es global o local
		var segmento *SegmentGroup
		if Scopes.Current() == Scopes.global {
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
		/*scope := "global"
		if Scopes.Current() != Scopes.global {
			scope = "local"
		}
		fmt.Printf("[DEBUG] Declaradas variables %v en scope %s (tipo %s)\n", ids, scope, tipo)*/
	}

	/*
		// Imprimir el contenido del scope actual
		fmt.Println(">>> Contenido del scope actual antes de declarar:")
		tabla.PrintOrdered()
		fmt.Println(">>> Fin del scope actual")
	*/

	return nil
}

// Registra el nombre del programa (no es una función)
func RegisterProgramName(name string) error {
	// Opcional: Almacenar en una estructura separada si es necesario
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
		ParamCount: 0,
		TempCount:  0,
		StartQuad:  len(Quads),
	})

	//fmt.Printf("[DEBUG] RegisterFunction: %d\n", len(Quads))

	return nil
}

// ValidateParams: Verifica que los parámetros de una función no estén duplicados
func ValidateParams(params []VariableStructure) error {
	// Diccionario temporal para llevar el control de nombres ya vistos
	paramSet := NewDictionary()

	// Recorre cada parametro recibido
	for _, param := range params {
		/*
			// Verifica si el parametro ya fue declarada, marca error
			if _, exists := paramSet.Get(param.Name); exists {
				return fmt.Errorf("error: parámetro '%s' duplicado en la función", param.Name)
			}
		*/

		// Si no existe, se agrega para futuras comparaciones
		paramSet.Put(param.Name, param)
	}
	return nil
}

// FuncDeclaration: Actualiza la entrada creada por RegisterFunction
func FuncDeclaration(name string, params []VariableStructure, localVarCount, startQuad, tempCount int) error {
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
	for i, param := range params {
		// paramName := fmt.Sprintf("param_%d", i)
		// Get the actual address from the local scope
		if raw, exists := Scopes.Current().Get(param.Name); exists {
			vs := raw.(VariableStructure)
			// Registrar con nombre especial para depuración			/
			// return fmt.Errorf("parameter %s already exists", param.Name)
			AddressToName[vs.Address] = fmt.Sprintf("%s_param_%d", name, i+1)
			// fmt.Printf("Registered param %s → %d (actual address)\n", param.Name, vs.Address)
		}
	}

	// Asigna los nuevos campos
	fs.ParamCount = len(params)      // Número de parámetros
	fs.LocalVarCount = localVarCount // Variables locales
	fs.StartQuad = startQuad         // Cuadruplo inicial
	fs.TempCount = tempCount         // Numero temporales

	// Asocia la tabla local de variables (scope actual donde se declararon los params)
	fs.VarTable = Scopes.Current()

	// Actualiza el directorio
	FunctionDirectory.Put(name, fs)

	return nil
}

// CountVars: Cuenta variables en el scope actual
func (d *Dictionary) CountVars() int {
	count := 0
	for _, v := range d.Items {
		if _, ok := v.(VariableStructure); ok {
			count++
		}
	}
	return count
}

// GetCurrentQuad: Obtiene el índice del último cuadruplo generado
func GetCurrentQuad() int {
	return len(Quads) - 1
}

func AssignAddressToParam(tipo string) (int, error) {
	switch tipo {
	case "int":
		return memory.Local.Ints.GetNext()
	case "float":
		return memory.Local.Floats.GetNext()
	default:
		return 0, fmt.Errorf("tipo no soportado: %s", tipo)
	}
}

func DeclareInCurrentScope(name, tipo string, address int) error {
	scope := Scopes.Current()

	scope.Put(name, VariableStructure{
		Name:    name,
		Type:    tipo,
		Address: address,
	})
	AddressToName[address] = name
	return nil
}
