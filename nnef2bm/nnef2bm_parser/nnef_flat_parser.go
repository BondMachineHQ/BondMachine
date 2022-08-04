// Code generated from Nnef_flat.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Nnef_flat

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type Nnef_flatParser struct {
	*antlr.BaseParser
}

var nnef_flatParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func nnef_flatParserInit() {
	staticData := &nnef_flatParserStaticData
	staticData.literalNames = []string{
		"", "'version'", "';'", "'graph'", "'('", "')'", "'->'", "','", "'{'",
		"'}'", "'='", "'<'", "'>'", "'['", "']'",
	}
	staticData.symbolicNames = []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "TYPE_NAME",
		"IDENTIFIER", "FLOAT", "STRING_LITERAL", "LOGICAL_LITERAL", "NUMERIC_LITERAL",
		"WHITESPACE",
	}
	staticData.ruleNames = []string{
		"document", "version", "graph_definition", "graph_declaration", "identifier_list",
		"body", "assignment", "invocation", "argument_list", "argument", "lvalue_expr",
		"array_lvalue_expr", "tuple_lvalue_expr", "rvalue_expr", "array_rvalue_expr",
		"tuple_rvalue_expr", "literal",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 21, 164, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1,
		2, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 3, 1, 4, 1,
		4, 1, 4, 5, 4, 58, 8, 4, 10, 4, 12, 4, 61, 9, 4, 1, 5, 1, 5, 4, 5, 65,
		8, 5, 11, 5, 12, 5, 66, 1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7,
		1, 7, 1, 7, 1, 7, 3, 7, 80, 8, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1,
		8, 5, 8, 89, 8, 8, 10, 8, 12, 8, 92, 9, 8, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9,
		98, 8, 9, 1, 10, 1, 10, 3, 10, 102, 8, 10, 1, 11, 1, 11, 1, 11, 1, 11,
		5, 11, 108, 8, 11, 10, 11, 12, 11, 111, 9, 11, 1, 11, 1, 11, 1, 12, 1,
		12, 1, 12, 1, 12, 4, 12, 119, 8, 12, 11, 12, 12, 12, 120, 1, 12, 1, 12,
		1, 12, 1, 12, 1, 12, 4, 12, 128, 8, 12, 11, 12, 12, 12, 129, 3, 12, 132,
		8, 12, 1, 13, 1, 13, 1, 13, 1, 13, 3, 13, 138, 8, 13, 1, 14, 1, 14, 1,
		14, 1, 14, 5, 14, 144, 8, 14, 10, 14, 12, 14, 147, 9, 14, 1, 14, 1, 14,
		1, 15, 1, 15, 1, 15, 1, 15, 5, 15, 155, 8, 15, 10, 15, 12, 15, 158, 9,
		15, 1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 0, 0, 17, 0, 2, 4, 6, 8, 10, 12,
		14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 0, 1, 1, 0, 18, 20, 161, 0, 34,
		1, 0, 0, 0, 2, 37, 1, 0, 0, 0, 4, 41, 1, 0, 0, 0, 6, 44, 1, 0, 0, 0, 8,
		54, 1, 0, 0, 0, 10, 62, 1, 0, 0, 0, 12, 70, 1, 0, 0, 0, 14, 75, 1, 0, 0,
		0, 16, 85, 1, 0, 0, 0, 18, 97, 1, 0, 0, 0, 20, 101, 1, 0, 0, 0, 22, 103,
		1, 0, 0, 0, 24, 131, 1, 0, 0, 0, 26, 137, 1, 0, 0, 0, 28, 139, 1, 0, 0,
		0, 30, 150, 1, 0, 0, 0, 32, 161, 1, 0, 0, 0, 34, 35, 3, 2, 1, 0, 35, 36,
		3, 4, 2, 0, 36, 1, 1, 0, 0, 0, 37, 38, 5, 1, 0, 0, 38, 39, 5, 17, 0, 0,
		39, 40, 5, 2, 0, 0, 40, 3, 1, 0, 0, 0, 41, 42, 3, 6, 3, 0, 42, 43, 3, 10,
		5, 0, 43, 5, 1, 0, 0, 0, 44, 45, 5, 3, 0, 0, 45, 46, 5, 16, 0, 0, 46, 47,
		5, 4, 0, 0, 47, 48, 3, 8, 4, 0, 48, 49, 5, 5, 0, 0, 49, 50, 5, 6, 0, 0,
		50, 51, 5, 4, 0, 0, 51, 52, 3, 8, 4, 0, 52, 53, 5, 5, 0, 0, 53, 7, 1, 0,
		0, 0, 54, 59, 5, 16, 0, 0, 55, 56, 5, 7, 0, 0, 56, 58, 5, 16, 0, 0, 57,
		55, 1, 0, 0, 0, 58, 61, 1, 0, 0, 0, 59, 57, 1, 0, 0, 0, 59, 60, 1, 0, 0,
		0, 60, 9, 1, 0, 0, 0, 61, 59, 1, 0, 0, 0, 62, 64, 5, 8, 0, 0, 63, 65, 3,
		12, 6, 0, 64, 63, 1, 0, 0, 0, 65, 66, 1, 0, 0, 0, 66, 64, 1, 0, 0, 0, 66,
		67, 1, 0, 0, 0, 67, 68, 1, 0, 0, 0, 68, 69, 5, 9, 0, 0, 69, 11, 1, 0, 0,
		0, 70, 71, 3, 20, 10, 0, 71, 72, 5, 10, 0, 0, 72, 73, 3, 14, 7, 0, 73,
		74, 5, 2, 0, 0, 74, 13, 1, 0, 0, 0, 75, 79, 5, 16, 0, 0, 76, 77, 5, 11,
		0, 0, 77, 78, 5, 15, 0, 0, 78, 80, 5, 12, 0, 0, 79, 76, 1, 0, 0, 0, 79,
		80, 1, 0, 0, 0, 80, 81, 1, 0, 0, 0, 81, 82, 5, 4, 0, 0, 82, 83, 3, 16,
		8, 0, 83, 84, 5, 5, 0, 0, 84, 15, 1, 0, 0, 0, 85, 90, 3, 18, 9, 0, 86,
		87, 5, 7, 0, 0, 87, 89, 3, 18, 9, 0, 88, 86, 1, 0, 0, 0, 89, 92, 1, 0,
		0, 0, 90, 88, 1, 0, 0, 0, 90, 91, 1, 0, 0, 0, 91, 17, 1, 0, 0, 0, 92, 90,
		1, 0, 0, 0, 93, 98, 3, 26, 13, 0, 94, 95, 5, 16, 0, 0, 95, 96, 5, 10, 0,
		0, 96, 98, 3, 26, 13, 0, 97, 93, 1, 0, 0, 0, 97, 94, 1, 0, 0, 0, 98, 19,
		1, 0, 0, 0, 99, 102, 5, 16, 0, 0, 100, 102, 3, 22, 11, 0, 101, 99, 1, 0,
		0, 0, 101, 100, 1, 0, 0, 0, 102, 21, 1, 0, 0, 0, 103, 104, 5, 13, 0, 0,
		104, 109, 3, 20, 10, 0, 105, 106, 5, 7, 0, 0, 106, 108, 3, 20, 10, 0, 107,
		105, 1, 0, 0, 0, 108, 111, 1, 0, 0, 0, 109, 107, 1, 0, 0, 0, 109, 110,
		1, 0, 0, 0, 110, 112, 1, 0, 0, 0, 111, 109, 1, 0, 0, 0, 112, 113, 5, 14,
		0, 0, 113, 23, 1, 0, 0, 0, 114, 115, 5, 4, 0, 0, 115, 118, 3, 20, 10, 0,
		116, 117, 5, 7, 0, 0, 117, 119, 3, 20, 10, 0, 118, 116, 1, 0, 0, 0, 119,
		120, 1, 0, 0, 0, 120, 118, 1, 0, 0, 0, 120, 121, 1, 0, 0, 0, 121, 122,
		1, 0, 0, 0, 122, 123, 5, 5, 0, 0, 123, 132, 1, 0, 0, 0, 124, 127, 3, 20,
		10, 0, 125, 126, 5, 7, 0, 0, 126, 128, 3, 20, 10, 0, 127, 125, 1, 0, 0,
		0, 128, 129, 1, 0, 0, 0, 129, 127, 1, 0, 0, 0, 129, 130, 1, 0, 0, 0, 130,
		132, 1, 0, 0, 0, 131, 114, 1, 0, 0, 0, 131, 124, 1, 0, 0, 0, 132, 25, 1,
		0, 0, 0, 133, 138, 5, 16, 0, 0, 134, 138, 3, 32, 16, 0, 135, 138, 3, 28,
		14, 0, 136, 138, 3, 30, 15, 0, 137, 133, 1, 0, 0, 0, 137, 134, 1, 0, 0,
		0, 137, 135, 1, 0, 0, 0, 137, 136, 1, 0, 0, 0, 138, 27, 1, 0, 0, 0, 139,
		140, 5, 13, 0, 0, 140, 145, 3, 26, 13, 0, 141, 142, 5, 7, 0, 0, 142, 144,
		3, 26, 13, 0, 143, 141, 1, 0, 0, 0, 144, 147, 1, 0, 0, 0, 145, 143, 1,
		0, 0, 0, 145, 146, 1, 0, 0, 0, 146, 148, 1, 0, 0, 0, 147, 145, 1, 0, 0,
		0, 148, 149, 5, 14, 0, 0, 149, 29, 1, 0, 0, 0, 150, 151, 5, 4, 0, 0, 151,
		156, 3, 26, 13, 0, 152, 153, 5, 7, 0, 0, 153, 155, 3, 26, 13, 0, 154, 152,
		1, 0, 0, 0, 155, 158, 1, 0, 0, 0, 156, 154, 1, 0, 0, 0, 156, 157, 1, 0,
		0, 0, 157, 159, 1, 0, 0, 0, 158, 156, 1, 0, 0, 0, 159, 160, 5, 5, 0, 0,
		160, 31, 1, 0, 0, 0, 161, 162, 7, 0, 0, 0, 162, 33, 1, 0, 0, 0, 13, 59,
		66, 79, 90, 97, 101, 109, 120, 129, 131, 137, 145, 156,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// Nnef_flatParserInit initializes any static state used to implement Nnef_flatParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewNnef_flatParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func Nnef_flatParserInit() {
	staticData := &nnef_flatParserStaticData
	staticData.once.Do(nnef_flatParserInit)
}

