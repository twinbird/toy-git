package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	REPOSITORY_DIR_NAME = ".toy-git"
)

func main() {
	hash_obj_flag := flag.NewFlagSet("hash-object", flag.ExitOnError)
	cat_file_flag := flag.NewFlagSet("cat-file", flag.ExitOnError)
	update_index_flag := flag.NewFlagSet("update-index", flag.ExitOnError)
	ls_files_flag := flag.NewFlagSet("ls-files", flag.ExitOnError)
	commit_tree_flag := flag.NewFlagSet("commit-tree", flag.ExitOnError)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `toy-git 
Git Plumbing commands implementation for learning git internal.
You **MUST NOT** use your production repository.
(toy-git manipulates the .toy-git directory)

toy-git is sets of subcommands.
 * toy-git init
 * toy-git hash-object
 * toy-git cat-file
 * toy-git update-index
 * toy-git ls-files
 * toy-git write-tree
 * toy-git commit-tree
 * toy-git update-ref

See also each subcommands help.

Ex) toy-git hash-object -h

`)
	}

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	switch os.Args[1] {
	case "init":
		init_cmd()
	case "hash-object":
		w := hash_obj_flag.Bool("w", false, "Actually write the object into the object database.")
		stdin := hash_obj_flag.Bool("stdin", false, "Read the object from standard input instead of from a file.")
		hash_obj_flag.Parse(os.Args[2:])

		hash_obj_cmd(*w, *stdin, hash_obj_flag.Args())
	case "cat-file":
		t := cat_file_flag.Bool("t", false, "Instead of the content, show the object type identified by <object>.")
		s := cat_file_flag.Bool("s", false, "Instead of the content, show the object size identified by <object>.")
		p := cat_file_flag.Bool("p", false, "Pretty-print the contents of <object> based on its type.")
		cat_file_flag.Parse(os.Args[2:])

		if *t == false && *s == false && *p == false {
			cat_file_flag.Usage()
			return
		}

		if len(cat_file_flag.Args()) < 1 {
			cat_file_flag.Usage()
			return
		}

		cat_file_cmd(*t, *s, *p, cat_file_flag.Args())
	case "update-index":
		add := update_index_flag.Bool("add", false, "If a specified file isn't in the index already then it's added. Default behaviour is to ignore new files.")
		remove := update_index_flag.Bool("remove", false, "If a specified file is in the index but is missing then it's removed. Default behavior is to ignore removed file.")
		update_index_flag.Parse(os.Args[2:])

		if *add == true && *remove == true {
			update_index_flag.Usage()
			return
		}

		if len(update_index_flag.Args()) < 1 {
			update_index_flag.Usage()
			return
		}

		update_index_cmd(*add, *remove, update_index_flag.Args())
	case "ls-files":
		cached := ls_files_flag.Bool("c", false, "Show cached files in the output (default)")
		deleted := ls_files_flag.Bool("d", false, "Show deleted files in the output")
		modified := ls_files_flag.Bool("m", false, "Show modified files in the output")
		ls_files_flag.Parse(os.Args[2:])

		if *cached == false && *deleted == false && *modified == false {
			*cached = true
		}

		ls_files_cmd(*cached, *deleted, *modified)
	case "write-tree":
		write_tree_cmd()
	case "commit-tree":
		if len(os.Args) < 3 {
			commit_tree_flag.Usage()
			return
		}

		tree_sha := os.Args[2]
		parent := commit_tree_flag.String("p", "", "Each -p indicates the id of a parent commit object")
		commit_tree_flag.Parse(os.Args[3:])

		commit_tree_cmd(tree_sha, *parent)
	case "update-ref":
		if len(os.Args) < 4 {
			fmt.Fprintf(os.Stderr, "toy-git update-ref <ref> <newvalue>")
			return
		}

		update_ref_cmd(os.Args[2], os.Args[3])
	default:
		flag.Usage()
	}
}

func find_git_repository(path string) (string, error) {
	repo_path := filepath.Join(path, REPOSITORY_DIR_NAME)
	if _, err := os.Stat(repo_path); os.IsNotExist(err) {
		if path == filepath.Dir(path) {
			return "", fmt.Errorf("fatal: not a git repository (or any of the parent directories): " + REPOSITORY_DIR_NAME)
		}

		return find_git_repository(filepath.Dir(path))
	}
	return repo_path, nil
}
