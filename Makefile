CPUS ?= $(shell nproc)
MAKEFLAGS += --jobs=$(CPUS)
.PHONY: stop clean all worker orchestrator

all: orchestrator worker

clean:
	$(MAKE) -C worker clean
	$(MAKE) -C orchestrator clean

stop:
	for f in `ls storage/.*.pid`; do read c < $$f; kill -- $$c; rm $$f; done

worker:
	$(MAKE) -C worker

orchestrator:
	$(MAKE) -C orchestrator
