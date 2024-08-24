.DEFAULT_GOAL := help
VERSION := v1.8.0

fmt:
	gofmt -l .

help:
	@echo 'Makefile for Crypt2go'
	@echo
	@echo 'Usage:'
	@echo '    make help       Display this help message'
	@echo '    make version    Display current package version'
	@echo '    make update     Update dependencies'
	@echo '    make release    Perform all checks (fmt, test, and lint)'
	@echo '    make fmt        Format Go files'
	@echo '    make test       Execute tests'
	@echo '    make lint       Lint code with golangci-lint'
	@echo '    make run        Execute program examples'
	@echo '    make tag        GitHub tag (after manual git commit)'

lint:
	golangci-lint run

release: update fmt test lint vet

run:
	go run examples/aes/main.go
	go run examples/blowfish/main.go

tag:
	git push
	git tag -a ${VERSION} -m 'Version ${VERSION}'
	git push --tags

test:
	go test -v ./...

version:
	@echo 'Crypt2go version: ${VERSION}'

vet:
	go list ./...
	go vet ./...

update:
	go get -u ./...
	go mod tidy
