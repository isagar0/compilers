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
		`program ejemplo1;
		 var x : int;
		 main {
			 x = 5;
		 }
		 end`,
	}, // Accept 1: Variable global declarada y usada correctamente
	/*
		{
			`program funcionesLocales;
			 void foo()[
				var a : int;
				{
					print(a);
				}
			 ];
			 void bar()[
				var a : float;
				{
					print(a);
				}
			 ];
			 main {
				 foo();
				 bar();
			 }
			 end`,
		}, // Accept 2: Misma variable local 'a' en distintas funciones
	*/
	{
		`program sinVars;
		 main {
			 print("hello");
		 }
		 end`,
	}, // Accept 3: Programa sin variables
	{
		`program prueba;
		 void test()[
			{}
		 ];
		 main {
			
		 }
		 end`,
	}, // Accept 4: Registro de función foo sin parámetros ni vars
	{
		`program conParams;
		 void paramsCheck(a: int, b: float)[
			{}
		 ];
		 main {
		 }
		 end`,
	}, // Accept 5: Función con parámetros `a:int` y `b:float`
}

var testDataFail = []*TI{
	{
		`program dupVar;
		 var x: int;
		 var x: float;
		 main {
			 x = 5;
		 }
		 end`,
	}, // Fail 1: Duplicación de variable global 'x'
	{
		`program dupVarLocal;
		 void algo()[ 
			var y: int;
			var y: float;
			{
				print(y);
			}
		 ];
		 main {
			 algo();
		 }
		 end`,
	}, // Fail 2: Duplicación de variable local 'y' dentro de una función
	{
		`program dupFunc;
		 void foo()[
		 	{}
		 ];
		 void foo()[
		 	{}
		 ];
		 main {
		 }
		 end`,
	}, // Fail 3: Duplicación de función 'foo'
	{
		`program dupParam;
		 void h(a: int, a: float)[
		 	{}
		 ];
		 main {
		 }
		 end`,
	}, // Fail 4: Duplicación de parámetro 'a' en la misma función
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
