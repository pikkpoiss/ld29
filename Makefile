.PHONY: build clean run

PROJECT = heavydrizzle
SOURCES = $(wildcard src/*.go)
RUNTIME_ASSETS = $(wildcard src/assets/*)
ICON_ASSETS = $(wildcard assets/*.icns)

BASEBUILD = build/$(PROJECT)-osx
OSXBUILD = $(BASEBUILD)/$(PROJECT).app/Contents

VERSION = $(shell cat VERSION)
REPLACE = s/9\.9\.9/$(VERSION)/g

clean:
	rm -rf build

$(OSXBUILD)/MacOS/launch.sh: scripts/launch.sh
	mkdir -p $(dir $@)
	cp $< $@

$(OSXBUILD)/Info.plist: pkg/osx/Info.plist
	mkdir -p $(OSXBUILD)
	sed $(REPLACE) $< > $@

$(OSXBUILD)/MacOS/$(PROJECT): $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go

$(OSXBUILD)/Resources/%.icns: assets/%.icns
	mkdir -p $(dir $@)
	cp $< $@

$(OSXBUILD)/Resources/assets/%: src/assets/%
	mkdir -p $(dir $@)
	cp -R $< $@

build/$(PROJECT)-osx-$(VERSION).zip: \
	$(OSXBUILD)/MacOS/launch.sh \
	$(OSXBUILD)/MacOS/$(PROJECT) \
	$(OSXBUILD)/Info.plist \
	$(subst src/assets/,$(OSXBUILD)/Resources/assets/,$(RUNTIME_ASSETS)) \
	$(subst assets/,$(OSXBUILD)/Resources/,$(ICON_ASSETS)) 
	cd build && zip -r $(notdir $@) $(PROJECT)-osx

build: build/$(PROJECT)-osx-$(VERSION).zip

run: build
	$(OSXBUILD)/MacOS/launch.sh