// NewNnef_flatParser produces a new parser instance for the optional input antlr.TokenStream.
func NewNnef_flatParser(input antlr.TokenStream) *Nnef_flatParser {
	Nnef_flatParserInit()
	this := new(Nnef_flatParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &nnef_flatParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "Nnef_flat.g4"

	return this
}

// Nnef_flatParser tokens.
const (
	Nnef_flatParserEOF             = antlr.TokenEOF
	Nnef_flatParserT__0            = 1
	Nnef_flatParserT__1            = 2
	Nnef_flatParserT__2            = 3
	Nnef_flatParserT__3            = 4
	Nnef_flatParserT__4            = 5
	Nnef_flatParserT__5            = 6
	Nnef_flatParserT__6            = 7
	Nnef_flatParserT__7            = 8
	Nnef_flatParserT__8            = 9
	Nnef_flatParserT__9            = 10
	Nnef_flatParserT__10           = 11
	Nnef_flatParserT__11           = 12
	Nnef_flatParserT__12           = 13
	Nnef_flatParserT__13           = 14
	Nnef_flatParserTYPE_NAME       = 15
	Nnef_flatParserIDENTIFIER      = 16
	Nnef_flatParserFLOAT           = 17
	Nnef_flatParserSTRING_LITERAL  = 18
	Nnef_flatParserLOGICAL_LITERAL = 19
	Nnef_flatParserNUMERIC_LITERAL = 20
	Nnef_flatParserWHITESPACE      = 21
)

// Nnef_flatParser rules.
const (
	Nnef_flatParserRULE_document          = 0
	Nnef_flatParserRULE_version           = 1
	Nnef_flatParserRULE_graph_definition  = 2
	Nnef_flatParserRULE_graph_declaration = 3
	Nnef_flatParserRULE_identifier_list   = 4
	Nnef_flatParserRULE_body              = 5
	Nnef_flatParserRULE_assignment        = 6
	Nnef_flatParserRULE_invocation        = 7
	Nnef_flatParserRULE_argument_list     = 8
	Nnef_flatParserRULE_argument          = 9
	Nnef_flatParserRULE_lvalue_expr       = 10
	Nnef_flatParserRULE_array_lvalue_expr = 11
	Nnef_flatParserRULE_tuple_lvalue_expr = 12
	Nnef_flatParserRULE_rvalue_expr       = 13
	Nnef_flatParserRULE_array_rvalue_expr = 14
	Nnef_flatParserRULE_tuple_rvalue_expr = 15
	Nnef_flatParserRULE_literal           = 16
)

// IDocumentContext is an interface to support dynamic dispatch.
type IDocumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsDocumentContext differentiates from other interfaces.
	IsDocumentContext()
}

type DocumentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyDocumentContext() *DocumentContext {
	var p = new(DocumentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_document
	return p
}

func (*DocumentContext) IsDocumentContext() {}

func NewDocumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DocumentContext {
	var p = new(DocumentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_document

	return p
}

func (s *DocumentContext) GetParser() antlr.Parser { return s.parser }

func (s *DocumentContext) Version() IVersionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IVersionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IVersionContext)
}

func (s *DocumentContext) Graph_definition() IGraph_definitionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGraph_definitionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGraph_definitionContext)
}

func (s *DocumentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *DocumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *DocumentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterDocument(s)
	}
}

func (s *DocumentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitDocument(s)
	}
}

func (p *Nnef_flatParser) Document() (localctx IDocumentContext) {
	this := p
	_ = this

	localctx = NewDocumentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, Nnef_flatParserRULE_document)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(34)
		p.Version()
	}
	{
		p.SetState(35)
		p.Graph_definition()
	}

	return localctx
}

// IVersionContext is an interface to support dynamic dispatch.
type IVersionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsVersionContext differentiates from other interfaces.
	IsVersionContext()
}

type VersionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyVersionContext() *VersionContext {
	var p = new(VersionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_version
	return p
}

func (*VersionContext) IsVersionContext() {}

func NewVersionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VersionContext {
	var p = new(VersionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_version

	return p
}

func (s *VersionContext) GetParser() antlr.Parser { return s.parser }

func (s *VersionContext) FLOAT() antlr.TerminalNode {
	return s.GetToken(Nnef_flatParserFLOAT, 0)
}

func (s *VersionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *VersionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *VersionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterVersion(s)
	}
}

func (s *VersionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitVersion(s)
	}
}

func (p *Nnef_flatParser) Version() (localctx IVersionContext) {
	this := p
	_ = this

	localctx = NewVersionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, Nnef_flatParserRULE_version)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(37)
		p.Match(Nnef_flatParserT__0)
	}
	{
		p.SetState(38)
		p.Match(Nnef_flatParserFLOAT)
	}
	{
		p.SetState(39)
		p.Match(Nnef_flatParserT__1)
	}

	return localctx
}

// IGraph_definitionContext is an interface to support dynamic dispatch.
type IGraph_definitionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGraph_definitionContext differentiates from other interfaces.
	IsGraph_definitionContext()
}

type Graph_definitionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGraph_definitionContext() *Graph_definitionContext {
	var p = new(Graph_definitionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_graph_definition
	return p
}

func (*Graph_definitionContext) IsGraph_definitionContext() {}

func NewGraph_definitionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Graph_definitionContext {
	var p = new(Graph_definitionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_graph_definition

	return p
}

func (s *Graph_definitionContext) GetParser() antlr.Parser { return s.parser }

func (s *Graph_definitionContext) Graph_declaration() IGraph_declarationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGraph_declarationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGraph_declarationContext)
}

