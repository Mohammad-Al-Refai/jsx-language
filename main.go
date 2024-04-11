package main

import (
	"encoding/json"
	"fmt"
	"os"

	"m.shebli.refaai/ht/lexer"
)

func main() {
	lex := lexer.Lexer{}
	file, err := os.Open("code.html")
	if err != nil {
		panic(err)
	}

	tokens := lex.LoadFileReader(file)
	fmt.Printf("%+v\n", tokens)
	ast := lexer.NewAST(tokens)
	program := ast.ProduceAST()
	// fmt.Printf("%+v\n", program)
	data, err := json.Marshal(program)
	if err != nil {
		panic(err)
	}
	os.WriteFile("AST.json", data, 0777)

}
