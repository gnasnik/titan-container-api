export GOPROXY=https://goproxy.cn

all: gen build
.PHONY: all

GOCC?=go

titan-container-api:
	rm -f titan-container-api
	$(GOCC) build $(GOFLAGS) -o titan-container-api .
.PHONY: titan-explorer

gen:
	sqlc generate
.PHONY: gen

build: titan-container-api
.PHONY: build
