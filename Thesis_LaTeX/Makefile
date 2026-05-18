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

.PHONY: help all oneside twoside clean count

help:
	@echo "Usage: make <COMMAND>"
	@echo ""
	@echo "Commands:"
	@echo "  make all      Build oneside and twoside thesis PDF"
	@echo "  make oneside  Build oneside thesis PDF"
	@echo "  make twoside  Build twoside thesis PDF"
	@echo "  make clean    Remove generated LaTeX intermediate files"
	@echo "  make count    Count thesis characters"

all: oneside twoside

oneside:
	@latexmk $(LATEXMK_OPT) oneside_$(THESIS).tex

twoside:
	@latexmk $(LATEXMK_OPT) twoside_$(THESIS).tex

clean:
	-@latexmk -c -bibtex -silent oneside_$(THESIS).tex 2> /dev/null
	-@latexmk -c -bibtex -silent twoside_$(THESIS).tex 2> /dev/null

count:
	@texcount oneside_$(THESIS).tex -inc -chinese | awk '/total/ {getline; print "Total characters:",$$4}'
