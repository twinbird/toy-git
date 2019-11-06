#!/bin/bash

REPOSITORY_DIR_NAME=".toy-git"

# initialize repository
./toy-git init > /dev/null

# create sha1 value from test-file
HASH_OBJECT_MESSAGE=`./toy-git hash-object test/test-target-file.txt`
TEST_TARGET_FILE_SHA1="4239a627fe3921e9beb24954158a9c47fb5683ec"

if [[ $HASH_OBJECT_MESSAGE != $TEST_TARGET_FILE_SHA1 ]]; then
  echo "[hash-object] created sha1 hash value was wrong."
  echo "Expect: $TEST_TARGET_FILE_SHA1"
  echo "Actual: $HASH_OBJECT_MESSAGE"
  exit 1
fi

# create sha1 value from STDIN
HASH_OBJECT_MESSAGE=`cat test/test-target-file.txt | ./toy-git hash-object --stdin`
TEST_TARGET_FILE_SHA1="4239a627fe3921e9beb24954158a9c47fb5683ec"

if [[ $HASH_OBJECT_MESSAGE != $TEST_TARGET_FILE_SHA1 ]]; then
  echo "[hash-object] created sha1 hash value(from STDIN) was wrong."
  echo "Expect: $TEST_TARGET_FILE_SHA1"
  echo "Actual: $HASH_OBJECT_MESSAGE"
  exit 1
fi

# store test-file to git repository
./toy-git hash-object -w test/test-target-file.txt

if [[ ! -e "$REPOSITORY_DIR_NAME/objects/42/39a627fe3921e9beb24954158a9c47fb5683ec" ]]; then
  echo "[hash-object] 'hash-object -w test/test-target-file.txt' does not stored git repository."
  exit 1
fi

# print file type
CAT_FILE_TYPE=$( ./toy-git cat-file -t $TEST_TARGET_FILE_SHA1 )

if [[ $CAT_FILE_TYPE != "blob" ]]; then
  echo "[hash-object, cat-file] 'cat-file -t $TEST_TARGET_FILE_SHA1' do not displayed 'blob'"
  echo "Expect: blob"
  echo "Actual: $CAT_FILE_TYPE"
  exit 1
fi

# print file size
CAT_FILE_SIZE=$( ./toy-git cat-file -s $TEST_TARGET_FILE_SHA1 )

if [[ $CAT_FILE_SIZE -ne 33 ]]; then
  echo "[hash-object, cat-file] 'cat-file -s $TEST_TARGET_FILE_SHA1' unmatched file size"
  echo "Expect: 33"
  echo "Actual: $CAT_FILE_SIZE"
  exit 1
fi

# print blob data
CAT_FILE_DATA=$( ./toy-git cat-file -p $TEST_TARGET_FILE_SHA1 )
EXPECT_FILE_DATA=$( cat test/test-target-file.txt )

if [[ $CAT_FILE_DATA != $EXPECT_FILE_DATA ]]; then
  echo "[hash-object, cat-file] 'cat-file -p $TEST_TARGET_FILE_SHA1' unmatched blob"
  echo "Expect: $EXPECT_FILE_DATA"
  echo "Actual: $CAT_FILE_DATA"
  exit 1
fi
