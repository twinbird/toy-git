package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"
)

func write_tree_object(d *Dircache, p string) ([]byte, error) {
	buf := new(bytes.Buffer)
	b := build_tree_bytes(d)

	// header
	header_str := fmt.Sprintf("tree %d\x00", len(b))
	buf.Write(append([]byte(header_str), b...))

	obj := buf.Bytes()

	// hashing(sha1)
	h := sha1.New()
	h.Write(obj)
	key := h.Sum(nil)

	// store object database
	if err := write_hash_object(key, obj); err != nil {
		return nil, err
	}
	return key, nil
}

func build_tree_bytes(d *Dircache) []byte {
	buf := new(bytes.Buffer)

	for _, e := range d.Entries {
		s := fmt.Sprintf("%o %s\x00", e.Mode, string(e.PathName))
		buf.Write([]byte(s))
		for _, v := range e.Sha1 {
			buf.WriteByte(v)
		}
	}
	return buf.Bytes()
}

func write_tree_cmd() {
	repop, err := find_git_repository(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: not a git repository (or any of the parent directories): %s", repop)
		os.Exit(128)
	}

	d, err := load_dircache(repop)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: internal error: %v\n", err)
		os.Exit(128)
	}

	key, err := write_tree_object(d, repop)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: internal error: %v\n", err)
		os.Exit(128)
	}

	fmt.Printf("%x\n", key)
}
