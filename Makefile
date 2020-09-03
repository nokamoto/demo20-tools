all:
	go install golang.org/x/tools/cmd/goimports
	goimports -d -w .
	go test ./...
	go mod tidy

	docker run -v $(PWD):/protobuf -w /protobuf --rm $$(docker build -q -f Dockerfile.protoc-gen-authz .) protoc --authz_out=paths=source_relative:. -I . test/test.proto
