#!/bin/bash

REPOSITORY_DIR_NAME=".toy-git"

# initialize repository
cd test
../toy-git init > /dev/null

# add index to test file
../toy-git update-index --add test-target-dir/test-target-file-nested.txt

# create alias for testing by git command
ln -s ./.toy-git .git

# testing by git command
EXPECT_LS_FILES_MESSAGE="test-target-dir/test-target-file-nested.txt"
LS_FILES_MESSAGE=`git ls-files`
if [[ "$LS_FILES_MESSAGE" != "$EXPECT_LS_FILES_MESSAGE" ]]; then
  echo "[update-index] index file is broken. 'git ls-files' output was"
  echo "Expect: $EXPECT_LS_FILES_MESSAGE"
  echo "Actual: $LS_FILES_MESSAGE"
  exit 1
fi
