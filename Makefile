all:
	go install golang.org/x/tools/cmd/goimports
	goimports -d -w .
	go test ./...
	go mod tidy
