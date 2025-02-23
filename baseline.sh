#!/bin/bash

# iterate over a list of directories and descend into each one
# and pull the latest changes from the remote repositor
for dir in $(ls -d */); do
  echo "Directory: $dir"
  cd "$dir" || exit 1
  echo "Pulling latest changes from remote repository"
  git pull
  cd ..
done
