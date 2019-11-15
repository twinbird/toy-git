#!/bin/bash

REPOSITORY_DIR_NAME=".toy-git"

# initialize repository
cd test
../toy-git init > /dev/null

# add test file to index
../toy-git update-index --add test-target-dir/test-target-file-nested.txt

# create alias for testing by git command
ln -s ./.toy-git .git

# testing by git command
EXPECT_LS_FILES_MESSAGE="test-target-dir/test-target-file-nested.txt"
EXPECT_GIT_LS_FILES_MESSAGE=`git ls-files`
if [[ "$EXPECT_LS_FILES_MESSAGE" != "$EXPECT_GIT_LS_FILES_MESSAGE" ]]; then
  echo "[update-index] 'update-index --add' failed."
  echo -e "Expect: \n$EXPECT_LS_FILES_MESSAGE"
  echo -e "Actual: \n$EXPECT_GIT_LS_FILES_MESSAGE"
  exit 1
fi

# testing by toy-git ls-files
LS_FILES_MESSAGE=`../toy-git ls-files`
if [[ "$LS_FILES_MESSAGE" != "$EXPECT_LS_FILES_MESSAGE" ]]; then
  echo "[update-index] 'ls-files' failed after update-index --add."
  echo -e "Expect: \n$EXPECT_LS_FILES_MESSAGE"
  echo -e "Actual: \n$LS_FILES_MESSAGE"
  exit 1
fi

# add more test file to index
../toy-git update-index --add test-target-file.txt

# testing by git command
EXPECT_LS_FILES_MESSAGE="test-target-dir/test-target-file-nested.txt
test-target-file.txt"
EXPECT_GIT_LS_FILES_MESSAGE=`git ls-files`
if [[ "$EXPECT_LS_FILES_MESSAGE" != "$EXPECT_GIT_LS_FILES_MESSAGE" ]]; then
  echo "[update-index] 'update-index --add' second file add failed."
  echo -e "Expect: \n$EXPECT_LS_FILES_MESSAGE"
  echo -e "Actual: \n$EXPECT_GIT_LS_FILES_MESSAGE"
  exit 1
fi

# testing by toy-git ls-files
LS_FILES_MESSAGE=`../toy-git ls-files`
if [[ "$LS_FILES_MESSAGE" != "$EXPECT_LS_FILES_MESSAGE" ]]; then
  echo "[update-index] 'ls-files' after second file added failed."
  echo -e "Expect: \n$EXPECT_LS_FILES_MESSAGE"
  echo -e "Actual: \n$LS_FILES_MESSAGE"
  exit 1
fi

# add and remove test file from index
touch a.txt
../toy-git update-index --add a.txt
../toy-git update-index --remove a.txt

# testing by git command
EXPECT_LS_FILES_MESSAGE="test-target-dir/test-target-file-nested.txt
test-target-file.txt
a.txt"
EXPECT_GIT_LS_FILES_MESSAGE=`git ls-files`
if [[ "$EXPECT_GIT_LS_FILES_MESSAGE" != "$EXPECT_LS_FILES_MESSAGE" ]]; then
  echo "[update-index] 'update-index --remove' without rm file is failed."
  echo -e "Expect: \n$EXPECT_LS_FILES_MESSAGE"
  echo -e "Actual: \n$LS_FILES_MESSAGE"
  exit 1
fi

# testing by toy-git command
LS_FILES_MESSAGE=`../toy-git ls-files`
if [[ "$LS_FILES_MESSAGE" != "$EXPECT_LS_FILES_MESSAGE" ]]; then
  echo "[update-index] 'update-index --remove' without rm file is failed."
  echo -e "Expect: \n$EXPECT_LS_FILES_MESSAGE"
  echo -e "Actual: \n$LS_FILES_MESSAGE"
  exit 1
fi

# remove test file from index and filesystem
rm a.txt
../toy-git update-index --remove a.txt

# testing by git command
EXPECT_LS_FILES_MESSAGE="test-target-dir/test-target-file-nested.txt
test-target-file.txt"
EXPECT_GIT_LS_FILES_MESSAGE=`git ls-files`
if [[ "$EXPECT_GIT_LS_FILES_MESSAGE" != "$EXPECT_LS_FILES_MESSAGE" ]]; then
  echo "[update-index] 'update-index --remove' is failed."
  echo -e "Expect: \n$EXPECT_LS_FILES_MESSAGE"
  echo -e "Actual: \n$LS_FILES_MESSAGE"
  exit 1
fi

# testing by toy-git command
LS_FILES_MESSAGE=`../toy-git ls-files`
if [[ "$LS_FILES_MESSAGE" != "$EXPECT_LS_FILES_MESSAGE" ]]; then
  echo "[update-index] 'update-index --remove' is failed."
  echo -e "Expect: \n$EXPECT_LS_FILES_MESSAGE"
  echo -e "Actual: \n$LS_FILES_MESSAGE"
  exit 1
fi

# remove test file and ls-files -d
touch b.txt
../toy-git update-index --add b.txt
rm b.txt

EXPECT_LS_FILES_MESSAGE="b.txt"
LS_FILES_MESSAGE=`../toy-git ls-files -d`
if [[ "$EXPECT_LS_FILES_MESSAGE" != "$LS_FILES_MESSAGE" ]]; then
  echo "[ls-files] 'ls-files -d' is failed."
  echo -e "Expect: \n$EXPECT_LS_FILES_MESSAGE"
  echo -e "Actual: \n$LS_FILES_MESSAGE"
  exit 1
fi

EXPECT_GIT_LS_FILES_MESSAGE=`git ls-files -d`
if [[ "$LS_FILES_MESSAGE" != "$EXPECT_GIT_LS_FILES_MESSAGE" ]]; then
  echo "[ls-files] 'ls-files -d' is failed."
  echo -e "Expect: \n$EXPECT_LS_FILES_MESSAGE"
  echo -e "Actual: \n$LS_FILES_MESSAGE"
  exit 1
fi

# modify test file and ls-files -m
touch c.txt
../toy-git update-index --add c.txt
echo "CHANGED" >> c.txt

EXPECT_LS_FILES_MESSAGE="b.txt
c.txt"
LS_FILES_MESSAGE=`../toy-git ls-files -m`
if [[ "$EXPECT_LS_FILES_MESSAGE" != "$LS_FILES_MESSAGE" ]]; then
  echo "[ls-files] 'ls-files -m' is failed"
  echo -e "Expect: \n$EXPECT_LS_FILES_MESSAGE"
  echo -e "Actual: \n$LS_FILES_MESSAGE"
  exit 1
fi

EXPECT_GIT_LS_FILES_MESSAGE=`git ls-files -m`
if [[ "$LS_FILES_MESSAGE" != "$EXPECT_GIT_LS_FILES_MESSAGE" ]]; then
  echo "[ls-files] 'ls-files -m' is failed"
  echo -e "Expect: \n$EXPECT_GIT_LS_FILES_MESSAGE"
  echo -e "Actual: \n$LS_FILES_MESSAGE"
  exit 1
fi

cd - > /dev/null
