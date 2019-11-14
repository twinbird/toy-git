package main

import (
	"fmt"
	"os"
)

func ls_files_cmd(cached bool, deleted bool, modified bool) {
	repop, err := find_git_repository(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: not a git repository (or any of the parent directories): %s", repop)
		os.Exit(128)
	}

	// read dircache
	d, err := load_dircache(repop)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: internal error: %v\n", err)
		os.Exit(128)
	}

	print_dircache(d, cached, deleted, modified)
}

func print_dircache(d *Dircache, cached bool, deleted bool, modified bool) {
	for _, e := range d.Entries {
		if cached {
			fmt.Println(string(e.PathName))
		}
	}
}
