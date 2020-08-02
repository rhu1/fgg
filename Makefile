##
# Basic assumptions:
# - ANTLR v4 runtime library for Go is on the $GOPATH
# - The output dir of `go install` is on the $PATH
##


.PHONY: help
help:
	@echo "Example targets:"
	@echo "        make install         Install fgg to bin dir of Go workspace"
	@echo "                             (Assumes parsers generated)"
	@echo "        make test-all        Run all tests (assumes install)"
	@echo "        make clean-test-all  Clean up temp test files"
	@echo ""
	@echo "Look inside the Makefile for more specific test targets."
	@echo "To bypass the ANTLR parser generation, try:"
	@echo "        make install-pregen-parser"
	@echo ""


# cf. `go run` is like `build`, doesn't use caching (unlike `install`)
.PHONY: install
install:
	go install github.com/rhu1/fgg
	go install github.com/rhu1/fgg/cmd/fg2fgg

#.PHONY: check-install
#check-install:
#	which fgg
#TODO: check ANTLR

# Needs an appropriate antlr4 command, e.g., java -jar [antlr-4.7.1-complete.jar]
.PHONY: install-antlr
install-antlr:
	antlr4 -Dlanguage=Go -o parser/fg parser/FG.g4
	antlr4 -Dlanguage=Go -o parser/fgg parser/FGG.g4

.PHONY: install-pregen-parser
install-pregen-parser:
	cp -r parser/pregen/fg parser
	cp -r parser/pregen/fgg parser
	go install github.com/rhu1/fgg
	go install github.com/rhu1/fgg/cmd/fg2fgg

