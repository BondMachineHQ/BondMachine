package basm

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func fragmentPruner(bi *BasmInstance) error {
	//TODO finish this

	if bi.debug {
		fmt.Println(green("\tProcessing fragments instances to be pruned:"))
	}

	pruneList := make([]int, 0)

	// Loop over the sections
	for p, fi := range bi.fis {
		toPrune := false
		if fi.GetMeta("pruned") == "true" {
			toPrune = true
		}
		if bi.debug {
			if toPrune {
				fmt.Println(red("\t\tFragment instance: "), fi)
			} else {
				fmt.Println(green("\t\tFragment instance: "), fi)
			}
		}

		if toPrune {
			// Get the input resources of the fragment
			fragment := fi.GetMeta("fragment")
			if bi.debug {
				fmt.Println("\t\t\tFragment: ", fragment)
			}

			resInS := bi.fragments[fragment].fragmentBody.GetMeta("resin")

			if bi.debug {
				fmt.Println("\t\t\tResources in: ", resInS)
			}

			resIn := make([]string, 0)
			for _, res := range strings.Split(resInS, ":") {
				if res != "" {
					resIn = append(resIn, res)
				}
			}

			reOuts := bi.fragments[fragment].fragmentBody.GetMeta("resout")

			if bi.debug {
				fmt.Println("\t\t\tResources out: ", reOuts)
			}

			resOut := make([]string, 0)
			for _, res := range strings.Split(reOuts, ":") {
				if res != "" {
					resOut = append(resOut, res)
				}
			}

			if len(resIn) != len(resOut) {
				return errors.New("The Fragment Instance " + fi.GetValue() + " has a different number of resources in and out and cannot be pruned")
			}

			// Inputs links
			inLinks := make([][]string, 0)
			for i, _ := range resIn {
				links, err := bi.GetLinks(FILINK, fi.GetValue(), strconv.Itoa(i), "input")
				if err != nil {
					return err
				}
				inLinks = append(inLinks, links)

				fmt.Println(bi.GetEndpoints(FILINK, links[0]))
			}

			if bi.debug {
				fmt.Println("\t\t\tInput links: ", inLinks)
			}

			// Outputs links
			outLinks := make([][]string, 0)
			for i, _ := range resOut {
				links, err := bi.GetLinks(FILINK, fi.GetValue(), strconv.Itoa(i), "output")
				if err != nil {
					return err
				}
				outLinks = append(outLinks, links)

				fmt.Println(bi.GetEndpoints(FILINK, links[0]))
			}

			if bi.debug {
				fmt.Println("\t\t\tOutput links: ", outLinks)
			}

			for i, inLink := range inLinks {
				// Get the endpoints of the link
				if inToRemoveEnd, outToKeepEnd, err := bi.GetEndpoints(FILINK, inLink[0]); err != nil {
					return err
				} else {
					for _, outLink := range outLinks[i] {
						// Get the endpoints of the link
						if _, outToChangeEnd, err := bi.GetEndpoints(FILINK, outLink); err != nil {
							return err
						} else {
							fiAttachToKeep := bi.fiLinkAttach[outToKeepEnd.seq]
							bi.fiLinkAttach[outToChangeEnd.seq].SetMeta("type", fiAttachToKeep.GetMeta("type"))
							bi.fiLinkAttach[outToChangeEnd.seq].SetMeta("index", fiAttachToKeep.GetMeta("index"))
							bi.fiLinkAttach[outToChangeEnd.seq].SetMeta("fi", fiAttachToKeep.GetMeta("fi"))
						}
					}

					if inToRemoveEnd.seq > outToKeepEnd.seq {
						toRemove := inToRemoveEnd.seq
						bi.fiLinkAttach = append(bi.fiLinkAttach[:toRemove], bi.fiLinkAttach[toRemove+1:]...)
						toRemove = outToKeepEnd.seq
						bi.fiLinkAttach = append(bi.fiLinkAttach[:toRemove], bi.fiLinkAttach[toRemove+1:]...)
					} else {
						toRemove := outToKeepEnd.seq
						bi.fiLinkAttach = append(bi.fiLinkAttach[:toRemove], bi.fiLinkAttach[toRemove+1:]...)
						toRemove = inToRemoveEnd.seq
						bi.fiLinkAttach = append(bi.fiLinkAttach[:toRemove], bi.fiLinkAttach[toRemove+1:]...)
					}

					pruneList = append(pruneList, p)

				}
			}
		}
	}

	for i := len(pruneList) - 1; i >= 0; i-- {
		bi.fis = append(bi.fis[:pruneList[i]], bi.fis[pruneList[i]+1:]...)
	}

	// panic("test")
	return nil
}
