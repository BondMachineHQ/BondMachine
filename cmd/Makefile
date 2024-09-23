SHELL := /bin/bash
.PHONY: all
all:
	@for i in `ls | grep -v Makefile` ; do cd ./$$i ; pwd ; if [[ $$i == "brvgasdl" || $$i == "bmgraph" || $$i == "bmdssim" ]]; then go install ; else CGO_ENABLED=0 go install ; fi ; cd .. ; done

