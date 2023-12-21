package main

import (
	"flag"
	"log"
	"os"

	graphviz "github.com/goccy/go-graphviz"
)

var debug = flag.Bool("d", false, "Debug")
var verbose = flag.Bool("v", false, "Verbose")

var graphFile = flag.String("graph-file", "", "Graph (DOT)")

func init() {
	flag.Parse()
}

func main() {
	g := graphviz.New()

	// Load a graph from a file
	if *graphFile != "" {
		data, err := os.ReadFile(*graphFile)
		if err != nil {
			log.Fatal(err)
		}

		graph, err := graphviz.ParseBytes(data)
		if err != nil {
			log.Fatal(err)
		}

		defer func() {
			if err := graph.Close(); err != nil {
				log.Fatal(err)
			}
			g.Close()
		}()
	} else {
		log.Fatalln("No graph file specified")
	}
}