func (s *Graph_definitionContext) Body() IBodyContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IBodyContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IBodyContext)
}

func (s *Graph_definitionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Graph_definitionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Graph_definitionContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterGraph_definition(s)
	}
}

func (s *Graph_definitionContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitGraph_definition(s)
	}
}

func (p *Nnef_flatParser) Graph_definition() (localctx IGraph_definitionContext) {
	this := p
	_ = this

	localctx = NewGraph_definitionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, Nnef_flatParserRULE_graph_definition)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(41)
		p.Graph_declaration()
	}
	{
		p.SetState(42)
		p.Body()
	}

	return localctx
}

// IGraph_declarationContext is an interface to support dynamic dispatch.
type IGraph_declarationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGraph_declarationContext differentiates from other interfaces.
	IsGraph_declarationContext()
}

type Graph_declarationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGraph_declarationContext() *Graph_declarationContext {
	var p = new(Graph_declarationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_graph_declaration
	return p
}

func (*Graph_declarationContext) IsGraph_declarationContext() {}

func NewGraph_declarationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Graph_declarationContext {
	var p = new(Graph_declarationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_graph_declaration

	return p
}

func (s *Graph_declarationContext) GetParser() antlr.Parser { return s.parser }

func (s *Graph_declarationContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(Nnef_flatParserIDENTIFIER, 0)
}

func (s *Graph_declarationContext) AllIdentifier_list() []IIdentifier_listContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IIdentifier_listContext); ok {
			len++
		}
	}

	tst := make([]IIdentifier_listContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IIdentifier_listContext); ok {
			tst[i] = t.(IIdentifier_listContext)
			i++
		}
	}

	return tst
}

func (s *Graph_declarationContext) Identifier_list(i int) IIdentifier_listContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IIdentifier_listContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IIdentifier_listContext)
}

func (s *Graph_declarationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Graph_declarationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Graph_declarationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterGraph_declaration(s)
	}
}

func (s *Graph_declarationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitGraph_declaration(s)
	}
}

func (p *Nnef_flatParser) Graph_declaration() (localctx IGraph_declarationContext) {
	this := p
	_ = this

	localctx = NewGraph_declarationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, Nnef_flatParserRULE_graph_declaration)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(44)
		p.Match(Nnef_flatParserT__2)
	}
	{
		p.SetState(45)
		p.Match(Nnef_flatParserIDENTIFIER)
	}
	{
		p.SetState(46)
		p.Match(Nnef_flatParserT__3)
	}
	{
		p.SetState(47)
		p.Identifier_list()
	}
	{
		p.SetState(48)
		p.Match(Nnef_flatParserT__4)
	}
	{
		p.SetState(49)
		p.Match(Nnef_flatParserT__5)
	}
	{
		p.SetState(50)
		p.Match(Nnef_flatParserT__3)
	}
	{
		p.SetState(51)
		p.Identifier_list()
	}
	{
		p.SetState(52)
		p.Match(Nnef_flatParserT__4)
	}

	return localctx
}

// IIdentifier_listContext is an interface to support dynamic dispatch.
type IIdentifier_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsIdentifier_listContext differentiates from other interfaces.
	IsIdentifier_listContext()
}

type Identifier_listContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyIdentifier_listContext() *Identifier_listContext {
	var p = new(Identifier_listContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_identifier_list
	return p
}

func (*Identifier_listContext) IsIdentifier_listContext() {}

func NewIdentifier_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Identifier_listContext {
	var p = new(Identifier_listContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_identifier_list

	return p
}

func (s *Identifier_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Identifier_listContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(Nnef_flatParserIDENTIFIER)
}

func (s *Identifier_listContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(Nnef_flatParserIDENTIFIER, i)
}

func (s *Identifier_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Identifier_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Identifier_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterIdentifier_list(s)
	}
}

func (s *Identifier_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitIdentifier_list(s)
	}
}

func (p *Nnef_flatParser) Identifier_list() (localctx IIdentifier_listContext) {
	this := p
	_ = this

	localctx = NewIdentifier_listContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, Nnef_flatParserRULE_identifier_list)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(54)
		p.Match(Nnef_flatParserIDENTIFIER)
	}
	p.SetState(59)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == Nnef_flatParserT__6 {
		{
			p.SetState(55)
			p.Match(Nnef_flatParserT__6)
		}
		{
			p.SetState(56)
			p.Match(Nnef_flatParserIDENTIFIER)
		}

		p.SetState(61)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IBodyContext is an interface to support dynamic dispatch.
type IBodyContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsBodyContext differentiates from other interfaces.
	IsBodyContext()
}

type BodyContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyBodyContext() *BodyContext {
	var p = new(BodyContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_body
	return p
}

func (*BodyContext) IsBodyContext() {}

func NewBodyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *BodyContext {
	var p = new(BodyContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_body

	return p
}

func (s *BodyContext) GetParser() antlr.Parser { return s.parser }

func (s *BodyContext) AllAssignment() []IAssignmentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAssignmentContext); ok {
			len++
		}
	}

	tst := make([]IAssignmentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAssignmentContext); ok {
			tst[i] = t.(IAssignmentContext)
			i++
		}
	}

	return tst
}

