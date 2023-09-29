package basm

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"text/template"

	"github.com/BondMachineHQ/BondMachine/pkg/bmline"
)

func templateResolver(bi *BasmInstance) error {

	if err := bi.templateAutoMark(); err != nil {
		return err
	}

	// Computing which CP needs a templated version of the code
	sort.Sort(bmline.ByName(bi.cps))

	sectionRem := make(map[string]struct{})
	fragmentRem := make(map[string]struct{})

	for _, cp := range bi.cps {

		if cp.GetMeta("templated") == "true" {
			if bi.debug {
				fmt.Print("\t\t" + green("CP: ") + yellow(cp.GetValue()) + " is templated")
			}

			romCode := cp.GetMeta("romcode")
			if romCode == "" {
				return errors.New("CP rom code not found")
			}
			if bi.debug {
				fmt.Println(" - " + green("rom code: ") + yellow(romCode))
			}

			// Generating a new name uniq for the section and adding it to the list of sections
			i := 0
			guessedName := romCode + "_templ_" + fmt.Sprint(i)
			for {
				if _, ok := bi.sections[guessedName]; ok {
					i++
					guessedName = romCode + "_templ_" + fmt.Sprint(i)
					continue
				}
				break
			}
			newSection := "%" + "section " + guessedName + " .romtext " + bi.sections[romCode].sectionBody.ListMeta() + "\n"
			newSection += bi.sections[romCode].writeText()
			newSection += "%" + "endsection\n"

			// fmt.Printf(newSection)

			td := createBasicTemplateData()

			for key, value := range cp.LoopMeta() {
				td.Params[key] = value
			}
			//fmt.Println(td.Params)
			var f bytes.Buffer

			t, err := template.New("template").Funcs(td.funcMap).Parse(newSection)
			if err != nil {
				return err
			}

			err = t.Execute(&f, *td)
			if err != nil {
				return err
			}

			newSection = f.String()

			if err := bi.ParseAssemblyString(newSection, basmParser); err != nil {
				return err
			}

			if isTemplate(newSection) {
				bi.sections[guessedName].sectionBody.SetMeta("template", "true")
			} else {
				bi.sections[guessedName].sectionBody.SetMeta("template", "false")
				sectionRem[romCode] = struct{}{}
			}

			cp.SetMeta("romcode", guessedName)

			// fmt.Printf(newSection)
		} else {
			if bi.debug {
				fmt.Print("\t\t" + green("CP: ") + yellow(cp.GetValue()) + " is not templated")
			}
		}
		if bi.debug {
			fmt.Println()
		}
	}

	// Computing which Fragment Instance needs a templated version of the code
	for _, fi := range bi.fis {

		if fi.GetMeta("templated") == "true" {
			if bi.debug {
				fmt.Print("\t\t" + green("FI: ") + yellow(fi.GetValue()) + " is templated")
			}

			fragCode := fi.GetMeta("fragment")
			if fragCode == "" {
				return errors.New("fragment rom code not found")
			}
			if bi.debug {
				fmt.Println(" - " + green("fragment code: ") + yellow(fragCode))
			}

			// Generating a new name uniq for the section and adding it to the list of sections
			i := 0
			guessedName := fragCode + "_templ_" + fmt.Sprint(i)
			for {
				if _, ok := bi.fragments[guessedName]; ok {
					i++
					guessedName = fragCode + "_templ_" + fmt.Sprint(i)
					continue
				}
				break
			}
			newFragment := "%" + "fragment " + guessedName + " " + bi.fragments[fragCode].fragmentBody.ListMeta() + "\n"
			newFragment += bi.fragments[fragCode].writeText()
			newFragment += "%" + "endfragment\n"

			// fmt.Printf(newFragment)

			td := createBasicTemplateData()

			for key, value := range fi.LoopMeta() {
				td.Params[key] = value
			}
			//fmt.Println(td.Params)
			var f bytes.Buffer

			t, err := template.New("template").Funcs(td.funcMap).Parse(newFragment)
			if err != nil {
				return err
			}

			err = t.Execute(&f, *td)
			if err != nil {
				return err
			}

			newFragment = f.String()

			if err := bi.ParseAssemblyString(newFragment, basmParser); err != nil {
				return err
			}

			if isTemplate(newFragment) {
				bi.fragments[guessedName].fragmentBody.SetMeta("template", "true")
			} else {
				bi.fragments[guessedName].fragmentBody.SetMeta("template", "false")
				fragmentRem[fragCode] = struct{}{}
			}

			fi.SetMeta("fragment", guessedName)

			// fmt.Printf(newSection)
		} else {
			if bi.debug {
				fmt.Print("\t\t" + green("FI: ") + yellow(fi.GetValue()) + " is not templated")
			}
		}
		if bi.debug {
			fmt.Println()
		}
	}

	// Remove all the templated sections
	for sectionName := range sectionRem {
		delete(bi.sections, sectionName)
	}

	// Remove all the templated fragments
	for fragmentName := range fragmentRem {
		delete(bi.fragments, fragmentName)
	}

	return nil
}
