INSTALL_PATH ?= /usr/local
INSTALL = /usr/bin/install -v -C
REV := $(shell git rev-parse --short HEAD)
GOARCH ?= amd64

all: build

build:
	go build -o gowiki

cross-compile:
	GOOS=darwin GOARCH=$(GOARCH) go build -o gowiki-osx
	GOOS=linux GOARCH=$(GOARCH) go build -o gowiki-linux
	GOOS=windows GOARCH=$(GOARCH) go build -o gowiki-windows

clean:
	rm -vf gowiki gowiki-* gowiki-*.tar.gz

install: gowiki
	@mkdir -p $(INSTALL_PATH)/bin
	@mkdir -p $(INSTALL_PATH)/share/gowiki/public/{css,js}
	@mkdir -p $(INSTALL_PATH)/share/gowiki/templates

	$(INSTALL) gowiki $(INSTALL_PATH)/bin/
	$(INSTALL) -m 0644 public/css/*.css $(INSTALL_PATH)/share/gowiki/public/css/
	$(INSTALL) -m 0644 public/js/*.js $(INSTALL_PATH)/share/gowiki/public/js/
	$(INSTALL) -m 0644 templates/*.html $(INSTALL_PATH)/share/gowiki/templates/

uninstall:
	rm -rvf $(INSTALL_PATH)/share/gowiki
	rm -vf $(INSTALL_PATH)/bin/gowiki

release: build
	@tar cvzf gowiki-$(REV).tar.gz gowiki public templates
