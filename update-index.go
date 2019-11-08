// See Also:
// https://github.com/git/git/blob/master/Documentation/technical/index-format.txt
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"syscall"

	"gopkg.in/djherbis/times.v1"
)

type Dircache struct {
	Header  DircacheHeader
	Entries []*DircacheEntry
}

type DircacheHeader struct {
	Signature       [4]byte
	Version         int32
	NumberOfEntries int32
}

type DircacheEntry struct {
	CTimeSeconds     uint32
	CTimeNanoSeconds uint32
	MTimeSeconds     uint32
	MTimeNanoSeconds uint32
	Dev              int32
	Inode            int32
	Mode             int32
	UID              uint32
	GID              uint32
	Size             int32
	Sha1             [20]byte
	Flags            int16  // [1-bit: assume-valid flag] [1-bit: extended flag(must be zero)] [2-bit: stage(during merge)] [12-bit: name length]
	PathName         []byte // variable length. size is 'Size'
	ZeroPaddingSize  int    // for 8 byte alignment
}

func load_dircache(path string) (*Dircache, error) {
	d := &Dircache{}
	d.Header.Signature = [4]byte{'D', 'I', 'R', 'C'}
	d.Header.Version = 2
	d.Header.NumberOfEntries = 0
	return d, nil
}

func build_dircache_bytes(d *Dircache) []byte {
	buf := new(bytes.Buffer)

	// Header
	binary.Write(buf, binary.LittleEndian, d.Header.Signature)
	binary.Write(buf, binary.LittleEndian, d.Header.Version)
	binary.Write(buf, binary.LittleEndian, d.Header.NumberOfEntries)

	// Entries
	for _, e := range d.Entries {
		s := 0
		binary.Write(buf, binary.LittleEndian, e.CTimeSeconds)
		s += 4
		binary.Write(buf, binary.LittleEndian, e.CTimeNanoSeconds)
		s += 4
		binary.Write(buf, binary.LittleEndian, e.MTimeSeconds)
		s += 4
		binary.Write(buf, binary.LittleEndian, e.MTimeNanoSeconds)
		s += 4
		binary.Write(buf, binary.LittleEndian, e.Dev)
		s += 4
		binary.Write(buf, binary.LittleEndian, e.Inode)
		s += 4
		binary.Write(buf, binary.LittleEndian, e.Mode)
		s += 4
		binary.Write(buf, binary.LittleEndian, e.UID)
		s += 4
		binary.Write(buf, binary.LittleEndian, e.GID)
		s += 4
		binary.Write(buf, binary.LittleEndian, e.Size)
		s += 4
		binary.Write(buf, binary.LittleEndian, e.Sha1)
		s += 20
		binary.Write(buf, binary.LittleEndian, e.Flags)
		s += 2
		binary.Write(buf, binary.LittleEndian, e.PathName)
		s += len(e.PathName)

		// padding
		for i := 0; i < (s % 8); i++ {
			buf.WriteByte(byte(0))
		}
	}

	return buf.Bytes()
}

func write_dircache(d *Dircache, repop string) error {
	// sort by filename

	b := build_dircache_bytes(d)

	indexp := filepath.Join(repop, "index")
	err := ioutil.WriteFile(indexp, b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func update_index_cmd(do_add bool, do_remove bool, paths []string) {
	repop, err := find_git_repository(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: not a git repository (or any of the parent directories): %s", repop)
		os.Exit(128)
	}

	// read dircache
	d, err := load_dircache(repop)

	// update or add or remove dircache
	for _, p := range paths {
		if do_add {
			add_dircache(d, p)
		} else if do_remove {
			remove_dircache(d, p)
		} else {
			update_dircache(d, p)
		}
	}

	// write dircache
	err = write_dircache(d, repop)
}

// update and remove error messages
// error: update-index.go: cannot add to the index - missing --add option?
// fatal: Unable to process path update-index.go

func add_dircache(d *Dircache, path string) {
	// file stat
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s: does not exist and --remove not passed\n", path)
		fmt.Fprintf(os.Stderr, "fatal: Unable to process path %s\n", path)
		os.Exit(128)
	}

	// file times
	t, err := times.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s: does not exist and --remove not passed\n", path)
		fmt.Fprintf(os.Stderr, "fatal: Unable to process path %s\n", path)
		os.Exit(128)
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s: does not exist and --remove not passed\n", path)
		fmt.Fprintf(os.Stderr, "fatal: Unable to process path %s\n", path)
		os.Exit(128)
	}
	defer f.Close()
	sha, err := hash_object(true, f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s: does not exist and --remove not passed\n", path)
		fmt.Fprintf(os.Stderr, "fatal: Unable to process path %s\n", path)
		os.Exit(128)
	}

	e := &DircacheEntry{}

	switch runtime.GOOS {
	case "windows":
		// TODO
	case "darwing":
		// TODO
		//internal_info := info.Sys().(*syscall.Stat_t)
	default:
		internal_info := info.Sys().(*syscall.Stat_t)

		ctime := t.ChangeTime()
		e.CTimeSeconds = uint32(ctime.Unix())
		e.CTimeNanoSeconds = uint32(ctime.UnixNano() - ctime.Unix())

		mtime := t.ModTime()
		e.MTimeSeconds = uint32(mtime.Unix())
		e.MTimeNanoSeconds = uint32(mtime.UnixNano() - mtime.Unix())

		e.Dev = int32(internal_info.Dev)
		e.Inode = int32(internal_info.Ino)
		e.Mode = int32(info.Mode())
		e.UID = internal_info.Uid
		e.GID = internal_info.Gid
		e.Size = int32(info.Size())
		for i, v := range sha {
			e.Sha1[i] = v
		}
		e.Flags = 0 // [1-bit: assume-valid flag] [1-bit: extended flag(must be zero)] [2-bit: stage(during merge)] [12-bit: name length]
		e.PathName = []byte(path)
		d.Entries = append(d.Entries, e)
	}
}

func update_dircache(d *Dircache, path string) {
}

func remove_dircache(d *Dircache, path string) {
}
