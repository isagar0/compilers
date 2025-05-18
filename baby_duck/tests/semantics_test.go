package main

import (
	"fmt"
	"testing"

	"baby_duck/lexer"
	"baby_duck/parser"
	"baby_duck/semantics"
)

type TI2 struct {
	src string
}

var testDataAccept2 = []*TI2{
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
	}, // Accept 3: Registro de función con parametros globales
	{
		`program prueba;
			 void test()[
				{}
			 ];
			 main {

			 }
			 end`,
	}, // Accept 4: Programa sin variables o parametros
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
			foo();
		 }
		 ];
		 main {
			bar();
			print(x);
		 }
		 end`,
	}, // Accept 8: Llamar funcion desde otra funcion
}

var testDataFail2 = []*TI2{
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
	}, // Fail 5: Misma variable dentro de funcion
	{
		`program useUndecl;
			main {
				x = 5;
			}
			end`,
	}, // Fail 6: Asignación a 'x' no declarada
	{
		`program useUndeclExpr;
			var x: int;
			main {
				x = y;
			}
			end`,
	}, // Fail 7: Uso de 'y' en expresión antes de declararla
	{
		`program funcUndecl;
			void foo()[{
				print(v);
			}];
			main {
			}
			end`,
	}, // Fail 8: Dentro de foo, 'v' no declarada
	{
		`program repetirFun;
		 void fun()[
		 var a : int;
		 {
			print(a);
		 }
		 ];
		 void fun()[
		 var a : float;
		 {
			print(a);
		 }
		 ];
		 main {
			fun();
		}
		 end`,
	}, // Fail 9: Nombre de funcion repetida
	{
		`program undeclaredCall;
			main {
				bar();
			}
			end`,
	}, // Fail 10: Llamada a función 'bar' no declarada
	{
		`program tooFewArgs;
			void sum(a: int, b: float)[{}];
			main {
				sum(1);
			}
			end`,
	}, // Fail 11: sum espera 2 args, recibe 1
	{
		`program diffType;
			void sum(a: int, b: float)[{}];
			main {
				sum(1, 2);
			}
			end`,
	}, // Fail 11: sum espera int float, recibe dos ints
}

func TestSemanticAccept2(t *testing.T) {
	p := parser.NewParser()
	for i, ts := range testDataAccept2 {
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

func TestSemanticFail2(t *testing.T) {
	p := parser.NewParser()
	for i, ts := range testDataFail2 {
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
