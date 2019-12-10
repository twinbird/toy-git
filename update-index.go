// See Also:
// https://github.com/git/git/blob/master/Documentation/technical/index-format.txt
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"syscall"
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
	Dev              uint32
	Inode            uint32
	Mode             uint32
	UID              uint32
	GID              uint32
	Size             uint32
	Sha1             [20]byte
	Flags            uint16 // [1-bit: assume-valid flag] [1-bit: extended flag(must be zero)] [2-bit: stage(during merge)] [12-bit: name length]
	PathName         []byte // variable length. size is 'Size'
	ZeroPaddingSize  int    // for 8 byte alignment
}

func load_dircache(path string) (*Dircache, error) {
	contents, err := ioutil.ReadFile(filepath.Join(path, "index"))
	if err != nil && os.IsNotExist(err) {
		// return new dircache
		d := &Dircache{}
		d.Header.Signature = [4]byte{'D', 'I', 'R', 'C'}
		d.Header.Version = 2
		d.Header.NumberOfEntries = 0
		return d, nil
	} else if err != nil {
		return nil, err
	}

	// read exist dircache
	return read_dircache_bytes(contents)
}

func read_dircache_bytes(b []byte) (*Dircache, error) {
	buf := bytes.NewReader(b)
	var d Dircache

	// Header
	binary.Read(buf, binary.BigEndian, &d.Header.Signature)
	binary.Read(buf, binary.BigEndian, &d.Header.Version)
	binary.Read(buf, binary.BigEndian, &d.Header.NumberOfEntries)

	// Entries
	for i := int32(0); i < d.Header.NumberOfEntries; i++ {
		var e DircacheEntry
		n := 0
		if err := binary.Read(buf, binary.BigEndian, &e.CTimeSeconds); err != nil {
			return nil, err
		}
		n += 4
		if err := binary.Read(buf, binary.BigEndian, &e.CTimeNanoSeconds); err != nil {
			return nil, err
		}
		n += 4
		if err := binary.Read(buf, binary.BigEndian, &e.MTimeSeconds); err != nil {
			return nil, err
		}
		n += 4
		if err := binary.Read(buf, binary.BigEndian, &e.MTimeNanoSeconds); err != nil {
			return nil, err
		}
		n += 4
		if err := binary.Read(buf, binary.BigEndian, &e.Dev); err != nil {
			return nil, err
		}
		n += 4
		if err := binary.Read(buf, binary.BigEndian, &e.Inode); err != nil {
			return nil, err
		}
		n += 4
		if err := binary.Read(buf, binary.BigEndian, &e.Mode); err != nil {
			return nil, err
		}
		n += 4
		if err := binary.Read(buf, binary.BigEndian, &e.UID); err != nil {
			return nil, err
		}
		n += 4
		if err := binary.Read(buf, binary.BigEndian, &e.GID); err != nil {
			return nil, err
		}
		n += 4
		if err := binary.Read(buf, binary.BigEndian, &e.Size); err != nil {
			return nil, err
		}
		n += 4
		if err := binary.Read(buf, binary.BigEndian, &e.Sha1); err != nil {
			return nil, err
		}
		n += 20
		if err := binary.Read(buf, binary.BigEndian, &e.Flags); err != nil {
			return nil, err
		}
		n += 2

		name_len := e.Flags
		name_len &= uint16(0b00001111111111111111) // [12-bit: name length]

		b := make([]byte, name_len)
		if _, err := io.ReadAtLeast(buf, b, int(name_len)); err != nil {
			return nil, err
		}
		e.PathName = b
		n += int(name_len)

		// padding
		for k := 0; k < (8 - n%8); k++ {
			var padding byte
			binary.Read(buf, binary.BigEndian, padding)
		}

		d.Entries = append(d.Entries, &e)
	}

	return &d, nil
}

func build_dircache_bytes(d *Dircache) []byte {
	buf := new(bytes.Buffer)

	// Header
	binary.Write(buf, binary.BigEndian, d.Header.Signature)
	binary.Write(buf, binary.BigEndian, d.Header.Version)
	binary.Write(buf, binary.BigEndian, d.Header.NumberOfEntries)

	// Entries
	for _, e := range d.Entries {
		s := 0
		binary.Write(buf, binary.BigEndian, e.CTimeSeconds)
		s += 4
		binary.Write(buf, binary.BigEndian, e.CTimeNanoSeconds)
		s += 4
		binary.Write(buf, binary.BigEndian, e.MTimeSeconds)
		s += 4
		binary.Write(buf, binary.BigEndian, e.MTimeNanoSeconds)
		s += 4
		binary.Write(buf, binary.BigEndian, e.Dev)
		s += 4
		binary.Write(buf, binary.BigEndian, e.Inode)
		s += 4
		binary.Write(buf, binary.BigEndian, e.Mode)
		s += 4
		binary.Write(buf, binary.BigEndian, e.UID)
		s += 4
		binary.Write(buf, binary.BigEndian, e.GID)
		s += 4
		binary.Write(buf, binary.BigEndian, e.Size)
		s += 4
		binary.Write(buf, binary.BigEndian, e.Sha1)
		s += 20
		binary.Write(buf, binary.BigEndian, e.Flags)
		s += 2
		binary.Write(buf, binary.BigEndian, e.PathName)
		s += len(e.PathName)

		// padding
		for i := 0; i < (8 - s%8); i++ {
			buf.WriteByte(byte(0))
		}
	}

	return buf.Bytes()
}

