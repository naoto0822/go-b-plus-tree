dep:
	brew install graphviz

genpng: clean
	go run main.go > graph.dot
	dot -Tpng -Gdpi=300 graph.dot > graph.png

clean:
	rm graph.dot graph.png test.btr

.PHONY: genpng dep clean
