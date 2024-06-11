CPUS ?= $(shell nproc)
MAKEFLAGS += --jobs=$(CPUS)
.PHONY: stop clean all crawler master

all: master crawler

clean:
	$(MAKE) -C crawler clean
	$(MAKE) -C master clean

stop:
	for f in `ls storage/.*.pid`; do read c < $$f; kill -- $$c; rm $$f; done

crawler:
	$(MAKE) -C crawler

master:
	$(MAKE) -C master
