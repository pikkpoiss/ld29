.PHONY: run

SOURCES = $(wildcard src/*.go)

run:
	cd src; go run $(patsubst src/%,%,$(SOURCES)); cd -
