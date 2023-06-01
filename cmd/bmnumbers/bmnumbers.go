package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

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

var getPrefix = flag.String("get-prefix", "", "Get prefix from type")
var OmitPrefix = flag.Bool("omit-prefix", false, "Omit prefix")

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

	conf := new(bmnumbers.BMNumberConfig)
	conf.OmitPrefix = *OmitPrefix

	if *getPrefix != "" {
		if _, err := bmnumbers.EventuallyCreateType(*getPrefix, nil); err != nil {
			log.Fatal(err)
		}
		if v := bmnumbers.GetType(*getPrefix); v == nil {
			log.Fatal("Error: Unknown type")
		} else {
			fmt.Println(v.ShowPrefix())
		}
	} else if *serve {
		bmnumbers.Serve(conf)
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
				// open the file
				f, err := os.Open(argTo)
				if err != nil {
					log.Fatal(err)
				}
				of, err := os.Create(argTo + ".out")
				if err != nil {
					log.Fatal(err)
				}
				r := csv.NewReader(f)
				w := csv.NewWriter(of)
				for {
					record, err := r.Read()
					if err == io.EOF {
						break
					}
					if err != nil {
						log.Fatal(err)
					}

					recordC := make([]string, len(record))
					for i, v := range record {
						if output, err := bmnumbers.ImportString(v); err != nil {
							fmt.Println("Error: ", err)
							recordC[i] = v
							continue
						} else {

							if *convertTo != "" {
								if err := newType.Convert(output); err != nil {
									fmt.Println("Error: ", err)
									recordC[i] = v
									continue
								}
							}

							if *castTo != "" {
								if err := bmnumbers.CastType(output, newType); err != nil {
									fmt.Println("Error: ", err)
									recordC[i] = v
									continue
								}
							}

							switch *showAs {
							case "native":
								if value, err := output.ExportString(conf); err != nil {
									fmt.Println("Error: ", err)
									log.Fatal(err)
								} else {
									recordC[i] = value
								}
							case "bin":
								if value, err := output.ExportBinary(*withSize); err != nil {
									fmt.Println("Error: ", err)
									log.Fatal(err)
								} else {
									recordC[i] = value
								}
							case "unsigned":
								if value, err := output.ExportUint64(); err != nil {
									fmt.Println("Error: ", err)
									log.Fatal(err)
								} else {
									recordC[i] = fmt.Sprintf("%d", value)
								}
							default:
								log.Fatal("Error: Unknown visualization format (native, bin, unsigned)")
							}

						}
					}

					if err := w.Write(recordC); err != nil {
						log.Fatalln("error writing record to csv:", err)
					}

					w.Flush()

					if err := w.Error(); err != nil {
						log.Fatal(err)
					}
				}
				f.Close()
				of.Close()
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
						if value, err := output.ExportString(conf); err != nil {
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
