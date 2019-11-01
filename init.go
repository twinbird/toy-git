package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func init_cmd() {
	//========================================
	// make repository dir path
	//========================================
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
		os.Exit(1)
	}

	repo_path := filepath.Join(wd, REPOSITORY_DIR_NAME)

	existed := true
	if _, err := os.Stat(repo_path); os.IsNotExist(err) {
		existed = false
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
		os.Exit(1)
	}

	//========================================
	// create repository file and directory
	//========================================
	var p string

	// repository dir
	if existed == false {
		p = filepath.Join(wd, REPOSITORY_DIR_NAME)
		err := os.Mkdir(p, 0775)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
			os.Exit(1)
		}
	}

	// HEAD
	p = filepath.Join(repo_path, "HEAD")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err = ioutil.WriteFile(p, []byte("ref: refs/heads/master"), 0664)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
			os.Exit(1)
		}
	}
	// branches(dir)
	p = filepath.Join(repo_path, "branches")
	err = os.Mkdir(p, 0775)
	if err != nil && os.IsExist(err) == false {
		fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
		os.Exit(1)
	}
	// config
	p = filepath.Join(repo_path, "config")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err = ioutil.WriteFile(p, []byte(`[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true`), 0664)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
			os.Exit(1)
		}
	}
	// description
	p = filepath.Join(repo_path, "description")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err = ioutil.WriteFile(p, []byte("Unnamed repository; edit this file 'description' to name the repository."), 0664)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
			os.Exit(1)
		}
	}
	// hooks(dir)
	p = filepath.Join(repo_path, "hooks")
	err = os.Mkdir(p, 0775)
	if err != nil && os.IsExist(err) == false {
		fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
		os.Exit(1)
	}
	// info(dir)
	p = filepath.Join(repo_path, "info")
	err = os.Mkdir(p, 0775)
	if err != nil && os.IsExist(err) == false {
		fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
		os.Exit(1)
	}
	p = filepath.Join(repo_path, "exclude")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err = ioutil.WriteFile(p, []byte(`# git ls-files --others --exclude-from=.git/info/exclude
# Lines that start with '#' are comments.
# For a project mostly in C, the following would be a good set of
# exclude patterns (uncomment them if you want to use them):
# *.[oa]
# *~`), 0664)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
			os.Exit(1)
		}
	}
	// objects(dir)
	p = filepath.Join(repo_path, "objects", "info")
	err = os.MkdirAll(p, 0775)
	if err != nil && os.IsExist(err) == false {
		fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
		os.Exit(1)
	}
	p = filepath.Join(repo_path, "objects", "pack")
	err = os.MkdirAll(p, 0775)
	if err != nil && os.IsExist(err) == false {
		fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
		os.Exit(1)
	}
	// refs(dir)
	p = filepath.Join(repo_path, "refs", "heads")
	err = os.MkdirAll(p, 0775)
	if err != nil && os.IsExist(err) == false {
		fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
		os.Exit(1)
	}
	p = filepath.Join(repo_path, "refs", "tags")
	err = os.MkdirAll(p, 0775)
	if err != nil && os.IsExist(err) == false {
		fmt.Fprintf(os.Stderr, "Error: 'git init' failed. %v\n", err)
		os.Exit(1)
	}

	//========================================
	// print message
	//========================================
	if existed {
		fmt.Printf("Reinitialized existing Git repository in %s\n", repo_path)
	} else {
		fmt.Printf("Initialized empty Git repository in %s\n", repo_path)
	}
}
