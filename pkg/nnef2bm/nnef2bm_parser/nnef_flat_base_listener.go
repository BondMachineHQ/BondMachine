// Code generated from Nnef_flat.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Nnef_flat

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseNnef_flatListener is a complete listener for a parse tree produced by Nnef_flatParser.
type BaseNnef_flatListener struct{}

var _ Nnef_flatListener = &BaseNnef_flatListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseNnef_flatListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseNnef_flatListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseNnef_flatListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseNnef_flatListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterDocument is called when production document is entered.
func (s *BaseNnef_flatListener) EnterDocument(ctx *DocumentContext) {}

// ExitDocument is called when production document is exited.
func (s *BaseNnef_flatListener) ExitDocument(ctx *DocumentContext) {}

// EnterVersion is called when production version is entered.
func (s *BaseNnef_flatListener) EnterVersion(ctx *VersionContext) {}

// ExitVersion is called when production version is exited.
func (s *BaseNnef_flatListener) ExitVersion(ctx *VersionContext) {}

// EnterGraph_definition is called when production graph_definition is entered.
func (s *BaseNnef_flatListener) EnterGraph_definition(ctx *Graph_definitionContext) {}

// ExitGraph_definition is called when production graph_definition is exited.
func (s *BaseNnef_flatListener) ExitGraph_definition(ctx *Graph_definitionContext) {}

// EnterGraph_declaration is called when production graph_declaration is entered.
func (s *BaseNnef_flatListener) EnterGraph_declaration(ctx *Graph_declarationContext) {}

// ExitGraph_declaration is called when production graph_declaration is exited.
func (s *BaseNnef_flatListener) ExitGraph_declaration(ctx *Graph_declarationContext) {}

// EnterIdentifier_list is called when production identifier_list is entered.
func (s *BaseNnef_flatListener) EnterIdentifier_list(ctx *Identifier_listContext) {}

// ExitIdentifier_list is called when production identifier_list is exited.
func (s *BaseNnef_flatListener) ExitIdentifier_list(ctx *Identifier_listContext) {}

// EnterBody is called when production body is entered.
func (s *BaseNnef_flatListener) EnterBody(ctx *BodyContext) {}

// ExitBody is called when production body is exited.
func (s *BaseNnef_flatListener) ExitBody(ctx *BodyContext) {}

// EnterAssignment is called when production assignment is entered.
func (s *BaseNnef_flatListener) EnterAssignment(ctx *AssignmentContext) {}

// ExitAssignment is called when production assignment is exited.
func (s *BaseNnef_flatListener) ExitAssignment(ctx *AssignmentContext) {}

// EnterInvocation is called when production invocation is entered.
func (s *BaseNnef_flatListener) EnterInvocation(ctx *InvocationContext) {}

// ExitInvocation is called when production invocation is exited.
func (s *BaseNnef_flatListener) ExitInvocation(ctx *InvocationContext) {}

// EnterArgument_list is called when production argument_list is entered.
func (s *BaseNnef_flatListener) EnterArgument_list(ctx *Argument_listContext) {}

// ExitArgument_list is called when production argument_list is exited.
func (s *BaseNnef_flatListener) ExitArgument_list(ctx *Argument_listContext) {}

// EnterArgument is called when production argument is entered.
func (s *BaseNnef_flatListener) EnterArgument(ctx *ArgumentContext) {}

// ExitArgument is called when production argument is exited.
func (s *BaseNnef_flatListener) ExitArgument(ctx *ArgumentContext) {}

// EnterLvalue_expr is called when production lvalue_expr is entered.
func (s *BaseNnef_flatListener) EnterLvalue_expr(ctx *Lvalue_exprContext) {}

// ExitLvalue_expr is called when production lvalue_expr is exited.
func (s *BaseNnef_flatListener) ExitLvalue_expr(ctx *Lvalue_exprContext) {}

// EnterArray_lvalue_expr is called when production array_lvalue_expr is entered.
func (s *BaseNnef_flatListener) EnterArray_lvalue_expr(ctx *Array_lvalue_exprContext) {}

// ExitArray_lvalue_expr is called when production array_lvalue_expr is exited.
func (s *BaseNnef_flatListener) ExitArray_lvalue_expr(ctx *Array_lvalue_exprContext) {}

// EnterTuple_lvalue_expr is called when production tuple_lvalue_expr is entered.
func (s *BaseNnef_flatListener) EnterTuple_lvalue_expr(ctx *Tuple_lvalue_exprContext) {}

// ExitTuple_lvalue_expr is called when production tuple_lvalue_expr is exited.
func (s *BaseNnef_flatListener) ExitTuple_lvalue_expr(ctx *Tuple_lvalue_exprContext) {}

// EnterRvalue_expr is called when production rvalue_expr is entered.
func (s *BaseNnef_flatListener) EnterRvalue_expr(ctx *Rvalue_exprContext) {}

// ExitRvalue_expr is called when production rvalue_expr is exited.
func (s *BaseNnef_flatListener) ExitRvalue_expr(ctx *Rvalue_exprContext) {}

// EnterArray_rvalue_expr is called when production array_rvalue_expr is entered.
func (s *BaseNnef_flatListener) EnterArray_rvalue_expr(ctx *Array_rvalue_exprContext) {}

// ExitArray_rvalue_expr is called when production array_rvalue_expr is exited.
func (s *BaseNnef_flatListener) ExitArray_rvalue_expr(ctx *Array_rvalue_exprContext) {}

// EnterTuple_rvalue_expr is called when production tuple_rvalue_expr is entered.
func (s *BaseNnef_flatListener) EnterTuple_rvalue_expr(ctx *Tuple_rvalue_exprContext) {}

// ExitTuple_rvalue_expr is called when production tuple_rvalue_expr is exited.
func (s *BaseNnef_flatListener) ExitTuple_rvalue_expr(ctx *Tuple_rvalue_exprContext) {}

// EnterLiteral is called when production literal is entered.
func (s *BaseNnef_flatListener) EnterLiteral(ctx *LiteralContext) {}

// ExitLiteral is called when production literal is exited.
func (s *BaseNnef_flatListener) ExitLiteral(ctx *LiteralContext) {}
