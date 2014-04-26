.PHONY: run

SOURCES = $(wildcard src/*.go)

run:
	go run $(SOURCES)
