.PHONY: build clean run

PROJECT = heavydrizzle
SOURCES = $(wildcard src/*.go)
RUNTIME_ASSETS = $(wildcard src/assets/*)

BASEBUILD = build/$(PROJECT)-osx
OSXBUILD = $(BASEBUILD)/$(PROJECT).app/Contents

VERSION = 0.1

clean:
	rm -rf build

$(OSXBUILD)/MacOS/launch.sh: scripts/launch.sh
	mkdir -p $(dir $@)
	cp $< $@

$(OSXBUILD)/MacOS/$(PROJECT): $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go

$(OSXBUILD)/Resources/assets/%: src/assets/%
	mkdir -p $(dir $@)
	cp -R $< $@

build/$(PROJECT)-osx-$(VERSION).zip: \
	$(OSXBUILD)/MacOS/launch.sh \
	$(OSXBUILD)/MacOS/$(PROJECT) \
	$(subst src/assets/,$(OSXBUILD)/Resources/assets/,$(RUNTIME_ASSETS))
	cd build && zip -r $(notdir $@) $(PROJECT)-osx

build: build/$(PROJECT)-osx-$(VERSION).zip

run: build
	$(OSXBUILD)/MacOS/launch.sh
