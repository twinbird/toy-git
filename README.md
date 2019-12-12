# toy-git

Git Plumbing commands implementation for learning git internal.

You **MUST NOT** use your production repository.

## Build

toy-git is written by go.

Get it [here](https://golang.org/).

```
$ git clone https://github.com/twinbird/toy-git.git
$ cd toy-git
$ make
```

## Test

```
$ make test
```

## Commands

toy-git is implement following commands subset.

> Ex) toy-git update-index --add README.md

 * git init
 * git hash-object
 * git cat-file
 * git update-index
 * git ls-file
 * git write-tree
 * git commit-tree
 * git update-ref

## Thanks & Reference

* [Pro Git](https://git-scm.com/book/en/v2)
* [GitHub - git/git](https://github.com/git/git/blob/master/Documentation/technical/index-format.txt)
* [Git User Manual](https://mirrors.edge.kernel.org/pub/software/scm/git/docs/user-manual.html#the-object-database)
* [Mercari Engineering Blog - Gitのステージング領域の正体を探る](https://tech.mercari.com/entry/2017/04/06/171430)

