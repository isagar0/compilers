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
	}, // 1. Aceptado
	{
		`program testProgram;
         main {
            print(\"Hello, world!\");
         }
         end`,
	}, // 2. Aceptado
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
            }
         }
         end`,
	}, // 3. Aceptado
	{
		`program withFunc;
         void sum(int a, int b)[
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
	}, // 4. Aceptado
	{
		`program funcOnly;
         void greet()[
            {
               print("Hi!");
            }
         ];
         main {
            greet();
         }
         end`,
	}, // 5. Aceptado
	{
		`program withCycle;
         var i: int;
         main {
            i = 0;
            while (i < 10) do {
               print(i);
               i = i + 1;
            }
         }
         end`,
	}, // 6. Aceptado
	{
		`program multiVars;
         var a: int, b: float, c: int;
         main {
            a = 1;
            b = 2.5;
            c = a + 2;
            print(a);
            print(b);
            print(c);
         }
         end`,
	}, // 7. Aceptado
	{
		`program missingMain;
         var x: int;
         {
            x = 5;
         }
         end`,
	},
}

var testDataFail = []*TI{
	{
		`program missingMain;
         var x: int;
         {
            x = 5;
         }
         end`,
	}, // 8. Falla
	{
		`program noEnd;
         var x: int;
         main {
            x = 1;
         }`,
	}, // 9. Falla
	{
		`program missingSemicolon;
         var a: int;
         main {
            a = 5
            print(a);
         }
         end`,
	}, // 10. Falla
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
	}, // 11. Falla
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
	}, // 12. Falla
}

func TestParserAccept(t *testing.T) {
	p := parser.NewParser()
	pass := true
	for _, ts := range testDataAccept {
		s := lexer.NewLexer([]byte(ts.src))
		_, err := p.Parse(s)
		if err != nil {
			pass = false
			t.Logf("Error parsing source (start: %.50s...):\nError: %s", ts.src, err.Error())
		}
	}
	if !pass {
		t.Fail()
	}
}

func TestParserFail(t *testing.T) {
	p := parser.NewParser()
	pass := true
	for _, ts := range testDataFail {
		s := lexer.NewLexer([]byte(ts.src))
		_, err := p.Parse(s)
		if err != nil {
			pass = false
			t.Logf("Error parsing source (start: %.50s...):\nError: %s", ts.src, err.Error())
		}
	}
	if !pass {
		t.Fail()
	}
}
