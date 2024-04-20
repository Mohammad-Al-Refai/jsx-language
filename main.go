package main

import (
	"fmt"
	"os"

	"m.shebli.refaai/ht/lexer"
	"m.shebli.refaai/ht/runtime"
)

func main() {
	file, err := RequestFile()
	// file, err := os.Open("./examples/for-loop.sog")
	if err != nil {
		fmt.Println(err)
		return
	}
	lex := lexer.Lexer{}
	tokens := lex.LoadFileReader(file)
	ast := lexer.NewAST(tokens)
	program := ast.ProduceAST()
	// program_ast, err := json.Marshal(program)
	// if err != nil {
	// 	panic(err)
	// }
	// os.WriteFile("AST.json", program_ast, 0777)
	interpreter := runtime.NewInterpreter(program)
	interpreter.Run()
}
func RequestFile() (*os.File, error) {
	args := os.Args
	if len(args) == 1 {
		return nil, fmt.Errorf("[Error]: Missing file path")
	}
	filePath := args[1]
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("[Error]: '%v' is not found", filePath)
	}
	return file, nil
}
