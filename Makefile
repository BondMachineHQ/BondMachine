SHELL := /bin/bash
.PHONY: all
all:
	make -C cmd installall

.PHONY: align
align:
	cp -r pkg/bmnumbers/*.md cmd/bmnumbers
	cp -r pkg/bmqsim/*.{md,png} cmd/bmqsim
	cp -r pkg/bmstack/*.{md,png} cmd/bmstack
	cp -r pkg/procbuilder/*.md cmd/procbuilder
