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
		`program myProgram;
         var x, y: int;
         main {
            x = 1.2 + 2;
            print(x);
         }
         end`,
	}, // Accept 1: If simple
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
		semantics.ExecuteQuads(semantics.Quads)

		fmt.Println("\n===========================================================")
	}
}
