package main

import (
	"encoding/json"
	"fmt"
	"interpreter/core"
)

func PrettyPrint(structure interface{}) string {
	s, _ := json.MarshalIndent(structure, "", "\t")
	return string(s)
}

func main() {
	code := `           for true {
		-true;
	}
	`
	parser := core.NewParser(code)
	ast, err := parser.ParseProgram()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(PrettyPrint(ast))

	interpreter := core.NewInterpreter()
	res := interpreter.Interpret(ast)

	fmt.Println()
	fmt.Println(res)
}
