package main

import (
	"bytes"
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

	if err := print_dircache(d, cached, deleted, modified); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: internal error: %v\n", err)
		os.Exit(128)
	}
}

func is_deleted(path string) (bool, error) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return true, nil
	} else if err != nil {
		return false, err
	}
	return false, nil
}

func is_modified(e *DircacheEntry) (bool, error) {
	path := e.PathName

	f, err := os.Open(string(path))
	if err != nil && os.IsNotExist(err) {
		return true, nil
	} else if err != nil {
		return false, err
	}

	b, err := hash_object(false, f)
	if err != nil {
		return false, err
	}

	return !bytes.Equal(b, e.Sha1[:]), nil
}

func print_dircache(d *Dircache, cached bool, deleted bool, modified bool) error {
	for _, e := range d.Entries {
		if cached {
			fmt.Println(string(e.PathName))
		}

		if deleted {
			d, err := is_deleted(string(e.PathName))
			if err != nil {
				return err
			}
			if d {
				fmt.Println(string(e.PathName))
			}
		}

		if modified {
			m, err := is_modified(e)
			if err != nil {
				return err
			}
			if m {
				fmt.Println(string(e.PathName))
			}
		}
	}
	return nil
}
