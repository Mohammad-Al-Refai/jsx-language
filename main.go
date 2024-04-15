package main

import (
	"encoding/json"
	"os"

	"m.shebli.refaai/ht/lexer"
	"m.shebli.refaai/ht/runtime"
)

func main() {
	lex := lexer.Lexer{}
	file, err := os.Open("./examples/test.ht")
	if err != nil {
		panic(err)
	}
	tokens := lex.LoadFileReader(file)
	// fmt.Printf("%+v\n", tokens)
	ast := lexer.NewAST(tokens)
	program := ast.ProduceAST()
	// fmt.Printf("%+v\n", program)
	program_ast, err := json.Marshal(program)
	if err != nil {
		panic(err)
	}
	os.WriteFile("AST.json", program_ast, 0777)
	interpreter := runtime.NewInterpreter(program)
	interpreter.Run()
}
