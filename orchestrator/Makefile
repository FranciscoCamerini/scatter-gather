.PHONY: all

target = orchestrator
srcs = orchestrator.go

all: $(target)

clean:
	rm -f $(target)

$(target): $(srcs)
	go build -o $(target).new
	cmp -s $(target).new $(target) || mv $(target).new $(target)
	rm -f $(target).new
