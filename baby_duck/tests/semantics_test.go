package tests

import (
	"baby_duck/lexer"
	"baby_duck/parser"
	"testing"
)

type TI struct {
	src string
}

var testDataAccept = []*TI{
	{
		`program withFuncs;
         void one()[
            var a: int;
            {
                print(a);
            }
         ];
		 void two()[
            var a: int;
            {
                print(a);
            }
         ];
         main {
            one();
			two();
         }
         end`,
	}, // Accept 1: Misma variable en diferentes funciones
	{
		`void sumar();
		var 
			x : int;
			y : float;
		{
			x = 3;
			y = 4.5;
		}
		`,
	}, // Accept 1: Misma variable en diferentes funciones
}

var testDataFail = []*TI{
	{
		`program dupVar;
         var x: int;
         main {
            x = 5;
         }
         end`,
	}, // Fail 1: Duplicación de variable global 'x'
	{
		`program withTwoLocs;
         void one()[
            var a: int;
			var a: float;
            {
                print(a);
            }
         ];
         main {
            one();
         }
         end`,
	}, // Fail 2: Duplicación de variable local 'x' dentro de la misma función
}

func TestSemanticAccept(t *testing.T) {
	p := parser.NewParser()
	for i, ts := range testDataAccept {
		s := lexer.NewLexer([]byte(ts.src))
		_, err := p.Parse(s)
		if err != nil {
			t.Errorf("Test %d (ACCEPT) failed: unexpected parse error.\nSource start: %.50s...\nError: %s", i+1, ts.src, err.Error())
		}
	}
}

func TestSemanticFail(t *testing.T) {
	p := parser.NewParser()
	for i, ts := range testDataFail {
		s := lexer.NewLexer([]byte(ts.src))
		_, err := p.Parse(s)
		if err == nil {
			t.Errorf("Test %d (FAIL) did not produce expected error.\nSource start: %.50s...", i+1, ts.src)
		} else {
			t.Logf("Test %d (FAIL): Expected fail. Error: %s", i+1, err.Error())
		}
	}
}
