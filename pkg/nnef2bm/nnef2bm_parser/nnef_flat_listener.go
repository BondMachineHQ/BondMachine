// Code generated from Nnef_flat.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Nnef_flat

import "github.com/antlr/antlr4/runtime/Go/antlr"

// Nnef_flatListener is a complete listener for a parse tree produced by Nnef_flatParser.
type Nnef_flatListener interface {
	antlr.ParseTreeListener

	// EnterDocument is called when entering the document production.
	EnterDocument(c *DocumentContext)

	// EnterVersion is called when entering the version production.
	EnterVersion(c *VersionContext)

	// EnterGraph_definition is called when entering the graph_definition production.
	EnterGraph_definition(c *Graph_definitionContext)

	// EnterGraph_declaration is called when entering the graph_declaration production.
	EnterGraph_declaration(c *Graph_declarationContext)

	// EnterIdentifier_list is called when entering the identifier_list production.
	EnterIdentifier_list(c *Identifier_listContext)

	// EnterBody is called when entering the body production.
	EnterBody(c *BodyContext)

	// EnterAssignment is called when entering the assignment production.
	EnterAssignment(c *AssignmentContext)

	// EnterInvocation is called when entering the invocation production.
	EnterInvocation(c *InvocationContext)

	// EnterArgument_list is called when entering the argument_list production.
	EnterArgument_list(c *Argument_listContext)

	// EnterArgument is called when entering the argument production.
	EnterArgument(c *ArgumentContext)

	// EnterLvalue_expr is called when entering the lvalue_expr production.
	EnterLvalue_expr(c *Lvalue_exprContext)

	// EnterArray_lvalue_expr is called when entering the array_lvalue_expr production.
	EnterArray_lvalue_expr(c *Array_lvalue_exprContext)

	// EnterTuple_lvalue_expr is called when entering the tuple_lvalue_expr production.
	EnterTuple_lvalue_expr(c *Tuple_lvalue_exprContext)

	// EnterRvalue_expr is called when entering the rvalue_expr production.
	EnterRvalue_expr(c *Rvalue_exprContext)

	// EnterArray_rvalue_expr is called when entering the array_rvalue_expr production.
	EnterArray_rvalue_expr(c *Array_rvalue_exprContext)

	// EnterTuple_rvalue_expr is called when entering the tuple_rvalue_expr production.
	EnterTuple_rvalue_expr(c *Tuple_rvalue_exprContext)

	// EnterLiteral is called when entering the literal production.
	EnterLiteral(c *LiteralContext)

	// ExitDocument is called when exiting the document production.
	ExitDocument(c *DocumentContext)

	// ExitVersion is called when exiting the version production.
	ExitVersion(c *VersionContext)

	// ExitGraph_definition is called when exiting the graph_definition production.
	ExitGraph_definition(c *Graph_definitionContext)

	// ExitGraph_declaration is called when exiting the graph_declaration production.
	ExitGraph_declaration(c *Graph_declarationContext)

	// ExitIdentifier_list is called when exiting the identifier_list production.
	ExitIdentifier_list(c *Identifier_listContext)

	// ExitBody is called when exiting the body production.
	ExitBody(c *BodyContext)

	// ExitAssignment is called when exiting the assignment production.
	ExitAssignment(c *AssignmentContext)

	// ExitInvocation is called when exiting the invocation production.
	ExitInvocation(c *InvocationContext)

	// ExitArgument_list is called when exiting the argument_list production.
	ExitArgument_list(c *Argument_listContext)

	// ExitArgument is called when exiting the argument production.
	ExitArgument(c *ArgumentContext)

	// ExitLvalue_expr is called when exiting the lvalue_expr production.
	ExitLvalue_expr(c *Lvalue_exprContext)

	// ExitArray_lvalue_expr is called when exiting the array_lvalue_expr production.
	ExitArray_lvalue_expr(c *Array_lvalue_exprContext)

	// ExitTuple_lvalue_expr is called when exiting the tuple_lvalue_expr production.
	ExitTuple_lvalue_expr(c *Tuple_lvalue_exprContext)

	// ExitRvalue_expr is called when exiting the rvalue_expr production.
	ExitRvalue_expr(c *Rvalue_exprContext)

	// ExitArray_rvalue_expr is called when exiting the array_rvalue_expr production.
	ExitArray_rvalue_expr(c *Array_rvalue_exprContext)

	// ExitTuple_rvalue_expr is called when exiting the tuple_rvalue_expr production.
	ExitTuple_rvalue_expr(c *Tuple_rvalue_exprContext)

	// ExitLiteral is called when exiting the literal production.
	ExitLiteral(c *LiteralContext)
}