.PHONY: clean-install
clean-install: 
	rm -f ../../../../bin/fgg
	rm -f ../../../../bin/fg2fgg
	rm -f parser/fg/*
	rm -f parser/fgg/*


.PHONY: clean
clean: clean-install clean-test-all


.PHONY: test-all
test-all: test test-against-go

.PHONY: clean-test-all
clean-test-all: clean-test-fg2fgg clean-test-monom-against-go clean-test-oblit

.PHONY: test
test: test-fg test-fgg test-fg2fgg

.PHONY: test-against-go
test-against-go: test-fg-examples-against-go test-monom-against-go


.PHONY: test-fg
test-fg: test-fg-unit test-fg-examples


.PHONY: test-fg-unit
test-fg-unit:
	go test github.com/rhu1/fgg/fg


.PHONY: test-fg-examples
test-fg-examples:
	$(call eval_fg,fg/examples/hello/hello.go,10)
	$(call eval_fg,fg/examples/hello/fmtprintf.go,10)

	$(call eval_fg,fg/examples/misc/booleans/booleans.go,-1)
	$(call eval_fg,fg/examples/misc/compose/compose.go,-1)
	$(call eval_fg,fg/examples/misc/equal/equal.go,-1)
	$(call eval_fg,fg/examples/misc/incr/incr.go,-1)
	$(call eval_fg,fg/examples/misc/map/map.go,-1)
	$(call eval_fg,fg/examples/misc/not/not.go,-1)

	$(call eval_fg,fg/examples/oopsla20/fig1/functions.go,-1)
	$(call eval_fg,fg/examples/oopsla20/fig2/lists.go,-1)
	$(call eval_fg,fg/examples/oopsla20/fig7/monom.go,-1)


# cf. [cmd] > output.txt
#     diff output.txt correct.txt
.PHONY: test-fg-examples-against-go
test-fg-examples-against-go:
	@$(call test_fg_against_go,fg/examples/misc/booleans/booleans.go)
	@$(call test_fg_against_go,fg/examples/misc/compose/compose.go)
	@$(call test_fg_against_go,fg/examples/misc/equal/equal.go)
	@$(call test_fg_against_go,fg/examples/misc/incr/incr.go)
	@$(call test_fg_against_go,fg/examples/misc/map/map.go)
	@$(call test_fg_against_go,fg/examples/misc/not/not.go)

	@$(call test_fg_against_go,fg/examples/oopsla20/fig1/functions.go)
	@$(call test_fg_against_go,fg/examples/oopsla20/fig2/lists.go)
	@$(call test_fg_against_go,fg/examples/oopsla20/fig7/monom.go)


.PHONY: test-fgg
#test-fgg: test-fgg-unit test-fgg-examples
test-fgg: test-fgg-unit test-nomono-bad simulate-monom simulate-oblit
# add monom, oblit?


.PHONY: test-fgg-unit
test-fgg-unit:
	go test github.com/rhu1/fgg/fgg


# Subsumed by, e.g., simulate-monom
.PHONY: test-fgg-examples
test-fgg-examples:
	$(call eval_fgg,fgg/examples/hello/hello.fgg,10)
	$(call eval_fgg,fgg/examples/hello/fmtprintf.fgg,10)

	$(call eval_fgg,fgg/examples/misc/booleans/booleans.fgg,-1)
	$(call eval_fgg,fgg/examples/misc/compose/compose.fgg,-1)
	$(call eval_fgg,fgg/examples/misc/graph/graph.fgg,-1)
	$(call eval_fgg,fgg/examples/misc/irregular/irregular.fgg,-1)
	$(call eval_fgg,fgg/examples/misc/map/map.fgg,-1)
	$(call eval_fgg,fgg/examples/misc/monomorph/monomorph.fgg,-1)

	$(call eval_fgg,fgg/examples/monom/box/box.fgg,10)
	$(call eval_fgg,fgg/examples/monom/box/box2.fgg,10)

	$(call eval_fgg,fgg/examples/monom/misc/ifacebox.fgg,-1)

	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/iface-embedding-simple.go,-1)
	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/iface-embedding.go,-1)
	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/rcver-iface.go,-1)
	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/one-pass-prob.go,-1)
	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/contamination.go,-1)
	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/struct-poly-rec.go,-1)
	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/Parameterised-Map.go,-1)
	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/alternate.go,10)
	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/i-closure.go,-1)
	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/i-closure-bad.go,-1)
	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/meth-clash.go,7)
	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/param-meth-cast.go,2)
	$(call eval_fgg,fgg/examples/monom/misc/mono-ok/poly-rec-iface.go,10)

	$(call eval_fgg,fgg/examples/oopsla20/fig3/functions.fgg,-1)
	$(call eval_fgg,fgg/examples/oopsla20/fig4/lists.fgg,-1)
	$(call eval_fgg,fgg/examples/oopsla20/fig5/graph.fgg,-1)
	$(call eval_fgg,fgg/examples/oopsla20/fig6/expression.fgg,-1)


.PHONY: test-nomono-bad
test-nomono-bad:
	@$(call nomono_bad,fgg/examples/monom/box/box.fgg)

	@$(call nomono_bad,fgg/examples/monom/misc/mono-ko/incompleteness-subtyping.go)
	@$(call nomono_bad,fgg/examples/monom/misc/mono-ko/monom-imp.go)
	@$(call nomono_bad,fgg/examples/monom/misc/mono-ko/mutual-poly-rec.go)
	@$(call nomono_bad,fgg/examples/monom/misc/mono-ko/mutual-rec-iface.go)
	@$(call nomono_bad,fgg/examples/monom/misc/mono-ko/nested-fix.go)
	@$(call nomono_bad,fgg/examples/monom/misc/mono-ko/two-type-param.go)


.PHONY: simulate-monom
simulate-monom:
	$(call sim_monom,fgg/examples/hello/hello.fgg,10)
	$(call sim_monom,fgg/examples/hello/fmtprintf.fgg,10)

	$(call sim_monom,fgg/examples/misc/booleans/booleans.fgg,-1)
	$(call sim_monom,fgg/examples/misc/compose/compose.fgg,-1)
	$(call sim_monom,fgg/examples/misc/graph/graph.fgg,-1)
	$(call sim_monom,fgg/examples/misc/irregular/irregular.fgg,-1)
	$(call sim_monom,fgg/examples/misc/map/map.fgg,-1)
	$(call sim_monom,fgg/examples/misc/monomorph/monomorph.fgg,-1)

	$(call sim_monom,fgg/examples/monom/box/box.fgg,10)
	$(call sim_monom,fgg/examples/monom/box/box2.fgg,10)

	$(call sim_monom,fgg/examples/monom/misc/ifacebox.fgg,-1)

	$(call sim_monom,fgg/examples/monom/misc/mono-ok/iface-embedding-simple.go,-1)
	$(call sim_monom,fgg/examples/monom/misc/mono-ok/iface-embedding.go,-1)
	$(call sim_monom,fgg/examples/monom/misc/mono-ok/rcver-iface.go,-1)
	$(call sim_monom,fgg/examples/monom/misc/mono-ok/one-pass-prob.go,-1)
	$(call sim_monom,fgg/examples/monom/misc/mono-ok/contamination.go,-1)

# TODO: add to oblit
	$(call sim_monom,fgg/examples/monom/misc/mono-ok/struct-poly-rec.go,-1)
	$(call sim_monom,fgg/examples/monom/misc/mono-ok/Parameterised-Map.go,-1)
	$(call sim_monom,fgg/examples/monom/misc/mono-ok/alternate.go,10)
	$(call sim_monom,fgg/examples/monom/misc/mono-ok/i-closure.go,-1)
	$(call sim_monom,fgg/examples/monom/misc/mono-ok/i-closure-bad.go,-1)
	$(call sim_monom,fgg/examples/monom/misc/mono-ok/meth-clash.go,7)
	$(call sim_monom,fgg/examples/monom/misc/mono-ok/param-meth-cast.go,2)
	$(call sim_monom,fgg/examples/monom/misc/mono-ok/poly-rec-iface.go,10)

	$(call sim_monom,fgg/examples/oopsla20/fig3/functions.fgg,-1)
	$(call sim_monom,fgg/examples/oopsla20/fig4/lists.fgg,-1)
	$(call sim_monom,fgg/examples/oopsla20/fig5/graph.fgg,-1)
	$(call sim_monom,fgg/examples/oopsla20/fig6/expression.fgg,-1)


# Non-terminating examples tested by simulate-monom
.PHONY: test-monom-against-go
test-monom-against-go:
	@$(call eval_monom_fgg_against_go,fgg/examples/misc/booleans/booleans.fgg,tmp/test/fg/booleans,booleans.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/misc/compose/compose.fgg,tmp/test/fg/compose,compose.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/misc/graph/graph.fgg,tmp/test/fg/graph,graph.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/misc/irregular/irregular.fgg,tmp/test/fg/irregular,irregular.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/misc/map/map.fgg,tmp/test/fg/map,map.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/misc/monomorph/monomorph.fgg,tmp/test/fg/monomorph,monomorph.go)

#@$(call eval_monom_fgg,fgg/examples/monom/box/box2.fgg,10,tmp/test/fg/monom/box,box2.go)

	@$(call eval_monom_fgg_against_go,fgg/examples/monom/misc/ifacebox.fgg,tmp/test/fg/monom/misc/ifacebox,ifacebox.go)

	@$(call eval_monom_fgg_against_go,fgg/examples/monom/misc/mono-ok/iface-embedding-simple.go,tmp/test/fg/monom/misc/mono-ok/iface-embedding-simple,iface-embedding-simple.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/monom/misc/mono-ok/iface-embedding.go,tmp/test/fg/monom/misc/mono-ok/iface-embedding,iface-embedding.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/monom/misc/mono-ok/rcver-iface.go,tmp/test/fg/monom/misc/mono-ok/rcver-iface,rcver-iface.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/monom/misc/mono-ok/one-pass-prob.go,tmp/test/fg/monom/misc/mono-ok/one-pass-prob,one-pass-prob.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/monom/misc/mono-ok/contamination.go,tmp/test/fg/monom/misc/mono-ok/contamination,contamination.go)

# TODO: add to oblit
	@$(call eval_monom_fgg_against_go,fgg/examples/monom/misc/mono-ok/struct-poly-rec.go,tmp/test/fg/monom/misc/mono-ok/struct-poly-rec,struct-poly-rec.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/monom/misc/mono-ok/Parameterised-Map.go,tmp/test/fg/monom/misc/mono-ok/Parameterised-Map,Parameterised-Map.go)
#@$(call eval_monom_fgg,fgg/examples/monom/misc/mono-ok/alternate.go,10,tmp/test/fg/monom/misc/mono-ok/alternate,alternate.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/monom/misc/mono-ok/i-closure.go,tmp/test/fg/monom/misc/mono-ok/i-closure,i-closure.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/monom/misc/mono-ok/i-closure-bad.go,tmp/test/fg/monom/misc/mono-ok/i-closure-bad,i-closure-bad.go)
#@$(call eval_monom_fgg,fgg/examples/monom/misc/mono-ok/meth-clash.go,7,tmp/test/fg/monom/misc/mono-ok/meth-clash,meth-clash.go)
#@$(call eval_monom_fgg,fgg/examples/monom/misc/mono-ok/param-meth-cast.go,2,tmp/test/fg/monom/misc/mono-ok/param-meth-cast,param-meth-cast.go)
#@$(call eval_monom_fgg,fgg/examples/monom/misc/mono-ok/poly-rec-iface.go,10,tmp/test/fg/monom/misc/mono-ok/poly-rec-iface,poly-rec-iface.go)

	#mkdir -p tmp/test/fg/monom/misc/mono-ko

	@$(call eval_monom_fgg_against_go,fgg/examples/oopsla20/fig3/functions.fgg,tmp/test/fg/oopsla20/functions,functions.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/oopsla20/fig4/lists.fgg,tmp/test/fg/oopsla20/lists,lists.go)
	@$(call eval_monom_fgg_against_go,fgg/examples/oopsla20/fig5/graph.fgg,tmp/test/fg/oopsla20/graph,graph.go)
#@$(call eval_monom_fgg_against_go,fgg/examples/oopsla20/fig6/expression.fgg,tmp/test/fg/oopsla20/expression,expression.go)  #go= "({{{}}}+{{}})" 

.PHONY: clean-test-monom-against-go
clean-test-monom-against-go:
	$(call rm_monom,tmp/test/fg/booleans,booleans.go)
	$(call rm_monom,tmp/test/fg/compose,compose.go)
	$(call rm_monom,tmp/test/fg/graph,graph.go)
	$(call rm_monom,tmp/test/fg/irregular,irregular.go)
	$(call rm_monom,tmp/test/fg/map,map.go)
	$(call rm_monom,tmp/test/fg/monomorph,monomorph.go)

#rm -f tmp/test/fg/monom/box/box2.go
#rm -fd tmp/test/fg/monom/box

	$(call rm_monom,tmp/test/fg/monom/misc/ifacebox,ifacebox.go)

	$(call rm_monom,tmp/test/fg/monom/misc/mono-ok/iface-embedding-simple,iface-embedding-simple.go)
	$(call rm_monom,tmp/test/fg/monom/misc/mono-ok/iface-embedding,iface-embedding.go)
	$(call rm_monom,tmp/test/fg/monom/misc/mono-ok/rcver-iface,rcver-iface.go)
	$(call rm_monom,tmp/test/fg/monom/misc/mono-ok/one-pass-prob,one-pass-prob.go)
	$(call rm_monom,tmp/test/fg/monom/misc/mono-ok/contamination,contamination.go)
	$(call rm_monom,tmp/test/fg/monom/misc/mono-ok/struct-poly-rec,struct-poly-rec.go)
	$(call rm_monom,tmp/test/fg/monom/misc/mono-ok/Parameterised-Map,Parameterised-Map.go)
	$(call rm_monom,tmp/test/fg/monom/misc/mono-ok/i-closure,i-closure.go)
	$(call rm_monom,tmp/test/fg/monom/misc/mono-ok/i-closure-bad,i-closure-bad.go)

	rm -fd tmp/test/fg/monom/misc/mono-ok
	rm -fd tmp/test/fg/monom/misc/mono-ko
	rm -fd tmp/test/fg/monom/misc
	rm -fd tmp/test/fg/monom

	$(call rm_monom,tmp/test/fg/oopsla20/functions,functions.go)
	$(call rm_monom,tmp/test/fg/oopsla20/lists,lists.go)
	$(call rm_monom,tmp/test/fg/oopsla20/graph,graph.go)
	$(call rm_monom,tmp/test/fg/oopsla20/expression,expression.go)

	rm -fd tmp/test/fg/oopsla20


##
# Aux
##

define eval_fg
	fgg -eval=$(2) $(1)
endef
	#RES=`fgg -eval=$(2) $(1)`; \
	#EXIT=$$?; if [ $$EXIT -ne 0 ]; then exit $$EXIT; fi; \
	#echo $$RES


# N.B. double-dollar
define test_fg_against_go
	echo "Testing "$(1)" against Go:"; \
	EXP=`go run $(1)`; \
	EXIT=$$?; if [ $$EXIT -ne 0 ]; then exit $$EXIT; fi; \
	echo "go="$$EXP; \
	ACT=`fgg -eval=-1 -printf $(1)`; \
	EXIT=$$?; if [ $$EXIT -ne 0 ]; then exit $$EXIT; fi; \
	echo "fg="$$ACT; \
	if [ "$$EXP" != "$$ACT" ]; then \
		echo "Not equal."; \
		exit 1; \
	fi
endef


define eval_fgg
	fgg -fgg -eval=$(2) $(1)
endef
	#RES=`fgg -fgg -eval=$(2) $(1)`; \
	#EXIT=$$?; if [ $$EXIT -ne 0 ]; then exit $$EXIT; fi; \
	#echo $$RES


# TODO: make error check more specific
define nomono_bad
	echo "Testing bad nomono in "$(1)":"
	RES=`fgg -fgg -monomc=-- $(1) 2> /dev/null`; \
	EXIT=$$?; if [ $$EXIT -eq 0 ]; then \
		echo "Expected nomono violation, but none occurred."; \
		exit 1; \
	fi; 
endef


define sim_monom
	fgg -test-monom -eval=$(2) $(1)
endef
	#`fgg -test-monom -eval=$(2) $(1)`; \
	#EXIT=$$?; if [ $$EXIT -ne 0 ]; then exit $$EXIT; fi


define eval_monom_fgg
	mkdir -p $(3); \
	RES=`fgg -fgg -eval=$(2) -monomc=$(3)/$(4) $(1)`; \
	EXIT=$$?; if [ $$EXIT -ne 0 ]; then exit $$EXIT; fi; \
	echo "fgg="$$RES; \
	EXP=`fgg -eval=$(2) $(3)/$(4)`; \
	EXIT=$$?; if [ $$EXIT -ne 0 ]; then exit $$EXIT; fi; \
	echo "fg ="$$EXP
endef


define eval_monom_fgg_against_go
	echo "Testing monom of "$(1)" against Go:"; \
	mkdir -p $(2); \
	RES=`fgg -fgg -eval=-1 -monomc=$(2)/$(3) $(1)`; \
	EXIT=$$?; if [ $$EXIT -ne 0 ]; then exit $$EXIT; fi; \
	echo "fgg="$$RES; \
	EXP=`go run $(2)/$(3)`; \
	echo "go ="$$EXP; \
	ACT=`fgg -eval=-1 -printf $(2)/$(3)`; \
	echo "fg ="$$ACT; \
	if [ "$$EXP" != "$$ACT" ]; then \
		echo "Not equal."; \
		exit 1; \
	fi
endef


define rm_monom
	rm -f $(1)/$(2); \
	rm -fd $(1)
endef


#.PHONY: foo
#foo:
#	declare -a arr=(\
#		"element1" \
#		"element2" \
#		"element3"); \
#	for i in "$${arr[@]}"; \
#	do \
#   	 echo "$$i"; \
#	done



##
# TODO: all the below need updating
##

.PHONY: simulate-oblit
simulate-oblit:
	fgg -test-oblit -eval=-1 fgg/examples/misc/booleans/booleans.fgg
	fgg -test-oblit -eval=-1 fgg/examples/misc/compose/compose.fgg
	fgg -test-oblit -eval=-1 fgg/examples/misc/graph/graph.fgg
	fgg -test-oblit -eval=-1 fgg/examples/misc/irregular/irregular.fgg
	fgg -test-oblit -eval=-1 fgg/examples/misc/map/map.fgg
	fgg -test-oblit -eval=-1 fgg/examples/misc/monomorph/monomorph.fgg
# TODO: currently trying to run to termination
#fgg -test-oblit -eval=10 fgg/examples/monom/box/box.fgg
#fgg -test-oblit -eval=10 fgg/examples/monom/box/box2.fgg

	fgg -test-oblit -eval=-1 fgg/examples/monom/misc/ifacebox.fgg

# TODO?
#fgg -test-oblit -eval=-1 fgg/examples/monom/misc/iface-embedding-simple.go
#fgg -test-oblit -eval=-1 fgg/examples/monom/misc/iface-embedding.go

	fgg -test-oblit -eval=-1 fgg/examples/monom/misc/mono-ok/rcver-iface.go
	fgg -test-oblit -eval=-1 fgg/examples/monom/misc/mono-ok/one-pass-prob.go
	fgg -test-oblit -eval=-1 fgg/examples/monom/misc/mono-ok/contamination.go


.PHONY: test-oblit
test-oblit:
	mkdir -p tmp/test-oblit/fgr/booleans
	fgg -fgg -oblitc=tmp/test-oblit/fgr/booleans/booleans.fgr -oblit-eval=-1 fgg/examples/misc/booleans/booleans.fgg
# TODO: standalone FGR execution (.fgr output currently unused)
# 
	mkdir -p tmp/test-oblit/fgr/compose
	fgg -fgg -oblitc=tmp/test-oblit/fgr/compose/compose.fgr -oblit-eval=-1 fgg/examples/misc/compose/compose.fgg

	mkdir -p tmp/test-oblit/fgr/graph
	fgg -fgg -oblitc=tmp/test-oblit/fgr/graph/graph.fgr -oblit-eval=-1 fgg/examples/misc/graph/graph.fgg

	mkdir -p tmp/test-oblit/fgr/irregular
	fgg -fgg -oblitc=tmp/test-oblit/fgr/irregular/irregular.fgr -oblit-eval=-1 fgg/examples/misc/irregular/irregular.fgg

	mkdir -p tmp/test-oblit/fgr/map
	fgg -fgg -oblitc=tmp/test-oblit/fgr/map/map.fgr -oblit-eval=-1 fgg/examples/misc/map/map.fgg

	mkdir -p tmp/test-oblit/fgr/monomorph
	fgg -fgg -oblitc=tmp/test-oblit/fgr/monomorph/monomorph.fgr -oblit-eval=-1 fgg/examples/misc/monomorph/monomorph.fgg

	mkdir -p tmp/test-oblit/fgr/box
	fgg -fgg -oblitc=tmp/test-oblit/fgr/box/box.fgr -oblit-eval=10 fgg/examples/monom/box/box.fgg
	fgg -fgg -oblitc=tmp/test-oblit/fgr/box/box2.fgr -oblit-eval=10 fgg/examples/monom/box/box2.fgg

	mkdir -p tmp/test-oblit/fgr/misc
	fgg -fgg -oblitc=tmp/test-oblit/fgr/misc/ifacebox.fgr -oblit-eval=-1 fgg/examples/monom/misc/ifacebox.fgg
# TODO: i/face embedding?
#fgg -fgg -oblitc=tmp/test-oblit/fgr/misc/iface-embedding-simple.fgr -oblit-eval=-1 fgg/examples/monom/misc/iface-embedding-simple.go
#fgg -fgg -oblitc=tmp/test-oblit/fgr/misc/iface-embedding.fgr -oblit-eval=-1 fgg/examples/monom/misc/iface-embedding.go

	mkdir -p tmp/test-oblit/fgr/misc/mono-ok
	fgg -fgg -oblitc=tmp/test-oblit/fgr/misc/mono-ok/rcver-iface.fgr -oblit-eval=-1 fgg/examples/monom/misc/mono-ok/rcver-iface.go
	fgg -fgg -oblitc=tmp/test-oblit/fgr/misc/mono-ok/one-pass-prob.fgr -oblit-eval=-1 fgg/examples/monom/misc/mono-ok/one-pass-prob.go
	fgg -fgg -oblitc=tmp/test-oblit/fgr/misc/mono-ok/contamination.fgr -oblit-eval=-1 fgg/examples/monom/misc/mono-ok/contamination.go

	mkdir -p tmp/test-oblit/fgr/misc/mono-ko

.PHONY: clean-test-oblit
clean-test-oblit:
	rm -f tmp/test-oblit/fgr/booleans/booleans.fgr
	rm -fd tmp/test-oblit/fgr/booleans

	rm -f tmp/test-oblit/fgr/compose/compose.fgr
	rm -fd tmp/test-oblit/fgr/compose

	rm -f tmp/test-oblit/fgr/graph/graph.fgr
	rm -fd tmp/test-oblit/fgr/graph

	rm -f tmp/test-oblit/fgr/irregular/irregular.fgr
	rm -fd tmp/test-oblit/fgr/irregular

	rm -f tmp/test-oblit/fgr/map/map.fgr
	rm -fd tmp/test-oblit/fgr/map

	rm -f tmp/test-oblit/fgr/monomorph/monomorph.fgr
	rm -fd tmp/test-oblit/fgr/monomorph

	rm -f tmp/test-oblit/fgr/box/box.fgr
	rm -f tmp/test-oblit/fgr/box/box2.fgr
	rm -fd tmp/test-oblit/fgr/box

	rm -f tmp/test-oblit/fgr/misc/mono-ok/rcver-iface.fgr
	rm -f tmp/test-oblit/fgr/misc/mono-ok/one-pass-prob.fgr
	rm -f tmp/test-oblit/fgr/misc/mono-ok/contamination.fgr
	rm -fd tmp/test-oblit/fgr/misc/mono-ok

	rm -fd tmp/test-oblit/fgr/misc/mono-ko

	rm -f tmp/test-oblit/fgr/misc/ifacebox.fgr
	rm -fd tmp/test-oblit/fgr/misc


.PHONY: test-fg2fgg
test-fg2fgg:
	mkdir -p tmp/test/fgg/booleans
	fg2fgg fg/examples/misc/booleans/booleans.go > tmp/test/fgg/booleans/booleans.fgg
	fgg -fgg -eval=-1 tmp/test/fgg/booleans/booleans.fgg

	mkdir -p tmp/test/fgg/compose
	fg2fgg fg/examples/misc/compose/compose.go > tmp/test/fgg/compose/compose.fgg
	fgg -fgg -eval=-1 tmp/test/fgg/compose/compose.fgg

	mkdir -p tmp/test/fgg/equal
	fg2fgg fg/examples/misc/equal/equal.go > tmp/test/fgg/equal/equal.fgg
	fgg -fgg -eval=-1 tmp/test/fgg/equal/equal.fgg

	mkdir -p tmp/test/fgg/incr
	fg2fgg fg/examples/misc/incr/incr.go > tmp/test/fgg/incr/incr.fgg
	fgg -fgg -eval=-1 tmp/test/fgg/incr/incr.fgg

	mkdir -p tmp/test/fgg/map
	fg2fgg fg/examples/misc/map/map.go > tmp/test/fgg/map/map.fgg
	fgg -fgg -eval=-1 tmp/test/fgg/map/map.fgg

	mkdir -p tmp/test/fgg/not
	fg2fgg fg/examples/misc/not/not.go > tmp/test/fgg/not/not.fgg
	fgg -fgg -eval=-1 tmp/test/fgg/not/not.fgg

# TODO: run fg_test.go unit tests through fg2fgg

.PHONY: clean-test-fg2fgg
clean-test-fg2fgg:
	rm -f tmp/test/fgg/booleans/booleans.fgg
	rm -fd tmp/test/fgg/booleans

	rm -f tmp/test/fgg/compose/compose.fgg
	rm -fd tmp/test/fgg/compose

	rm -f tmp/test/fgg/equal/equal.fgg
	rm -fd tmp/test/fgg/equal

	rm -f tmp/test/fgg/incr/incr.fgg
	rm -fd tmp/test/fgg/incr

	rm -f tmp/test/fgg/map/map.fgg
	rm -fd tmp/test/fgg/map

	rm -f tmp/test/fgg/not/not.fgg
	rm -fd tmp/test/fgg/not

