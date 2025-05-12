package semantics

import "fmt"

// SemanticCube: Estructura de 3 niveles que define el tipo de resultado
// Estructura: [tipoIzq][tipoDer][operador] → tipoResultado
type SemanticCube map[string]map[string]map[string]string

// Define las reglas de tipos del lenguaje
// Si una operación no está permitida, marca error
var cube = SemanticCube{
	"int": {
		"int": {
			"+":  "int",
			"-":  "int",
			"*":  "int",
			"/":  "float",
			"=":  "int",
			"!=": "int",
			"<":  "int",
			">":  "int",
		},
		"float": {
			"+":  "float",
			"-":  "float",
			"*":  "float",
			"/":  "float",
			"=":  "int",
			"!=": "int",
			"<":  "int",
			">":  "int",
		},
	},
	"float": {
		"int": {
			"+":  "float",
			"-":  "float",
			"*":  "float",
			"/":  "float",
			"=":  "int",
			"!=": "int",
			"<":  "int",
			">":  "int",
		},
		"float": {
			"+":  "float",
			"-":  "float",
			"*":  "float",
			"/":  "float",
			"=":  "int",
			"!=": "int",
			"<":  "int",
			">":  "int",
		},
	},
	"string": {
		"string": {
			"+":  "error",
			"-":  "error",
			"*":  "error",
			"/":  "error",
			"=":  "error",
			"!=": "error",
			"<":  "error",
			">":  "error",
		},
		"int": {
			"+":  "error",
			"-":  "error",
			"*":  "error",
			"/":  "error",
			"=":  "error",
			"!=": "error",
			"<":  "error",
			">":  "error",
		},
		"float": {
			"+":  "error",
			"-":  "error",
			"*":  "error",
			"/":  "error",
			"=":  "error",
			"!=": "error",
			"<":  "error",
			">":  "error",
		},
	},
}

// GetResultType: Devuelve el tipo del resultado según el cubo
func GetResultType(leftType, rightType, operator string) (string, error) {
	// Busca directamente en el cubo semántico
	if res, ok := cube[leftType][rightType][operator]; ok {
		// Si la operación está marcada como inválida, marca error
		if res == "error" {
			return "", fmt.Errorf("semántico: operación inválida %s %s %s", leftType, operator, rightType)
		}

		// Devuelve el tipo
		return res, nil
	}

	// Si no existe esa combinación en el cubo, marca error
	return "", fmt.Errorf("semántico: combinación no definida %s %s %s", leftType, operator, rightType)
}
