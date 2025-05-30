package semantics

import "fmt"

// Define las reglas de tipos del lenguaje
var cube = SemanticCube{
	"int": {
		"int": {
			"+":  "int",
			"-":  "int",
			"*":  "int",
			"/":  "float",
			"=":  "int",
			"!=": "bool",
			"<":  "bool",
			">":  "bool",
		},
		"float": {
			"+":  "float",
			"-":  "float",
			"*":  "float",
			"/":  "float",
			"=":  "int",
			"!=": "bool",
			"<":  "bool",
			">":  "bool",
		},
	},
	"float": {
		"int": {
			"+":  "float",
			"-":  "float",
			"*":  "float",
			"/":  "float",
			"=":  "int",
			"!=": "bool",
			"<":  "bool",
			">":  "bool",
		},
		"float": {
			"+":  "float",
			"-":  "float",
			"*":  "float",
			"/":  "float",
			"=":  "int",
			"!=": "bool",
			"<":  "bool",
			">":  "bool",
		},
	},
	"bool": {
		"bool": {
			"!=": "bool",
		},
	},
}

// GetResultType: Devuelve el tipo del resultado según el cubo
func GetResultType(leftType, rightType, operator string) (string, error) {
	// Busca directamente en el cubo semántico
	if res, ok := cube[leftType][rightType][operator]; ok {

		// Devuelve el tipo
		return res, nil
	}

	// Si no existe esa combinación en el cubo, marca error
	return "", fmt.Errorf("semántico: combinación no definida %s %s %s", leftType, operator, rightType)
}
