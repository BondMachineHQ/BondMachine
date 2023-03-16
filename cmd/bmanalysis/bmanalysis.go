package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"github.com/BondMachineHQ/BondMachine/pkg/bmanalysis"
)

var projectLists = flag.String("projectLists", "", "Comma separeted lists of projects")
var hdlFile = flag.String("ipynb-file", "analysis.ipynb", "Name of the file to write the Python")


func init() {
	flag.Parse()
}

func main() {
	bmanalysis := bmanalysis.CreateAnalysisTemplate()

	if *projectLists != "" {
		for _, project := range strings.Split(*projectLists, ",") {
			bmanalysis.ProjectLists = append(bmanalysis.ProjectLists, project)
		}
	} else {
		log.Fatal("No project lists specified")
	}

	if *hdlFile != "" {
		hdl, err := bmanalysis.WritePython()
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(*hdlFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err = f.WriteString(hdl)
		if err != nil {
			log.Fatal(err)
		}
	}

}
