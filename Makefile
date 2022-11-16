API_PROTO_FILES=$(shell find api -name "*.proto")

errors:
	protoc --proto_path=. \
             --proto_path=$(API_PROTO_FILES) \
             --go_out=paths=source_relative:. \
             --go-errors_out=paths=source_relative:. \
             $(API_PROTO_FILES)

