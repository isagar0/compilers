package test

import (
	"testing"

	"baby_duck/lexer"
	"baby_duck/parser"
)

type TI struct {
	src string
}

var testDataAccept = []*TI{
	{
		`program myProgram;
         var x, y: int;
         main {
            x = 1 + 2 * 3;
            print(x);
         }
         end`,
	}, // Accept 1: Declaración de variables, uso de operadores, función print
	{
		`program testProgram;
         main {
            print("Hello, world!");
         }
         end`,
	}, // Accept 2: Función main simplificada, variables y funciones opcionales, print de string
	{
		`program simple;
         var a: float;
         main {
            a = (1 + 2) * 3.5;
            if (a > 5) {
                print(a);
            }
            else {
                print("Small");
            };
         }
         end`,
	}, // Accept 3: Asignación de variables, condicionales if-else, comparadores
	{
		`program withFunc;
         void sum(a: int, b: int)[
            var result: int;
            {
                result = a + b;
                print(result);
            }
         ];
         main {
            sum(2,3);
         }
         end`,
	}, // Accept 4: Funciones que reciben parámetros, declaración de variables locales, llamada a funciones con parámetros
	{
		`program funcOnly;
         void greetAndPrint()[
            {
               print("Hi!", 1+2);
            }
         ];
         main {
            greetAndPrint();
         }
         end`,
	}, // Accept 5: Función print con strings y expresiones
	{
		`program withCycle;
         var i: int;
         main {
            i = 0;
            while (i < 10) do {
               print(i);
               i = i + 1;
            };
         }
         end`,
	}, // Accept 6: Uso de ciclo while-do
	{
		`program multiVars;
         var a, b, c: int;
         main {
            a = 1;
            b = 2;
            c = a + 2;
            print(a);
            print(b);
            print(c);
         }
         end`,
	}, // Accept 7: Asignación de variables, suma con constantes y variables
	{
		`program funcWithVars;
        void calc()[ var a, b: int; {
            a = 2;
            b = a * 3;
            print(b);
        }];
        main {
            calc();
        }
        end`,
	}, // Accept 8: Declaración local de variables, asignación de variables, llamada a una función
	{
		`program parentheses;
        var result: float;
        main {
            result = (2 + 3) * (4 - 1);
            print(result);
        }
        end`,
	}, // Accept 9: Combinación de expresiones con paréntesis
	{
		`program complexExpr;
        var a: int;
        main {
            a = 2 + 3 * -(-4 + 5);
            print(a);
        }
        end`,
	}, // Accept 10: Combinación de expresiones con operadores aritméticos
	{
		`program nestedIfs;
        var x: int;
        main {
            x = 5;
            if (x > 0) {
                print("Positive");
                if (x != 10) {
                    print("X is not equal to 10");
                }
                else {
                    print("10 or more");
                };
            }
            else {
                print("Non-positive");
            };
        }
        end`,
	}, // Accept 11: Uso de if-else anidados, condiciones con comparadores, print de string
	{
		`program wordTest;
         void mainA()[
            {
               print("Magic");
            }
         ];
         main {
            mainA();
         }
         end`,
	}, // Accept 12: Funciones que tienen palabras reservadas
}

var testDataFail = []*TI{
	{
		`program missingMain;
         var x: int;
         {
            x = 5;
         }
         end`,
	}, // Fail 1: Falta main, body inválido
	{
		`program noEnd;
         var x: int;
         main {
            x = 1;
         }`,
	}, // Fail 2: Falta end
	{
		`program missingSemicolon;
         var a: int;
         main {
            a = 5
            print(a);
         }
         end`,
	}, // Fail 3: Falta punto y coma
	{
		`program badFunc;
         void add(int a, int b){
            var result: int;
            {
               result = a + b;
               print(result);
            }
         ]
         main {
            add(1,2);
         }
         end`,
	}, // Fail 4: FUNCS mal declarada (paréntesis y llaves)
	{
		`program badWhile;
         var i: int;
         main {
            i = 0;
            while i < 10 do {
               print(i);
            }
         }
         end`,
	}, // Fail 5: while sin paréntesis
	{
		`program badIf;
        var a: int;
        main {
            if (a == 2);
        }
        end`,
	}, // Fail 6: if con operador inválido (==)
	{
		`program badWhile2;
        var i: int;
        main {
            while (i < 10) {
                print(i);
            }
        }
        end`,
	}, // Fail 7: while sin 'do'
}

func TestParserAccept(t *testing.T) {
	p := parser.NewParser()
	for i, ts := range testDataAccept {
		s := lexer.NewLexer([]byte(ts.src))
		_, err := p.Parse(s)
		if err != nil {
			t.Errorf("Test %d (ACCEPT) failed: unexpected parse error.\nSource start: %.50s...\nError: %s", i+1, ts.src, err.Error())
		}
	}
}

func TestParserFail(t *testing.T) {
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
