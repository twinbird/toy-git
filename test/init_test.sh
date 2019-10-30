#!/bin/bash

REPOSITORY_DIR_NAME=".toy-git"

MESSAGE=`./toy-git init`

# repository dir should exist
if [ ! -d $REPOSITORY_DIR_NAME ]; then
  echo "[init] REPOSITORY DIRECTORY wasn't created."
  exit
fi

# initialize message should output
if [ "Initialized empty Git repository in $PWD/$REPOSITORY_DIR_NAME" != "$MESSAGE" ]; then
  echo "[init] initialize message was wrong."
  echo "Expect: Initialized empty Git repository in $PWD/$REPOSITORY_DIR_NAME"
  echo "Actual: $MESSAGE"
  exit
fi

# HEAD should exist
if [ ! -e "$REPOSITORY_DIR_NAME/HEAD" ]; then
  echo "[init] 'HEAD' wasn't created."
  exit
fi

# config should exist
if [ ! -e "$REPOSITORY_DIR_NAME/config" ]; then
  echo "[init] 'config' wasn't created."
  exit
fi

# description should exist
if [ ! -e "$REPOSITORY_DIR_NAME/description" ]; then
  echo "[init] 'description' wasn't created."
  exit
fi

# hooks dir should exist
if [ ! -d "$REPOSITORY_DIR_NAME/hooks" ]; then
  echo "[init] 'hooks' wasn't created."
  exit
fi

# info dir should exist
if [ ! -d "$REPOSITORY_DIR_NAME/info" ]; then
  echo "[init] 'info' wasn't created."
  exit
fi

# objects dir should exist
if [ ! -d "$REPOSITORY_DIR_NAME/objects" ]; then
  echo "[init] 'objects' wasn't created."
  exit
fi

# refs dir should exist
if [ ! -d "$REPOSITORY_DIR_NAME/refs" ]; then
  echo "[init] 'refs' wasn't created."
  exit
fi

# HEAD should be 'ref: refs/heads/master'
HEAD_TEXT=`cat "$REPOSITORY_DIR_NAME"/HEAD`
if [ "$HEAD_TEXT" != "ref: refs/heads/master" ]; then
  echo "[init] In 'HEAD' text was wrong."
  echo "Expect: ref: refs/heads/master"
  echo "Actual: $HEAD_TEXT"
  exit
fi
