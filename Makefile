PLUGIN = $(shell ls -d ./plugin/*/)

build:
	go build
all: $(PLUGIN)
	for dir in $(PLUGIN);do \
		cd $$dir && make && cd ../../;\
	done

clean:
	rm gowiki
