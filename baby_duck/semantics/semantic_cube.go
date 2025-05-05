package semantics

import (
	"fmt"
)

type SemanticCube map[string]map[string]map[string]string

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

func GetResultType(leftType, rightType, operator string) (string, error) {
	if res, ok := cube[leftType][rightType][operator]; ok {
		if res == "error" {
			return "", fmt.Errorf("semántico: operación inválida %s %s %s", leftType, operator, rightType)
		}
		return res, nil
	}
	return "", fmt.Errorf("semántico: combinación no definida %s %s %s", leftType, operator, rightType)
}
