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
				print(sum);
            }
            end`,
	}, // Accept 2: Ver que el orden de jerarquía funciona correctamente
	{
		`program SimpleIf;
			var a, b : int;
            main {
                a = 3;
                b = 2;
                print(a != b);
            }
            end`,
	}, // Accept 1: If simple
	{
		`program SimpleIf;
			var a : int;
            main {
                a = 11;
				if (a < 10) {
					print(a);
				}
				else {
					print("ELSE", 5>2, 8);
				};
            }
            end`,
	}, // Accept 1: If simple
	{
		`program withCycle;
         var a, b : int;
         main {
			a = 1;
			b = 4;
            while (a < b) do {
               a = a + 1;
			   print(a);
            };
			print("Out of While");
         }
         end`,
	}, // Accept 4: Uso de ciclo while-do
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
		vm := semantics.NewVirtualMachine(semantics.Quads)
		vm.Run() // ¡Esto ejecutará todos los cuádruplos!

		fmt.Println("\n===========================================================")
	}
}
