.DEFAULT_GOAL := help
VERSION := v0.13.0

fmt:
	gofmt -l .

help:
	@echo 'Makefile for Crypt2go'
	@echo
	@echo 'Usage:'
	@echo '    make fmt        Format Go files'
	@echo '    make help       Display this help message'
	@echo '    make run        Execute program examples'
	@echo '    make tag        GitHub tag (after manual git commit)'
	@echo '    make test       Execute tests'
	@echo '    make version    Display current package version'

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

update:
	go get -u -d ./...
