THESIS := main
ENGINE ?= -xelatex
LATEXMK_OPT := \
	-quiet \
	-file-line-error \
	-halt-on-error \
	-interaction=nonstopmode \
	-shell-escape \
	-synctex=1 \
	-recorder \
	-usepretex="\listfiles" \
	$(ENGINE)

.DEFAULT_GOAL := help

.PHONY: help all clean count

help:
	@echo "Usage: make <COMMAND>"
	@echo ""
	@echo "Commands:"
	@echo "  make all    Build thesis PDF"
	@echo "  make clean  Remove generated LaTeX intermediate files"
	@echo "  make count  Count thesis characters"

all:
	@latexmk $(LATEXMK_OPT) $(THESIS).tex

clean:
	-@latexmk -c -bibtex -silent $(THESIS).tex 2> /dev/null

count:
	@texcount $(THESIS).tex -inc -chinese | awk '/total/ {getline; print "Total characters:",$$4}'
