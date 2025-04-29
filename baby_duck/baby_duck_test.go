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
	}, // Accept 1
	{
		`program testProgram;
         main {
            print("Hello, world!");
         }
         end`,
	}, // Accept 2
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
	}, // Accept 3
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
	}, // Accept 4
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
	}, // Accept 5
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
	}, // Accept 6
	{
		`program multiVars;
         var a, b: int;
         main {
            a = 1;
            b = 2.5;
            c = a + 2;
            print(a);
            print(b);
            print(c);
         }
         end`,
	}, // Accept 7
}

var testDataFail = []*TI{
	{
		`program missingMain;
         var x: int;
         {
            x = 5;
         }
         end`,
	}, // Fail 1
	{
		`program missingMain;
         var x: int;
         {
            x = 5;
         }
         end`,
	}, // Fail 2
	{
		`program noEnd;
         var x: int;
         main {
            x = 1;
         }`,
	}, // Fail 3
	{
		`program missingSemicolon;
         var a: int;
         main {
            a = 5
            print(a);
         }
         end`,
	}, // Fail 4
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
	}, // Fail 5
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
	}, // Fail 6
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