func write_dircache(d *Dircache, repop string) error {
	// sort by filename
	sort.Slice(d.Entries, func(i, k int) bool {
		return bytes.Compare(d.Entries[i].PathName, d.Entries[k].PathName) < 0
	})

	// set entry nunber
	d.Header.NumberOfEntries = int32(len(d.Entries))

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
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: internal error: %v\n", err)
		os.Exit(128)
	}

	// update or add or remove dircache
	for _, p := range paths {
		if do_add {
			update_dircache(d, p, true)
		} else if do_remove {
			if err := remove_dircache(d, p); err != nil {
				fmt.Fprintf(os.Stderr, "fatal: internal error: %v\n", err)
				os.Exit(128)
			}
		} else {
			update_dircache(d, p, false)
		}
	}

	// write dircache
	err = write_dircache(d, repop)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: internal error: %v\n", err)
		os.Exit(1)
	}
}

func update_dircache(d *Dircache, path string, do_add bool) {
	// already added?
	if find_dircache_entry(d, path) < 0 && do_add == false {
		fmt.Fprintf(os.Stderr, "error: %s: does not exist and --remove not passed\n", path)
		fmt.Fprintf(os.Stderr, "fatal: Unable to process path %s\n", path)
		os.Exit(128)
	}

	// file stat
	info, err := os.Lstat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		fmt.Fprintf(os.Stderr, "fatal: Unable to process path %s\n", path)
		os.Exit(128)
	}

	// create hash-object
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		fmt.Fprintf(os.Stderr, "fatal: Unable to process path %s\n", path)
		os.Exit(128)
	}
	defer f.Close()
	sha, err := hash_object(true, f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "internal error: %v\n", err)
		fmt.Fprintf(os.Stderr, "fatal: Unable to process path %s\n", path)
		os.Exit(128)
	}

	e := &DircacheEntry{}

	switch runtime.GOOS {
	case "windows":
		log.Fatal("windows not suppported")
	default:
		internal_info := info.Sys().(*syscall.Stat_t)

		e.CTimeSeconds = uint32(internal_info.Ctim.Sec)
		e.CTimeNanoSeconds = uint32(internal_info.Ctim.Nsec)

		e.MTimeSeconds = uint32(internal_info.Mtim.Sec)
		e.MTimeNanoSeconds = uint32(internal_info.Mtim.Nsec)

		e.Dev = uint32(internal_info.Dev)
		e.Inode = uint32(internal_info.Ino)

		var modeFlag uint32
		// 4-bit object type valid values in binary are 1000 (regular file), 1010 (symbolic link) and 1110 (gitlink)
		if info.Mode()&os.ModeSymlink != 0 {
			modeFlag |= uint32(0b00000000000000001010000000000000)
		} else {
			// regular file (git link is current unsupported)
			modeFlag |= uint32(0b00000000000000001000000000000000)
		}
		// 3-bit unused
		// 9-bit unix permission. Only 0755 and 0644 are valid for regular files. Symbolic links and gitlinks have value 0 in this field.
		if info.Mode()&os.ModeSymlink == 0 {
			perm := uint32(0644)
			perm |= (uint32(0111) & uint32(info.Mode()))
			modeFlag |= perm
		}
		e.Mode = modeFlag

		e.UID = internal_info.Uid
		e.GID = internal_info.Gid
		e.Size = uint32(info.Size())

		for i, v := range sha {
			e.Sha1[i] = v
		}

		var flag uint16
		flag |= uint16(0b0000000000000000)     // [1-bit: assume-valid flag]
		flag |= uint16(0b0000000000000000)     // [1-bit: extended flag(must be zero)]
		flag |= uint16(0b0000000000000000)     // [2-bit: stage(during merge)]
		mask := uint16(0b00001111111111111111) // [12-bit: name length]
		flag |= mask & uint16(len(path))
		e.Flags = flag

		e.PathName = []byte(path)
		d.Entries = append(d.Entries, e)
	}
}

func find_dircache_entry(d *Dircache, path string) int {
	pathb := []byte(path)

	for i, d := range d.Entries {
		if bytes.Equal(d.PathName, pathb) {
			return i
		}
	}
	return -1
}

func remove_dircache(d *Dircache, path string) error {
	// already added in cache?
	idx := find_dircache_entry(d, path)
	if idx < 0 {
		return fmt.Errorf("%s is not found in index\n", path)
	}

	// already deleted on filesystem?
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if err != nil && os.IsNotExist(err) == false {
		return err
	}

	d.Entries = append(d.Entries[:idx], d.Entries[idx+1:]...)

	return nil
}