func (s *BodyContext) Assignment(i int) IAssignmentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAssignmentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAssignmentContext)
}

func (s *BodyContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *BodyContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *BodyContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterBody(s)
	}
}

func (s *BodyContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitBody(s)
	}
}

func (p *Nnef_flatParser) Body() (localctx IBodyContext) {
	this := p
	_ = this

	localctx = NewBodyContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, Nnef_flatParserRULE_body)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(62)
		p.Match(Nnef_flatParserT__7)
	}
	p.SetState(64)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for ok := true; ok; ok = _la == Nnef_flatParserT__12 || _la == Nnef_flatParserIDENTIFIER {
		{
			p.SetState(63)
			p.Assignment()
		}

		p.SetState(66)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(68)
		p.Match(Nnef_flatParserT__8)
	}

	return localctx
}

// IAssignmentContext is an interface to support dynamic dispatch.
type IAssignmentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAssignmentContext differentiates from other interfaces.
	IsAssignmentContext()
}

type AssignmentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAssignmentContext() *AssignmentContext {
	var p = new(AssignmentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_assignment
	return p
}

func (*AssignmentContext) IsAssignmentContext() {}

func NewAssignmentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AssignmentContext {
	var p = new(AssignmentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_assignment

	return p
}

func (s *AssignmentContext) GetParser() antlr.Parser { return s.parser }

func (s *AssignmentContext) Lvalue_expr() ILvalue_exprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILvalue_exprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILvalue_exprContext)
}

func (s *AssignmentContext) Invocation() IInvocationContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IInvocationContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IInvocationContext)
}

func (s *AssignmentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AssignmentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AssignmentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterAssignment(s)
	}
}

func (s *AssignmentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitAssignment(s)
	}
}

func (p *Nnef_flatParser) Assignment() (localctx IAssignmentContext) {
	this := p
	_ = this

	localctx = NewAssignmentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, Nnef_flatParserRULE_assignment)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(70)
		p.Lvalue_expr()
	}
	{
		p.SetState(71)
		p.Match(Nnef_flatParserT__9)
	}
	{
		p.SetState(72)
		p.Invocation()
	}
	{
		p.SetState(73)
		p.Match(Nnef_flatParserT__1)
	}

	return localctx
}

// IInvocationContext is an interface to support dynamic dispatch.
type IInvocationContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsInvocationContext differentiates from other interfaces.
	IsInvocationContext()
}

type InvocationContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyInvocationContext() *InvocationContext {
	var p = new(InvocationContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_invocation
	return p
}

func (*InvocationContext) IsInvocationContext() {}

func NewInvocationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *InvocationContext {
	var p = new(InvocationContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_invocation

	return p
}

func (s *InvocationContext) GetParser() antlr.Parser { return s.parser }

func (s *InvocationContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(Nnef_flatParserIDENTIFIER, 0)
}

func (s *InvocationContext) Argument_list() IArgument_listContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArgument_listContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArgument_listContext)
}

func (s *InvocationContext) TYPE_NAME() antlr.TerminalNode {
	return s.GetToken(Nnef_flatParserTYPE_NAME, 0)
}

func (s *InvocationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *InvocationContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *InvocationContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterInvocation(s)
	}
}

func (s *InvocationContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitInvocation(s)
	}
}

func (p *Nnef_flatParser) Invocation() (localctx IInvocationContext) {
	this := p
	_ = this

	localctx = NewInvocationContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, Nnef_flatParserRULE_invocation)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(75)
		p.Match(Nnef_flatParserIDENTIFIER)
	}
	p.SetState(79)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == Nnef_flatParserT__10 {
		{
			p.SetState(76)
			p.Match(Nnef_flatParserT__10)
		}
		{
			p.SetState(77)
			p.Match(Nnef_flatParserTYPE_NAME)
		}
		{
			p.SetState(78)
			p.Match(Nnef_flatParserT__11)
		}

	}
	{
		p.SetState(81)
		p.Match(Nnef_flatParserT__3)
	}
	{
		p.SetState(82)
		p.Argument_list()
	}
	{
		p.SetState(83)
		p.Match(Nnef_flatParserT__4)
	}

	return localctx
}

// IArgument_listContext is an interface to support dynamic dispatch.
type IArgument_listContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArgument_listContext differentiates from other interfaces.
	IsArgument_listContext()
}

type Argument_listContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArgument_listContext() *Argument_listContext {
	var p = new(Argument_listContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_argument_list
	return p
}

func (*Argument_listContext) IsArgument_listContext() {}

func NewArgument_listContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Argument_listContext {
	var p = new(Argument_listContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_argument_list

	return p
}

func (s *Argument_listContext) GetParser() antlr.Parser { return s.parser }

func (s *Argument_listContext) AllArgument() []IArgumentContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IArgumentContext); ok {
			len++
		}
	}

	tst := make([]IArgumentContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IArgumentContext); ok {
			tst[i] = t.(IArgumentContext)
			i++
		}
	}

	return tst
}

func (s *Argument_listContext) Argument(i int) IArgumentContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArgumentContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArgumentContext)
}

func (s *Argument_listContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Argument_listContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Argument_listContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterArgument_list(s)
	}
}

func (s *Argument_listContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitArgument_list(s)
	}
}

