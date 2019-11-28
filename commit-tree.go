package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	AUTHOR          = "twinbird"
	AUTHOR_EMAIL    = "ixa2063@gmail.com"
	COMMITTER       = "twinbird"
	COMMITTER_EMAIL = "ixa2063@gmail.com"
	TIMEZONE        = "+0900"
)

func build_commit_bytes(tree string, parent string, message string) []byte {
	buf := new(bytes.Buffer)

	t := time.Now()
	tstr := fmt.Sprintf("%d", t.Unix())

	buf.Write([]byte(fmt.Sprintf("tree %s\n", tree)))
	if len(parent) > 0 {
		buf.Write([]byte(fmt.Sprintf("parent %s\n", parent)))
	}
	buf.Write([]byte(fmt.Sprintf("author %s <%s> %s %s\n", AUTHOR, AUTHOR_EMAIL, tstr, TIMEZONE)))
	buf.Write([]byte(fmt.Sprintf("committer %s <%s> %s %s\n", COMMITTER, COMMITTER_EMAIL, tstr, TIMEZONE)))
	buf.Write([]byte(fmt.Sprintf("\n")))
	buf.Write([]byte(fmt.Sprintf("%s", message)))

	return buf.Bytes()
}

func commit_tree_object(tree string, parent string, message string, repop string) ([]byte, error) {
	buf := new(bytes.Buffer)
	b := build_commit_bytes(tree, parent, message)

	// header
	header_str := fmt.Sprintf("commit %d\x00", len(b))
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

func commit_tree_cmd(tree string, parent string) {
	repop, err := find_git_repository(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: not a git repository (or any of the parent directories): %s", repop)
		os.Exit(128)
	}

	message, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: internal error: %v\n", err)
		os.Exit(128)
	}

	key, err := commit_tree_object(tree, parent, string(message), repop)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: internal error: %v\n", err)
		os.Exit(128)
	}

	fmt.Printf("%x\n", key)
}
