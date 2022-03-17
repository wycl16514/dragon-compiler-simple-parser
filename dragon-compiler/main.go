package main

import (
	"fmt"
	"io"
	"lexer"
	"simple_parser"
)

func main() {
	source := "(1+(2+3))"
	my_lexer := lexer.NewLexer(source)
	parser := simple_parser.NewSimpleParser(my_lexer)
	err := parser.Parse()
	if err != nil && err != io.EOF {
		fmt.Println("source error: ", err)
	} else {
		fmt.Println("source is legal expression")
	}
}
