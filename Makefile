GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
API_PROTO_FILES=$(shell find api -name "*.proto")
KRATOS_VERSION=$(shell go mod graph |grep go-kratos/kratos/v2 |head -n 1 |awk -F '@' '{print $$2}')
KRATOS=$(GOPATH)/pkg/mod/github.com/go-kratos/kratos/v2@$(KRATOS_VERSION)

COMPILE_TARGET="./"

.PHONY: all

all:
	make errors;
	make validate;


.PHONY: errors
errors:
	protoc --proto_path=. \
           	 --proto_path=./third_party \
             --proto_path=$(API_PROTO_FILES) \
             --go_out=paths=source_relative:. \
             --go-errors_out=paths=source_relative:. \
             $(API_PROTO_FILES)


.PHONY: validate
# generate validate code
validate:
	protoc --proto_path=. \
           --proto_path=./third_party \
           --go_out=paths=source_relative:$(COMPILE_TARGET)  \
           --validate_out=paths=source_relative,lang=go:$(COMPILE_TARGET) \
           $(API_PROTO_FILES)
