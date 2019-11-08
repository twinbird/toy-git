package main

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func hash_obj_cmd(write bool, stdin bool, files []string) {
	if stdin {
		hash_obj_cmd_sub(write, os.Stdin)
		return
	}

	for _, p := range files {
		f, err := os.Open(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: 'git hash-object' target file open failed. %v\n", err)
			os.Exit(1)
		}
		defer f.Close()

		hash_obj_cmd_sub(write, f)
	}
}

func hash_obj_cmd_sub(write bool, f *os.File) {
	sha, err := hash_object(write, f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: 'git hash-object' failed. %v\n", err)
		os.Exit(1)
	}
	if write == false {
		fmt.Printf("%x\n", sha)
	}
}

func hash_object(write bool, f *os.File) ([]byte, error) {
	// contents
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// header
	// 'blob' + ' ' + <contents size> + <null byte>
	header_str := fmt.Sprintf("blob %d\x00", len(content))
	header := []byte(header_str)

	// object
	// header + contents
	obj := append(header, content...)

	// hashing(sha1)
	h := sha1.New()
	h.Write(obj)
	key := h.Sum(nil)

	if write == false {
		return key, nil
	}

	// store object database
	if err := write_hash_object(key, obj); err != nil {
		return nil, err
	}

	return key, nil
}

func write_hash_object(key []byte, obj []byte) error {
	//==========================================
	// Create object store directory
	// dirname is 'key' prefix 2 chars
	//==========================================
	p, err := os.Getwd()
	if err != nil {
		return err
	}
	repo_dir, err := find_git_repository(p)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(128)
	}
	dirname := fmt.Sprintf("%x", key)[:2]

	obj_dir := filepath.Join(repo_dir, "objects", dirname)
	if err := os.Mkdir(obj_dir, 0755); err != nil && os.IsExist(err) == false {
		return err
	}

	//==========================================
	// Store object
	// filename is 'key' suffix 38 chars
	//==========================================
	fname := fmt.Sprintf("%x", key)[2:]
	fpath := filepath.Join(obj_dir, fname)

	// compress by zlib
	var b bytes.Buffer
	w, err := zlib.NewWriterLevel(&b, flate.BestSpeed) // default compression level of loose object is BestSpeed
	if err != nil {
		return err
	}
	w.Write(obj)
	w.Close()

	// store to file
	err = ioutil.WriteFile(fpath, b.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
