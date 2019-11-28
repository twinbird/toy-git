#!/bin/bash

REPOSITORY_DIR_NAME=".toy-git"

cd test
rm -rf .toy-git
unlink .git > /dev/null 2>&1
rm -rf .git

########################
# initialize repository
########################
../toy-git init > /dev/null
git init > /dev/null

################
# create index
################
../toy-git update-index --add test-target-file.txt
git update-index --add test-target-file.txt

################
# write-tree
################

# test generated SHA1 value
EXPECT_SHA1=`git write-tree`
ACTUAL_SHA1=`../toy-git write-tree`

# test generated tree object
EXPECT_PREFIX=${EXPECT_SHA1:0:2}
EXPECT_SUFFIX=${EXPECT_SHA1:2}
ACTUAL_PREFIX=${ACTUAL_SHA1:0:2}
ACTUAL_SUFFIX=${ACTUAL_SHA1:2}

EXPECT_TREE_PATH="$REPOSITORY_DIR_NAME/objects/${EXPECT_PREFIX}/${EXPECT_SUFFIX}"
ACTUAL_TREE_PATH="$REPOSITORY_DIR_NAME/objects/${ACTUAL_PREFIX}/${ACTUAL_SUFFIX}"

################
# commit-tree
################

EXPECT_COMMIT_HASH=$( echo "first commit" | git commit-tree $EXPECT_SHA1 )
ACTUAL_COMMIT_HASH=$( echo "first commit" | ../toy-git commit-tree $ACTUAL_SHA1 )

if [[ $EXPECT_COMMIT_HASH != $ACTUAL_COMMIT_HASH ]]; then
  echo "[commit-tree] generated commit hash is wrong."
  echo -e "Expect: \n$EXPECT_COMMIT_HASH"
  echo -e "Actual: \n$ACTUAL_COMMIT_HASH"
  exit 1
fi

EXPECT_COMMIT_PREFIX=${EXPECT_COMMIT_HASH:0:2}
EXPECT_COMMIT_SUFFIX=${EXPECT_COMMIT_HASH:2}
ACTUAL_COMMIT_PREFIX=${ACTUAL_COMMIT_HASH:0:2}
ACTUAL_COMMIT_SUFFIX=${ACTUAL_COMMIT_HASH:2}

EXPECT_COMMIT_PATH="$REPOSITORY_DIR_NAME/objects/${EXPECT_COMMIT_PREFIX}/${EXPECT_COMMIT_SUFFIX}"
ACTUAL_COMMIT_PATH="$REPOSITORY_DIR_NAME/objects/${ACTUAL_COMMIT_PREFIX}/${ACTUAL_COMMIT_SUFFIX}"

EXPECT_COMMIT_OBJECT=$( python3 zlib-decode.py $EXPECT_COMMIT_PATH )
ACTUAL_COMMIT_OBJECT=$( python3 zlib-decode.py $ACTUAL_COMMIT_PATH )

cmp <( python3 zlib-decode.py $EXPECT_COMMIT_PATH ) <( python3 zlib-decode.py $ACTUAL_COMMIT_PATH )
if [[ "$?" -ne 0 ]]; then
  echo "[commit-tree] generated commit object is wrong."
  echo -e "Commit hash: \n$ACTUAL_COMMIT_HASH"
  echo -e "Expect: \n$EXPECT_COMMIT_OBJECT"
  echo -e "Actual: \n$ACTUAL_COMMIT_OBJECT"
  exit 1
fi
