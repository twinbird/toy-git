#!/bin/bash

REPOSITORY_DIR_NAME=".toy-git"

cd test
rm -rf .toy-git
unlink .git > /dev/null
rm -rf .git

# initialize repository
../toy-git init > /dev/null
git init > /dev/null

# create index
../toy-git update-index --add test-target-file.txt
#../toy-git update-index --add test-target-dir/test-target-file-nested.txt

git update-index --add test-target-file.txt
#git update-index --add test-target-dir/test-target-file-nested.txt

# write-tree

# test generated SHA1 value
EXPECT_SHA1=`git write-tree`
ACTUAL_SHA1=`../toy-git write-tree`

if [[ "$EXPECT_SHA1" != "$ACTUAL_SHA1" ]]; then
  echo "[write-tree] generated SHA1 value is wrong."
  echo -e "Expect: \n$EXPECT_SHA1"
  echo -e "Actual: \n$ACTUAL_SHA1"
  exit 1
fi

# test generated tree object
EXPECT_PREFIX=${EXPECT_SHA1:0:2}
EXPECT_SUFFIX=${EXPECT_SHA1:2}
ACTUAL_PREFIX=${ACTUAL_SHA1:0:2}
ACTUAL_SUFFIX=${ACTUAL_SHA1:2}

EXPECT_TREE="$( cat -v $REPOSITORY_DIR_NAME/objects/${EXPECT_PREFIX}/${EXPECT_SUFFIX} )"
ACTUAL_TREE="$( cat -v $REPOSITORY_DIR_NAME/objects/${ACTUAL_PREFIX}/${ACTUAL_SUFFIX} )"

if [[ "$EXPECT_TREE" != "$ACTUAL_TREE" ]]; then
  echo "[write-tree] generated tree object is wrong."
  echo -e "Expect: \n$EXPECT_TREE"
  echo -e "Actual: \n$ACTUAL_TREE"
  exit 1
fi

