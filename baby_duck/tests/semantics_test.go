package tests

import (
	"baby_duck/lexer"
	"baby_duck/parser"
	"baby_duck/semantics"
	"fmt"
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
	{
		`program funcionesLocales;
		 var x : int;
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
			print(x);
		 }
		 end`,
	}, // Accept 2: Misma variable local 'a' en distintas funciones
	{
		`program sinVars;
			 var a : int;
			 var b: float;
			 var c: int;
			 main {
				 print("hello");
			 }
			 end`,
	}, // Accept 3: Programa sin variables o parametros
	{
		`program prueba;
			 void test()[
				{}
			 ];
			 main {

			 }
			 end`,
	}, // Accept 4: Registro de función foo con parametros globales
	{
		`program conParams;
			 void paramsCheck(a: int, b: float)[
				{}
			 ];
			 main {
			 }
			 end`,
	}, // Accept 5: Función con parámetros `a:int` y `b:float`
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
	}, // Accept 6: Funcion con variables globales y locales
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
	{
		`program funcionesLocales;
		 var x : int;
		 void foo()[
		 var a : int;
		 var a: float;
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
			print(x);
		 }
		 end`,
	}, // Accept 5: Misma variable dentro de funcion
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

		// —— Al final, imprimimos las tablas ——
		fmt.Println("=== Tabla de variables globales ===")
		semantics.Current().PrintOrdered()

		fmt.Println("\n=== Directorio de funciones y sus tablas locales ===")
		for _, fname := range semantics.FunctionDirectory.Keys() {
			fmt.Printf("Función %s:\n", fname)
			entry, _ := semantics.FunctionDirectory.Get(fname)
			fs := entry.(semantics.FunctionStructure)
			fs.VarTable.PrintOrdered()
			fmt.Println()
		}
	}
}

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
