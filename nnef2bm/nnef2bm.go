package nnef2bm

import (
	"fmt"
	parser "github.com/BondMachineHQ/BondMachine/nnef2bm/nnef2bm_parser"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type Nnef_flatListener struct {
	*parser.BaseNnef_flatListener
	// TODO other BM data
}

func (l *Nnef_flatListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	fmt.Println("Entering new rule")
	fmt.Println("\tcontext:", ctx.GetText())
}

func (l *Nnef_flatListener) EnterBody(c *parser.BodyContext) {
	//	fmt.Println("Entering Body")
}

func (l *Nnef_flatListener) ExitBody(c *parser.BodyContext) {
	//	fmt.Println("Exiting Body")
}

func NnefBuildBM(nnefModel string) {

	// Setup the input
	is := antlr.NewInputStream(nnefModel)

	// Create the Lexer
	lexer := parser.NewNnef_flatLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := parser.NewNnef_flatParser(stream)

	// Finally parse the expression (by walking the tree)
	var listener Nnef_flatListener
	antlr.ParseTreeWalkerDefault.Walk(&listener, p.Document())

}
