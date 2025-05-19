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
		`program SimpleIf;
			var a, b, c, d : int;
            main {
                a = 1;
                b = 2;
                c = 3;
                d = 4;
                if (a + b > c * d) {
					a = b + d;
				};
				b = a * c;
            }
            end`,
	}, // Accept 1: If simple
	{
		`program IfWithElse;
			var a, b, c, d : int;
            main {
                a = 1;
                b = 2;
                c = 3;
                d = 4;
                if (a + b != c * d) {
					a = b + c;
				}
				else {
					a = d - c;
				};
				b = a * c + d;
            }
            end`,
	}, // Accept 2: If - else
	{
		`program IfWithElse;
			var a, b, c, d : int;
			main {
				a = 1;
				b = 2;
				c = 3;
				d = 4;
				if (a > b) {
					b = c * d;
					if (b < c + d){
						c = a + b;
						print(b);
					};
				}
				else {
					c = a + b;
					if (a > c){
						d = b + a;
						print(a + b);
					};
				};
				c = b - d * a;
			}
			end`,
	}, // Accept 3: Complicated If - else
	{
		`program withCycle;
         var a, b, c, d : int;
         main {
            while (a > b * c) do {
               a = a - d;
			   print(a);
            };
			b = c + a;
         }
         end`,
	}, // Accept 4: Uso de ciclo while-do
	{
		`program withBoth;
         var a, b, c, d, e, f, g, h, j, k : int;
         main {
		 	a = b + c * (d - e / f) * h;
			b = e - f;
            while (a * b - c > d * e / (g + h) ) do {
               h = j * k + b;
			   if (b < h){
					b = h + j;
					while (b > a + c) do {
						print(a + b * c, d - e);
						b = b - j;
					};
			   }
			   else {
			   		while (a - d < c + b) do {
						a = a + b;
						print(b - d);
					};
			   };
            };
			f = a + b;
         }
         end`,
	}, // Accept 5: Uso de if con ciclo while-do
}

func TestSemanticAccept(t *testing.T) {
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
