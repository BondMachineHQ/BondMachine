.PHONY: installall
installall:
	@for i in `ls | grep -v Makefile` ; do cd ./$$i ; pwd ; if [[ $$i == "brvgasdl" || $$i == "bmgraph" || $$i == "bmdssim" ]]; then go install ; else CGO_ENABLED=0 go install ; fi ; cd .. ; done

