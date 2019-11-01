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