func (p *Nnef_flatParser) Argument_list() (localctx IArgument_listContext) {
	this := p
	_ = this

	localctx = NewArgument_listContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, Nnef_flatParserRULE_argument_list)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(85)
		p.Argument()
	}
	p.SetState(90)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == Nnef_flatParserT__6 {
		{
			p.SetState(86)
			p.Match(Nnef_flatParserT__6)
		}
		{
			p.SetState(87)
			p.Argument()
		}

		p.SetState(92)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IArgumentContext is an interface to support dynamic dispatch.
type IArgumentContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArgumentContext differentiates from other interfaces.
	IsArgumentContext()
}

type ArgumentContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArgumentContext() *ArgumentContext {
	var p = new(ArgumentContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_argument
	return p
}

func (*ArgumentContext) IsArgumentContext() {}

func NewArgumentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ArgumentContext {
	var p = new(ArgumentContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_argument

	return p
}

func (s *ArgumentContext) GetParser() antlr.Parser { return s.parser }

func (s *ArgumentContext) Rvalue_expr() IRvalue_exprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRvalue_exprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRvalue_exprContext)
}

func (s *ArgumentContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(Nnef_flatParserIDENTIFIER, 0)
}

func (s *ArgumentContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ArgumentContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ArgumentContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterArgument(s)
	}
}

func (s *ArgumentContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitArgument(s)
	}
}

func (p *Nnef_flatParser) Argument() (localctx IArgumentContext) {
	this := p
	_ = this

	localctx = NewArgumentContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, Nnef_flatParserRULE_argument)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(97)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 4, p.GetParserRuleContext()) {
	case 1:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(93)
			p.Rvalue_expr()
		}

	case 2:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(94)
			p.Match(Nnef_flatParserIDENTIFIER)
		}
		{
			p.SetState(95)
			p.Match(Nnef_flatParserT__9)
		}
		{
			p.SetState(96)
			p.Rvalue_expr()
		}

	}

	return localctx
}

// ILvalue_exprContext is an interface to support dynamic dispatch.
type ILvalue_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLvalue_exprContext differentiates from other interfaces.
	IsLvalue_exprContext()
}

type Lvalue_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLvalue_exprContext() *Lvalue_exprContext {
	var p = new(Lvalue_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_lvalue_expr
	return p
}

func (*Lvalue_exprContext) IsLvalue_exprContext() {}

func NewLvalue_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Lvalue_exprContext {
	var p = new(Lvalue_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_lvalue_expr

	return p
}

func (s *Lvalue_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Lvalue_exprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(Nnef_flatParserIDENTIFIER, 0)
}

func (s *Lvalue_exprContext) Array_lvalue_expr() IArray_lvalue_exprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArray_lvalue_exprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArray_lvalue_exprContext)
}

func (s *Lvalue_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Lvalue_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Lvalue_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterLvalue_expr(s)
	}
}

func (s *Lvalue_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitLvalue_expr(s)
	}
}

func (p *Nnef_flatParser) Lvalue_expr() (localctx ILvalue_exprContext) {
	this := p
	_ = this

	localctx = NewLvalue_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, Nnef_flatParserRULE_lvalue_expr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(101)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case Nnef_flatParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(99)
			p.Match(Nnef_flatParserIDENTIFIER)
		}

	case Nnef_flatParserT__12:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(100)
			p.Array_lvalue_expr()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IArray_lvalue_exprContext is an interface to support dynamic dispatch.
type IArray_lvalue_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArray_lvalue_exprContext differentiates from other interfaces.
	IsArray_lvalue_exprContext()
}

type Array_lvalue_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArray_lvalue_exprContext() *Array_lvalue_exprContext {
	var p = new(Array_lvalue_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_array_lvalue_expr
	return p
}

func (*Array_lvalue_exprContext) IsArray_lvalue_exprContext() {}

func NewArray_lvalue_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Array_lvalue_exprContext {
	var p = new(Array_lvalue_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_array_lvalue_expr

	return p
}

func (s *Array_lvalue_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Array_lvalue_exprContext) AllLvalue_expr() []ILvalue_exprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILvalue_exprContext); ok {
			len++
		}
	}

	tst := make([]ILvalue_exprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILvalue_exprContext); ok {
			tst[i] = t.(ILvalue_exprContext)
			i++
		}
	}

	return tst
}

func (s *Array_lvalue_exprContext) Lvalue_expr(i int) ILvalue_exprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILvalue_exprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILvalue_exprContext)
}

func (s *Array_lvalue_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Array_lvalue_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Array_lvalue_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterArray_lvalue_expr(s)
	}
}

func (s *Array_lvalue_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitArray_lvalue_expr(s)
	}
}

func (p *Nnef_flatParser) Array_lvalue_expr() (localctx IArray_lvalue_exprContext) {
	this := p
	_ = this

	localctx = NewArray_lvalue_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, Nnef_flatParserRULE_array_lvalue_expr)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(103)
		p.Match(Nnef_flatParserT__12)
	}
	{
		p.SetState(104)
		p.Lvalue_expr()
	}
	p.SetState(109)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == Nnef_flatParserT__6 {
		{
			p.SetState(105)
			p.Match(Nnef_flatParserT__6)
		}
		{
			p.SetState(106)
			p.Lvalue_expr()
		}

		p.SetState(111)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(112)
		p.Match(Nnef_flatParserT__13)
	}

	return localctx
}

