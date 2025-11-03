.PHONY: generate clean

# Generate code from proto files
generate:
	@echo "Generating Go code from proto files..."
	@mkdir -p gen/go
	@protoc --go_out=gen/go --go_opt=paths=source_relative \
		--go-grpc_out=gen/go --go-grpc_opt=paths=source_relative \
		-I. \
		*.proto
	@echo "Generation complete!"

# Generate TypeScript code
generate-ts:
	@echo "Generating TypeScript code from proto files..."
	@mkdir -p gen/ts
	@protoc --ts_out=gen/ts \
		--plugin=protoc-gen-ts=./node_modules/.bin/protoc-gen-ts \
		-I. \
		*.proto
	@echo "TypeScript generation complete!"

# Clean generated files
clean:
	@rm -rf gen/
	@echo "Cleaned generated files"

# Install dependencies
install:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

