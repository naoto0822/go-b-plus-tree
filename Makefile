test:
	go test -v ./...

dep:
	brew install graphviz

.PHONY: test
