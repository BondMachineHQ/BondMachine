package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"github.com/BondMachineHQ/BondMachine/pkg/bmanalysis"
)

var projectsList = flag.String("projects-list", "", "Comma separeted lists of projects")
var pythonFile = flag.String("ipynb-file", "analysis.ipynb", "Name of the file to write the Python")
var pivotRun = flag.Int("pivot-run", 0, "Index of run to use as pivot to compare with other results")

func init() {
	flag.Parse()
}

func main() {
	bmanalysis := bmanalysis.CreateAnalysisTemplate()

	if *projectsList != "" {
		for _, project := range strings.Split(*projectsList, ",") {
			projectT := strings.TrimSpace(project)
			if projectT != "" {
				bmanalysis.ProjectsList = append(bmanalysis.ProjectsList, projectT)
			}
		}
	} else {
		log.Fatal("No project lists specified")
	}

	if *pivotRun < 0 && *pivotRun > len(bmanalysis.ProjectsList) {
		log.Fatal("Invalid data width")
	}
	bmanalysis.PivotRun = *pivotRun

	if *pythonFile != "" {
		python, err := bmanalysis.WritePython()
		if err != nil {
			log.Fatal(err)
		}

		f, err := os.Create(*pythonFile)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err = f.WriteString(python)
		if err != nil {
			log.Fatal(err)
		}
	}

}
