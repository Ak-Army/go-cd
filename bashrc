# Bash wrapper to change directory to the output of gocd
gocd () {
  dir=$($GOPATH/bin/go-cd $@)
  if [ -d $dir ]; then
    cd "$dir"
  else
    echo -e "$dir"
  fi
} 