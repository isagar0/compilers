package main

import (
	"baby_duck/lexer"
	"baby_duck/parser"
	"baby_duck/semantics"
	"fmt"
	"testing"
)

type TI3 struct {
	src string
}

var testDataAccept3 = []*TI3{
	{
		`program withFunc;
			 var x : int;
	         void sum(a: int, b: int)[
	            var result: int;
	            {
	                result = a + b;
	                print(result);
	            }
	         ];
	         main {
			 	x = 2;
	            sum(x,3);
	         }
	         end`,
	}, // Accept 1: Uso de variables con una función
	{
		`program sumTest;
			var a, b, c, d, e, f, g, h, j, k, l : int;
            var sum : int;
            main {
                a = 1;
                b = 2;
                c = 3;
                d = 4;
                e = 5;
                f = 6;
                g = 7;
                h = 8;
                j = 10;
                k = 11;
                l = 12;
                sum = ( ( a + b ) * c + d * e * f + k / h * j ) + g * l + h + j + ( a - c * d ) / f;
            }
            end`,
	}, // Accept 2: Ver que el orden de jerarquía funciona correctamente
	{
		`program simple;
         main {
            if (3 + 5 > 2 * 8) {
                print(1);
            };
         }
         end`,
	}, // Accept 3: Prueba condicional con simbolos relaciones
	{
		`program withCycle;
         var i: int;
         main {
            i = 0;
            while (i < 10/5) do {
               print(i);
               i = i + 1;
            };
         }
         end`,
	}, // Accept 4: Prueba ciclo con simbolos relaciones
}

func TestSemanticAccept3(t *testing.T) {
	p := parser.NewParser()
	for i, ts := range testDataAccept3 {
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
		semantics.PrintAddressTable()
		fmt.Println("\n===========================================================")
	}
}
