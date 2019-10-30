toy-git: main.go
	go build

.PHONY: test
.SILENT:
test: clean toy-git
	test/init_test.sh

.PHONY: clean
clean:
	rm -f toy-git
	rm -rf .toy-git
