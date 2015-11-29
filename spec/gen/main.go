package main

import (
	"fmt"
	"os"

	"h12.me/wipro/gen"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("gen (bnf | bnfj | goj | go | gof)")
		fmt.Println("bnf: from HTML to BNF")
		fmt.Println("bnfj: from BNF to BNF JSON")
		fmt.Println("goj: from BNF to Go JSON")
		fmt.Println("go: from BNF to Go")
		fmt.Println("gof: from BNF to Go funcs")
		fmt.Println("goe: from HTML to Go errors")
		return
	}
	file := os.Args[2]
	switch os.Args[1] {
	case "bnf":
		fromHTMLToBNF(file, os.Stdout)
	case "bnfj":
		bnf := gen.ParseBNF(file)
		fmt.Println(bnf.JSON())
	case "goj":
		bnf := gen.ParseBNF(file)
		fmt.Println(bnf.GoTypes().JSON())
	case "go":
		bnf := gen.ParseBNF(file)
		goTypes := bnf.GoTypes().RemoveDecl("RequestMessage")
		goTypes.Fprint(os.Stdout)
	case "gof":
		bnf := gen.ParseBNF(file)
		goTypes := bnf.GoTypes()
		goTypes.GoFuncs(os.Stdout)
	case "goe":
		genErrorCodes(file, os.Stdout)
	}
}
