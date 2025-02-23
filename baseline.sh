#!/bin/bash

# iterate over a list of directories and descend into each one
# and pull the latest changes from the remote repositor

# REF:
# - https://stackoverflow.com/questions/2107945/how-to-loop-over-directories-in-linux

for dir in ./*/; do
  dir=${dir%*/} # remove trailing slash
  echo "Directory: $dir"
  cd "$dir" || exit 1
  echo "Pulling latest changes from remote repository"
  git pull
  cd ..
done
