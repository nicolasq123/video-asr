#VERSION=$(shell echo "$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)")
#GOBUILD=go build -ldflags "-w -s -X github.com/nicolasq123/videoasr.version=$(VERSION)"
GOBUILD=go build -ldflags "-w -s "


build:
	$(GOBUILD) -o bin/videoasr/videoasr github.com/nicolasq123/videoasr/cmd/main/

run:
	bin/videoasr/videoasr --conf "./cmd/main/conf.yml"

build_run: \
	build \
	run \
