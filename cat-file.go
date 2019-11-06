package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type GitObject interface {
	obj_type() string
	obj_size() int
	obj_string() string
}

type BlobObject struct {
	type_str string
	size     int
	data     []byte
}

func (obj BlobObject) obj_type() string {
	return obj.type_str
}

func (obj BlobObject) obj_size() int {
	return obj.size
}

func (obj BlobObject) obj_string() string {
	return string(obj.data)
}

func cat_file_cmd(opt_t bool, opt_s bool, opt_p bool, sha_strs []string) {
	repo_path, err := find_git_repository(".")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(128)
	}

	for _, s := range sha_strs {
		// open object
		p := filepath.Join(repo_path, "objects", s[:2], s[2:])
		f, err := os.Open(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal: Not a valid object name %s\n%v\n", s, err)
			os.Exit(128)
		}
		defer f.Close()

		// uncompress zlib
		zreader, err := zlib.NewReader(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal: zlib uncompress failed.")
			os.Exit(1)
		}
		defer zreader.Close()

		// read file
		obj, err := read_object(zreader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal: Not a valid object name %s\n%v\n", s, err)
			os.Exit(128)
		}

		// print type
		if opt_t == true {
			fmt.Println(obj.obj_type())
		}
		// print size
		if opt_s == true {
			fmt.Println(obj.obj_size())
		}
		// pretty print
		if opt_p == true {
			fmt.Println(obj.obj_string())
		}
	}
}

func read_file_type(f io.Reader) (string, error) {
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	sep := bytes.IndexAny(b, " ")
	if sep < 0 {
		return "", fmt.Errorf("Type not found.%v", b)
	}
	return string(b[:sep]), nil
}

func read_object(f io.Reader) (GitObject, error) {
	var obj BlobObject

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// type
	type_sep := bytes.IndexAny(b, " ")
	if type_sep < 0 {
		return nil, fmt.Errorf("Type not found.%v", b)
	}
	obj.type_str = string(b[:type_sep])

	// size
	size_sep := bytes.Index(b, []byte("\x00"))
	if size_sep < 0 {
		return nil, fmt.Errorf("Size not found.%v", b)
	}
	obj.size, err = strconv.Atoi(string(b[type_sep+1 : size_sep]))
	if err != nil {
		return nil, fmt.Errorf("Size not found.%v", err)
	}

	// blob
	obj.data = b[size_sep+1 : size_sep+1+obj.size-1]

	return obj, nil
}
