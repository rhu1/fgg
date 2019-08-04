# Temporary makefile for testing, e.g., -monom
# TODO: ^should be integrated with `go test` instead
# TODO: but if this makefile is to be retained, needs refactoring


#.PHONY: build
#build:
#	go build github.com/rhu1/fgg
#
#
#.PHONY: install
#install:
#	go install github.com/rhu1/fgg


.PHONY: test
test: test-fg test-fgg test-monom test-fg2fgg


.PHONY: clean
clean: clean-test

#rm -f ../../../../bin/fgg.exe
#make test-clean


.PHONY: test-fg
test-fg: test-fg-unit test-fg-examples


.PHONY: test-fg-unit
test-fg-unit:
	go test github.com/rhu1/fgg/fg


.PHONY: test-fg-examples
test-fg-examples:
	go run github.com/rhu1/fgg -eval=-1 fg/examples/popl20/booleans/booleans.go
	go run github.com/rhu1/fgg -eval=-1 fg/examples/popl20/compose/compose.go
	go run github.com/rhu1/fgg -eval=-1 fg/examples/popl20/equal/equal.go
	go run github.com/rhu1/fgg -eval=-1 fg/examples/popl20/incr/incr.go
	go run github.com/rhu1/fgg -eval=-1 fg/examples/popl20/map/map.go
	go run github.com/rhu1/fgg -eval=-1 fg/examples/popl20/not/not.go

# TODO: currently examples testing limited to "good" examples


.PHONY: test-fgg
test-fgg: test-fgg-unit test-fgg-examples


.PHONY: test-fgg-unit
test-fgg-unit:
	go test github.com/rhu1/fgg/fgg


.PHONY: test-fgg-examples
test-fgg-examples:
	go run github.com/rhu1/fgg -fgg -eval=-1 fgg/examples/popl20/booleans/booleans.fgg
	go run github.com/rhu1/fgg -fgg -eval=-1 fgg/examples/popl20/compose/compose.fgg
	go run github.com/rhu1/fgg -fgg -eval=-1 fgg/examples/popl20/graph/graph.fgg
	go run github.com/rhu1/fgg -fgg -eval=-1 fgg/examples/popl20/irregular/irregular.fgg
	go run github.com/rhu1/fgg -fgg -eval=-1 fgg/examples/popl20/map/map.fgg
	go run github.com/rhu1/fgg -fgg -eval=-1 fgg/examples/popl20/monomorph/monomorph.fgg
	go run github.com/rhu1/fgg -fgg -eval=10 fgg/examples/monom/box/box.fgg
	go run github.com/rhu1/fgg -fgg -eval=10 fgg/examples/monom/box/box2.fgg


#program.exe < input.txt > output.txt
#diff correct.txt output.txt

.PHONY: test-monom
test-monom:
	mkdir -p tmp/test/fg/booleans
	go run github.com/rhu1/fgg -fgg -eval=-1 -compile=tmp/test/fg/booleans/booleans.go fgg/examples/popl20/booleans/booleans.fgg
	go run github.com/rhu1/fgg -eval=-1 tmp/test/fg/booleans/booleans.go

	mkdir -p tmp/test/fg/compose
	go run github.com/rhu1/fgg -fgg -eval=-1 -compile=tmp/test/fg/compose/compose.go fgg/examples/popl20/compose/compose.fgg
	go run github.com/rhu1/fgg -eval=-1 tmp/test/fg/compose/compose.go

	mkdir -p tmp/test/fg/graph
	go run github.com/rhu1/fgg -fgg -eval=-1 -compile=tmp/test/fg/graph/graph.go fgg/examples/popl20/graph/graph.fgg
	go run github.com/rhu1/fgg -eval=-1 tmp/test/fg/graph/graph.go

	mkdir -p tmp/test/fg/irregular
	go run github.com/rhu1/fgg -fgg -eval=-1 -compile=tmp/test/fg/irregular/irregular.go fgg/examples/popl20/irregular/irregular.fgg
	go run github.com/rhu1/fgg -eval=-1 tmp/fg/irregular/irregular.go

	mkdir -p tmp/test/fg/map
	go run github.com/rhu1/fgg -fgg -eval=-1 -compile=tmp/test/fg/map/map.go fgg/examples/popl20/map/map.fgg
	go run github.com/rhu1/fgg -eval=-1 tmp/test/fg/map/map.go

	mkdir -p tmp/test/fg/monomorph
	go run github.com/rhu1/fgg -fgg -eval=-1 -compile=tmp/test/fg/monomorph/monomorph.go fgg/examples/popl20/monomorph/monomorph.fgg
	go run github.com/rhu1/fgg -eval=-1 tmp/fg/monomorph/monomorph.go

	mkdir -p tmp/test/fg/box
	go run github.com/rhu1/fgg -fgg -eval=10 -compile=tmp/test/fg/box/box2.go fgg/examples/monom/box/box2.fgg
	go run github.com/rhu1/fgg -eval=10 tmp/test/fg/box/box2.go

# TODO: check simulation of -monom output (same result and number of eval steps)


.PHONY: clean-test-monom
clean-test-monom:
	rm -f tmp/test/fg/booleans/booleans.go
	rm -fd tmp/test/fg/booleans

	rm -f tmp/test/fg/compose/compose.go
	rm -fd tmp/test/fg/compose

	rm -f tmp/test/fg/graph/graph.go
	rm -fd tmp/test/fg/graph

	rm -f tmp/test/fg/irregular/irregular.go
	rm -fd tmp/test/fg/irregular

	rm -f tmp/test/fg/map/map.go
	rm -fd tmp/test/fg/map

	rm -f tmp/test/fg/monomorph/monomorph.go
	rm -fd tmp/test/fg/monomorph

	rm -f tmp/test/fg/box/box2.go
	rm -fd tmp/test/fg/box


.PHONY: test-fg2fgg
test-fg2fgg:
	mkdir -p tmp/test/fgg/compose
	go run github.com/rhu1/fgg/cmd/fg2fgg fg/examples/popl20/compose/compose.go > tmp/test/fgg/compose/compose.fgg
	go run github.com/rhu1/fgg -fgg -eval=-1 tmp/test/fgg/compose/compose.fgg

	mkdir -p tmp/test/fgg/map
	go run github.com/rhu1/fgg/cmd/fg2fgg fg/examples/popl20/map/map.go > tmp/test/fgg/map/map.fgg
	go run github.com/rhu1/fgg -fgg -eval=-1 tmp/test/fgg/map/map.fgg


.PHONY: clean-test-fg2fgg
clean-test-fg2fgg:
	rm -f tmp/test/fgg/compose/compose.fgg
	rm -fd tmp/test/fgg/compose

	rm -f tmp/test/fgg/map/map.fgg
	rm -fd tmp/test/fgg/map


.PHONY: clean-test
clean-test: clean-test-monom clean-test-fg2fgg

