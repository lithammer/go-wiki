INSTALL_PATH ?= /usr/local
INSTALL = /usr/bin/install -v -C
REV := $(shell git rev-parse --short HEAD)

all: build

build:
	go build -o gowiki

clean:
	/bin/rm -vf gowiki *.tar.gz

install: gowiki
	@/bin/mkdir -p $(INSTALL_PATH)/bin
	@/bin/mkdir -p $(INSTALL_PATH)/share/gowiki/public/{css,js}
	@/bin/mkdir -p $(INSTALL_PATH)/share/gowiki/templates

	$(INSTALL) gowiki $(INSTALL_PATH)/bin/
	$(INSTALL) -m 0644 public/css/*.css $(INSTALL_PATH)/share/gowiki/public/css/
	$(INSTALL) -m 0644 public/js/*.js $(INSTALL_PATH)/share/gowiki/public/js/
	$(INSTALL) -m 0644 templates/*.html $(INSTALL_PATH)/share/gowiki/templates/

uninstall:
	/bin/rm -rvf $(INSTALL_PATH)/share/gowiki
	/bin/rm -vf $(INSTALL_PATH)/bin/gowiki

release: build
	@/usr/bin/tar cvzf gowiki-$(REV).tar.gz gowiki public templates
