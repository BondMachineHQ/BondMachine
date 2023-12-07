package main

import (
	"flag"
	"log"

	graphviz "github.com/goccy/go-graphviz"
	lua "github.com/yuin/gopher-lua"
)

var debug = flag.Bool("d", false, "Debug")
var verbose = flag.Bool("v", false, "Verbose")

var protocolFile = flag.String("protocol-file", "", "Protocol in (LUA)")
var initialFile = flag.String("initial-file", "", "Initial state (LUA)")
var graphFile = flag.String("graph-file", "", "Graph (DOT)")

func init() {
	flag.Parse()
}

func main() {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	l := lua.NewState()
	defer l.Close()
	if err := l.DoString(`print("hello")`); err != nil {
		panic(err)
	}

}
