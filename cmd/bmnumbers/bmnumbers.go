package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/BondMachineHQ/BondMachine/pkg/bmnumbers"
)

var verbose = flag.Bool("v", false, "Verbose")
var debug = flag.Bool("d", false, "Verbose")

var convertTo = flag.String("convert-to", "", "Convert to type")
var overrideTo = flag.String("override-to", "", "Override type")
var dumpAs = flag.String("dump-as", "native", "Dump as (native, hex, bin)")

var withSize = flag.Bool("with-size", false, "With size")

var useFiles = flag.Bool("use-files", false, "Load files instead of command line argument")

var serve = flag.Bool("serve", false, "Serve as REST API")

func init() {
	flag.Parse()
}

func main() {

	if *serve {
		bmnumbers.Serve()
	} else {

		var newType bmnumbers.BMNumberType
		if *overrideTo != "" && *convertTo != "" {
			log.Fatal("Error: Cannot override and convert at the same time")
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
				fmt.Println(v)
			}
		case *overrideTo != "":
			if _, err := bmnumbers.EventuallyCreateType(*overrideTo, nil); err != nil {
				log.Fatal(err)
			}
			if v := bmnumbers.GetType(*overrideTo); v == nil {
				log.Fatal("Error: Unknown type")
			} else {
				newType = v
				fmt.Println(v)
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

					if *overrideTo != "" {
						if err := bmnumbers.OverrideType(output, newType); err != nil {
							fmt.Println("Error: ", err)
						}
					}

					switch *dumpAs {
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
						fmt.Println("Error: Unknown dump format")
					}
				}
			}
		}
	}
}
