# Bash wrapper to change directory to the output of gocd
gocd () {
  if dir=$($GOPATH/bin/go-cd $@); then
    cd "$dir"
  fi
} 