package main

import (
	"baby_duck/lexer"
	"baby_duck/parser"
	"baby_duck/semantics"
	"fmt"
	"testing"
)

type TI4 struct {
	src string
}

var testDataAccept4 = []*TI4{
	{
		`program Recursion;
		 void factorial(x: int, y: int) [
			{
				print(x);				
			}
		 ];

		 main {
			factorial(3, 7);
		 }

		 end`,
	},
}

func TestSemanticAccept(t *testing.T) {
	p := parser.NewParser()
	for i, ts := range testDataAccept4 {
		// Reiniciamos semántica antes de empezar
		semantics.ResetSemanticState()

		s := lexer.NewLexer([]byte(ts.src))
		_, err := p.Parse(s)
		if err != nil {
			t.Errorf("Test %d (ACCEPT) failed: unexpected parse error.\nSource start: %.50s...\nError: %s",
				i+1, ts.src, err.Error())
			continue
		}
		semantics.PrintQuads()
		// semantics.PrintAddressTable()
		fmt.Println("\n===========================================================")

		// Ejecutar cuádruplos
		// print(semantics.FunctionDirectory)
		// fmt.Println("\n!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		vm := semantics.NewVirtualMachine(semantics.Quads, semantics.FunctionDirectory)
		vm.Run()

		fmt.Println("\n===========================================================")

		/*
			// Después de parsear:
			fmt.Println("\n=== Funciones registradas ===")
			for name := range semantics.FunctionDirectory.Items {
				fmt.Println("Función:", name)
			}

			// En TestSemanticAccept:
			mainEntry, exists := semantics.FunctionDirectory.Get("main")
			if !exists {
				t.Errorf("Test %d: 'main' no está en el directorio", i+1)
				return
			}

			// Verificar que 'main' tiene 0 parámetros y variables locales correctas
			fs := mainEntry.(semantics.FunctionStructure)
			if fs.ParamCount != 0 || fs.LocalVarCount != 0 { // Ejemplo: 1 variable local
				t.Errorf("Test %d: Params=%d (esperaba 0), Locales=%d (esperaba 1)", i+1, fs.ParamCount, fs.LocalVarCount)
			}
		*/
	}
}
