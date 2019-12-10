package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func update_ref_cmd(ref string, nvalue string) {
	p, err := find_git_repository(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: not a git repository (or any of the parent directories): "+REPOSITORY_DIR_NAME)
	}

	if check_new_value(p, nvalue) == false {
		fmt.Fprintf(os.Stderr, "%s is invalid git commit object.\n", nvalue)
		os.Exit(128)
	}

	if err := update_ref(p, ref, nvalue); err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		os.Exit(128)
	}
}

func check_new_value(repo string, nvalue string) bool {
	// [TODO] check git object type
	p := filepath.Join(repo, "objects", nvalue[:2], nvalue[2:])
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	} else if err != nil {
		return false
	}
	return true
}

func update_ref(repo string, ref string, nvalue string) error {
	p := filepath.Join(repo, ref)
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "%s", nvalue)
	defer f.Close()

	return nil
}