// ITuple_lvalue_exprContext is an interface to support dynamic dispatch.
type ITuple_lvalue_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTuple_lvalue_exprContext differentiates from other interfaces.
	IsTuple_lvalue_exprContext()
}

type Tuple_lvalue_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTuple_lvalue_exprContext() *Tuple_lvalue_exprContext {
	var p = new(Tuple_lvalue_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_tuple_lvalue_expr
	return p
}

func (*Tuple_lvalue_exprContext) IsTuple_lvalue_exprContext() {}

func NewTuple_lvalue_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Tuple_lvalue_exprContext {
	var p = new(Tuple_lvalue_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_tuple_lvalue_expr

	return p
}

func (s *Tuple_lvalue_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Tuple_lvalue_exprContext) AllLvalue_expr() []ILvalue_exprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILvalue_exprContext); ok {
			len++
		}
	}

	tst := make([]ILvalue_exprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILvalue_exprContext); ok {
			tst[i] = t.(ILvalue_exprContext)
			i++
		}
	}

	return tst
}

func (s *Tuple_lvalue_exprContext) Lvalue_expr(i int) ILvalue_exprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILvalue_exprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILvalue_exprContext)
}

func (s *Tuple_lvalue_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Tuple_lvalue_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Tuple_lvalue_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterTuple_lvalue_expr(s)
	}
}

func (s *Tuple_lvalue_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitTuple_lvalue_expr(s)
	}
}

func (p *Nnef_flatParser) Tuple_lvalue_expr() (localctx ITuple_lvalue_exprContext) {
	this := p
	_ = this

	localctx = NewTuple_lvalue_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 24, Nnef_flatParserRULE_tuple_lvalue_expr)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(131)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case Nnef_flatParserT__3:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(114)
			p.Match(Nnef_flatParserT__3)
		}
		{
			p.SetState(115)
			p.Lvalue_expr()
		}
		p.SetState(118)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == Nnef_flatParserT__6 {
			{
				p.SetState(116)
				p.Match(Nnef_flatParserT__6)
			}
			{
				p.SetState(117)
				p.Lvalue_expr()
			}

			p.SetState(120)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(122)
			p.Match(Nnef_flatParserT__4)
		}

	case Nnef_flatParserT__12, Nnef_flatParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(124)
			p.Lvalue_expr()
		}
		p.SetState(127)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == Nnef_flatParserT__6 {
			{
				p.SetState(125)
				p.Match(Nnef_flatParserT__6)
			}
			{
				p.SetState(126)
				p.Lvalue_expr()
			}

			p.SetState(129)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IRvalue_exprContext is an interface to support dynamic dispatch.
type IRvalue_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsRvalue_exprContext differentiates from other interfaces.
	IsRvalue_exprContext()
}

type Rvalue_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyRvalue_exprContext() *Rvalue_exprContext {
	var p = new(Rvalue_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_rvalue_expr
	return p
}

func (*Rvalue_exprContext) IsRvalue_exprContext() {}

func NewRvalue_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Rvalue_exprContext {
	var p = new(Rvalue_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_rvalue_expr

	return p
}

func (s *Rvalue_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Rvalue_exprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(Nnef_flatParserIDENTIFIER, 0)
}

func (s *Rvalue_exprContext) Literal() ILiteralContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILiteralContext)
}

func (s *Rvalue_exprContext) Array_rvalue_expr() IArray_rvalue_exprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IArray_rvalue_exprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IArray_rvalue_exprContext)
}

func (s *Rvalue_exprContext) Tuple_rvalue_expr() ITuple_rvalue_exprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITuple_rvalue_exprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITuple_rvalue_exprContext)
}

func (s *Rvalue_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Rvalue_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Rvalue_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterRvalue_expr(s)
	}
}

func (s *Rvalue_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitRvalue_expr(s)
	}
}

func (p *Nnef_flatParser) Rvalue_expr() (localctx IRvalue_exprContext) {
	this := p
	_ = this

	localctx = NewRvalue_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, Nnef_flatParserRULE_rvalue_expr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(137)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case Nnef_flatParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(133)
			p.Match(Nnef_flatParserIDENTIFIER)
		}

	case Nnef_flatParserSTRING_LITERAL, Nnef_flatParserLOGICAL_LITERAL, Nnef_flatParserNUMERIC_LITERAL:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(134)
			p.Literal()
		}

	case Nnef_flatParserT__12:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(135)
			p.Array_rvalue_expr()
		}

	case Nnef_flatParserT__3:
		p.EnterOuterAlt(localctx, 4)
		{
			p.SetState(136)
			p.Tuple_rvalue_expr()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IArray_rvalue_exprContext is an interface to support dynamic dispatch.
type IArray_rvalue_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsArray_rvalue_exprContext differentiates from other interfaces.
	IsArray_rvalue_exprContext()
}

type Array_rvalue_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyArray_rvalue_exprContext() *Array_rvalue_exprContext {
	var p = new(Array_rvalue_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_array_rvalue_expr
	return p
}

func (*Array_rvalue_exprContext) IsArray_rvalue_exprContext() {}

func NewArray_rvalue_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Array_rvalue_exprContext {
	var p = new(Array_rvalue_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_array_rvalue_expr

	return p
}

func (s *Array_rvalue_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Array_rvalue_exprContext) AllRvalue_expr() []IRvalue_exprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRvalue_exprContext); ok {
			len++
		}
	}

	tst := make([]IRvalue_exprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRvalue_exprContext); ok {
			tst[i] = t.(IRvalue_exprContext)
			i++
		}
	}

	return tst
}

