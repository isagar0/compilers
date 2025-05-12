package main

import (
	"baby_duck/lexer"
	"baby_duck/parser"
	"baby_duck/semantics"
	"testing"
)

type TI struct {
	src string
}

var testDataAccept = []*TI{
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
	}, // Accept 7: Uso de variable en expresion
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
	}, // Accept 7: Uso de variable en expresion
}

func TestSemanticAccept(t *testing.T) {
	p := parser.NewParser()
	for i, ts := range testDataAccept {
		// Reiniciamos semántica antes de empezar
		semantics.ResetSemanticState()

		s := lexer.NewLexer([]byte(ts.src))
		_, err := p.Parse(s)
		if err != nil {
			t.Errorf("Test %d (ACCEPT) failed: unexpected parse error.\nSource start: %.50s...\nError: %s",
				i+1, ts.src, err.Error())
			continue
		}

	}
}

/*
func TestSemanticFail(t *testing.T) {
	p := parser.NewParser()
	for i, ts := range testDataFail {
		// Reiniciamos semántica antes de empezar
		semantics.ResetSemanticState()

		s := lexer.NewLexer([]byte(ts.src))
		_, err := p.Parse(s)
		if err == nil {
			t.Errorf("Test %d (FAIL) did not produce expected error.\nSource start: %.50s...", i+1, ts.src)
		} else {
			t.Logf("Test %d (FAIL): Expected fail. Error: %s", i+1, err.Error())
		}
	}
}
*/
