.DEFAULT_GOAL := help
VERSION := v0.11.0

clean:
	rm examples/aes/aes
	rm examples/blowfish/blowfish

fmt:
	gofmt -l

help:
	@echo 'Makefile for Crypt2go'
	@echo
	@echo 'Usage:'
	@echo '    make clean      Delete executables'
	@echo '    make fmt        Format Go files'
	@echo '    make help       Display this help message'
	@echo '    make tag        GitHub tag (after manual git commit)'
	@echo '    make test       Execute tests'
	@echo '    make version    Display current package version'

tag:
	git push
	git tag -a ${VERSION} -m 'Version ${VERSION}'
	git push --tags

test:
	go test -v ./...

version:
	@echo 'Crypt2go version: ${VERSION}'

.PHONY: check clean deploy fmt help lint test version