func (s *Array_rvalue_exprContext) Rvalue_expr(i int) IRvalue_exprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRvalue_exprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRvalue_exprContext)
}

func (s *Array_rvalue_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Array_rvalue_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Array_rvalue_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterArray_rvalue_expr(s)
	}
}

func (s *Array_rvalue_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitArray_rvalue_expr(s)
	}
}

func (p *Nnef_flatParser) Array_rvalue_expr() (localctx IArray_rvalue_exprContext) {
	this := p
	_ = this

	localctx = NewArray_rvalue_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, Nnef_flatParserRULE_array_rvalue_expr)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(139)
		p.Match(Nnef_flatParserT__12)
	}
	{
		p.SetState(140)
		p.Rvalue_expr()
	}
	p.SetState(145)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == Nnef_flatParserT__6 {
		{
			p.SetState(141)
			p.Match(Nnef_flatParserT__6)
		}
		{
			p.SetState(142)
			p.Rvalue_expr()
		}

		p.SetState(147)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(148)
		p.Match(Nnef_flatParserT__13)
	}

	return localctx
}

// ITuple_rvalue_exprContext is an interface to support dynamic dispatch.
type ITuple_rvalue_exprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTuple_rvalue_exprContext differentiates from other interfaces.
	IsTuple_rvalue_exprContext()
}

type Tuple_rvalue_exprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTuple_rvalue_exprContext() *Tuple_rvalue_exprContext {
	var p = new(Tuple_rvalue_exprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_tuple_rvalue_expr
	return p
}

func (*Tuple_rvalue_exprContext) IsTuple_rvalue_exprContext() {}

func NewTuple_rvalue_exprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *Tuple_rvalue_exprContext {
	var p = new(Tuple_rvalue_exprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_tuple_rvalue_expr

	return p
}

func (s *Tuple_rvalue_exprContext) GetParser() antlr.Parser { return s.parser }

func (s *Tuple_rvalue_exprContext) AllRvalue_expr() []IRvalue_exprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IRvalue_exprContext); ok {
			len++
		}
	}

	tst := make([]IRvalue_exprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IRvalue_exprContext); ok {
			tst[i] = t.(IRvalue_exprContext)
			i++
		}
	}

	return tst
}

func (s *Tuple_rvalue_exprContext) Rvalue_expr(i int) IRvalue_exprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IRvalue_exprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IRvalue_exprContext)
}

func (s *Tuple_rvalue_exprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *Tuple_rvalue_exprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *Tuple_rvalue_exprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterTuple_rvalue_expr(s)
	}
}

func (s *Tuple_rvalue_exprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitTuple_rvalue_expr(s)
	}
}

func (p *Nnef_flatParser) Tuple_rvalue_expr() (localctx ITuple_rvalue_exprContext) {
	this := p
	_ = this

	localctx = NewTuple_rvalue_exprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 30, Nnef_flatParserRULE_tuple_rvalue_expr)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(150)
		p.Match(Nnef_flatParserT__3)
	}
	{
		p.SetState(151)
		p.Rvalue_expr()
	}
	p.SetState(156)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == Nnef_flatParserT__6 {
		{
			p.SetState(152)
			p.Match(Nnef_flatParserT__6)
		}
		{
			p.SetState(153)
			p.Rvalue_expr()
		}

		p.SetState(158)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}
	{
		p.SetState(159)
		p.Match(Nnef_flatParserT__4)
	}

	return localctx
}

// ILiteralContext is an interface to support dynamic dispatch.
type ILiteralContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLiteralContext differentiates from other interfaces.
	IsLiteralContext()
}

type LiteralContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralContext() *LiteralContext {
	var p = new(LiteralContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = Nnef_flatParserRULE_literal
	return p
}

func (*LiteralContext) IsLiteralContext() {}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	var p = new(LiteralContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = Nnef_flatParserRULE_literal

	return p
}

func (s *LiteralContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralContext) NUMERIC_LITERAL() antlr.TerminalNode {
	return s.GetToken(Nnef_flatParserNUMERIC_LITERAL, 0)
}

func (s *LiteralContext) STRING_LITERAL() antlr.TerminalNode {
	return s.GetToken(Nnef_flatParserSTRING_LITERAL, 0)
}

func (s *LiteralContext) LOGICAL_LITERAL() antlr.TerminalNode {
	return s.GetToken(Nnef_flatParserLOGICAL_LITERAL, 0)
}

func (s *LiteralContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LiteralContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.EnterLiteral(s)
	}
}

func (s *LiteralContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(Nnef_flatListener); ok {
		listenerT.ExitLiteral(s)
	}
}

func (p *Nnef_flatParser) Literal() (localctx ILiteralContext) {
	this := p
	_ = this

	localctx = NewLiteralContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 32, Nnef_flatParserRULE_literal)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(161)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<Nnef_flatParserSTRING_LITERAL)|(1<<Nnef_flatParserLOGICAL_LITERAL)|(1<<Nnef_flatParserNUMERIC_LITERAL))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}
