package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"
	"strings"
)

type TreeEntry interface {
	GetName() string
	RecordBytes() []byte
}

type FileEntry struct {
	Name string
	Hash [20]byte
	Mode uint32
}

func (f *FileEntry) GetName() string {
	return f.Name
}

func (f *FileEntry) RecordBytes() []byte {
	buf := new(bytes.Buffer)

	s := fmt.Sprintf("%o %s\x00", f.Mode, string(f.Name))
	buf.Write([]byte(s))
	for _, v := range f.Hash {
		buf.WriteByte(v)
	}
	return buf.Bytes()
}

type DirectoryEntry struct {
	Name    string
	Hash    [20]byte
	Mode    uint32
	Entries []TreeEntry
}

func (d *DirectoryEntry) GetName() string {
	return d.Name
}

func (d *DirectoryEntry) RecordBytes() []byte {
	buf := new(bytes.Buffer)

	s := fmt.Sprintf("%o %s\x00", d.Mode, string(d.Name))
	buf.Write([]byte(s))
	for _, v := range d.Hash {
		buf.WriteByte(v)
	}
	return buf.Bytes()
}

func addPath(d *DirectoryEntry, p string, sha [20]byte, mode uint32) {
	if strings.Contains(p, "/") {
		// dir
		left := strings.Split(p, "/")[0]

		var found bool
		for _, e := range d.Entries {
			if e.GetName() == left {
				st := strings.Index(p, "/") + 1
				t := e.(*DirectoryEntry)
				addPath(t, p[st:], sha, mode)
				found = true
				break
			}
		}

		if found == false {
			// new dir
			newd := &DirectoryEntry{}
			newd.Name = left
			// newd.Hash is set after
			newd.Mode = 040000
			d.Entries = append(d.Entries, newd)

			st := strings.Index(p, "/") + 1
			addPath(newd, p[st:], sha, mode)
		}
	} else {
		// file
		fe := &FileEntry{}
		fe.Name = p
		fe.Mode = mode
		fe.Hash = sha
		d.Entries = append(d.Entries, fe)
	}
}

func build_tree(d *Dircache) *DirectoryEntry {
	root := &DirectoryEntry{}
	root.Name = "root"

	for _, e := range d.Entries {
		addPath(root, string(e.PathName), e.Sha1, e.Mode)
	}
	return root
}

func print_tree(d *DirectoryEntry, nest int) {
	for _, e := range d.Entries {
		for i := 0; i < nest; i++ {
			fmt.Printf("\t")
		}
		fmt.Println(e.GetName())
		x, ok := e.(*DirectoryEntry)
		if ok {
			print_tree(x, nest+1)
		}
	}
}

func write_tree_object(d *DirectoryEntry) ([]byte, error) {
	buf := new(bytes.Buffer)
	b, err := build_tree_bytes(d)
	if err != nil {
		return nil, err
	}

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

func build_tree_bytes(d *DirectoryEntry) ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, e := range d.Entries {
		x, ok := e.(*DirectoryEntry)
		if ok {
			sha, err := write_tree_object(x)
			if err != nil {
				return nil, err
			}
			for i, v := range sha {
				x.Hash[i] = v
			}
		}

		buf.Write(e.RecordBytes())
	}
	return buf.Bytes(), nil
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

	t := build_tree(d)

	key, err := write_tree_object(t)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: internal error: %v\n", err)
		os.Exit(128)
	}

	fmt.Printf("%x\n", key)
}
