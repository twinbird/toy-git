toy-git: *.go
	go build

.PHONY: test
.SILENT:
test: clean toy-git
	test/init_test.sh
	test/hash_cat_test.sh
	test/update-index.sh

.PHONY: clean
clean:
	rm -f toy-git
	rm -rf .toy-git
	rm -rf test/.toy-git
	-unlink test/.git 2>/dev/null
	rm -f test/[a-z].txt
