package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var convertTo = flag.String("convert", "", "Convert to type")
var castTo = flag.String("cast", "", "Cast to type")
var showAs = flag.String("show", "native", "Show as (native, hex, bin)")

var withSize = flag.Bool("with-size", false, "With size")

var useFiles = flag.Bool("use-files", false, "Load files instead of command line argument")

var serve = flag.Bool("serve", false, "Serve as REST API")

// Custom types
var linearDataRange = flag.String("linear-data-range", "", "Load a linear data range file (with the syntax index,filename)")

func init() {
	flag.Parse()
	// Load custom types
	if *linearDataRange != "" {
		if err := bmnumbers.LoadLinearDataRangesFromFile(*linearDataRange); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {

	if *serve {
		bmnumbers.Serve()
	} else {

		var newType bmnumbers.BMNumberType
		if *castTo != "" && *convertTo != "" {
			log.Fatal("Error: Cannot cast and convert at the same time")
		}

		switch {
		case *convertTo != "":
			if _, err := bmnumbers.EventuallyCreateType(*convertTo, nil); err != nil {
				log.Fatal(err)
			}
			if v := bmnumbers.GetType(*convertTo); v == nil {
				log.Fatal("Error: Unknown type")
			} else {
				newType = v
				// fmt.Println(v)
			}
		case *castTo != "":
			if _, err := bmnumbers.EventuallyCreateType(*castTo, nil); err != nil {
				log.Fatal(err)
			}
			if v := bmnumbers.GetType(*castTo); v == nil {
				log.Fatal("Error: Unknown type")
			} else {
				newType = v
				// fmt.Println(v)
			}
		}

		for _, argTo := range flag.Args() {
			if *useFiles {
				fmt.Println("Load files")
			} else {
				if output, err := bmnumbers.ImportString(argTo); err != nil {
					fmt.Println("Error: ", err)
				} else {

					if *convertTo != "" {
						if err := newType.Convert(output); err != nil {
							fmt.Println("Error: ", err)
						}
					}

					if *castTo != "" {
						if err := bmnumbers.CastType(output, newType); err != nil {
							fmt.Println("Error: ", err)
						}
					}

					switch *showAs {
					case "native":
						if value, err := output.ExportString(); err != nil {
							fmt.Println("Error: ", err)
						} else {
							fmt.Println(value)
						}
					case "bin":
						if value, err := output.ExportBinary(*withSize); err != nil {
							fmt.Println("Error: ", err)
						} else {
							fmt.Println(value)
						}
					case "unsigned":
						if value, err := output.ExportUint64(); err != nil {
							fmt.Println("Error: ", err)
						} else {
							fmt.Println(value)
						}
					default:
						fmt.Println("Error: Unknown visualization format (native, bin, unsigned)")
					}
				}
			}
		}
	}
}
