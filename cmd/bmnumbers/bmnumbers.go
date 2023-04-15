package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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

// Custom types
var linearDataRange = flag.String("linear-data-range", "", "Load a linear data range file (with the syntax index,filename)")

func init() {
	flag.Parse()
	// Load custom types
	if *linearDataRange != "" {

		// Get the linear quantizer ranges struct
		var lqRanges *map[int]bmnumbers.LinearDataRange
		for _, t := range bmnumbers.AllDynamicalTypes {
			if t.GetName() == "dyn_linear_quantizer" {
				lqRanges = t.(bmnumbers.DynLinearQuantizer).Ranges
			}
		}

		splitted := strings.Split(*linearDataRange, ",")
		if len(splitted)%2 != 0 {
			log.Fatal("Error: Invalid linear data range files")
		}

		// Load a file for each index
		for i := 0; i < len(splitted); i += 2 {
			index, err := strconv.Atoi(splitted[i])
			if err != nil {
				log.Fatal(err)
			}

			if index == 0 {
				log.Fatal("Error: Index cannot be 0 (reserved)")
			}

			// Check if the index is already present
			if _, ok := (*lqRanges)[index]; ok {
				log.Fatal("Error: Index already present")
			}

			filename := splitted[i+1]

			// Read all the lines of the file
			f, err := os.Open(filename)
			if err != nil {
				log.Fatal(err)
			}

			// Read all the lines of the file
			var min, max float64
			scanner := bufio.NewScanner(f)
			first := true
			for scanner.Scan() {
				line := scanner.Text()

				// Parse the min and max values
				if val, err := strconv.ParseFloat(line, 64); err == nil {
					if first {
						min = val
						max = val
						first = false
					}

					if val < min {
						min = val
					}
					if val > max {
						max = val
					}
				}
			}

			// Add the range to the map
			(*lqRanges)[index] = bmnumbers.LinearDataRange{Min: min, Max: max}
			f.Close()
		}
	}
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
				// fmt.Println(v)
			}
		case *overrideTo != "":
			if _, err := bmnumbers.EventuallyCreateType(*overrideTo, nil); err != nil {
				log.Fatal(err)
			}
			if v := bmnumbers.GetType(*overrideTo); v == nil {
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
