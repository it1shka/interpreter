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
	code := `
	@@@
	let counter = 100;
	for counter <= 100 {
		say counter;
		counter += 1
	};
	say "Finished!"
	`
	parser := core.NewParser(code)
	ast, err := parser.ParseProgram()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(PrettyPrint(ast))
}
