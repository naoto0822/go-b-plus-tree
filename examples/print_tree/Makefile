dep:
	brew install graphviz

genpng:
	go run main.go > graph.dot
	dot -Tpng -Gdpi=300 graph.dot > graph.png

.PHONY: genpng dep
