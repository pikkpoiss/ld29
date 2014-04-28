.PHONY: build clean run

PROJECT = heavydrizzle
SOURCES = $(wildcard src/*.go)

BASEBUILD = build/$(PROJECT)-osx
OSXBUILD = $(BASEBUILD)/$(PROJECT).app/Contents

VERSION = 0.1

clean:
	rm -rf build

$(OSXBUILD)/MacOS/$(PROJECT): $(SOURCES)
	mkdir -p $(dir $@)
	go build -o $@ src/*.go

#build/$(PROJECT)-osx-$(VERSION).zip: \
#	$(OSXBUILD)/Info.plist \
#	$(OSXBUILD)/MacOS/launch.sh \
#	$(BASEBUILD)/README \


#build: build/$(PROJECT)-osx-$(VERSION).zip
build: $(OSXBUILD/MacOS/$(PROJECT)

run: build
	$(OSXBUILD)/MacOS/launch.sh
#run:
#	cd src; go run $(patsubst src/%,%,$(SOURCES)); cd -
