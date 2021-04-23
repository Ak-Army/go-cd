package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const VendorToken = "^"

func getGoPath() (string, error) {
	if path := os.Getenv("GOPATH"); path != "" {
		return filepath.Join(path, "src"), nil
	}
	path := filepath.Join(build.Default.GOPATH, "src")
	_, err := os.Stat(path)
	return path, err
}

func tryGoToVendorParent() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if !strings.Contains(cwd, "vendor") {
		return
	}

	components := strings.Split(cwd, string(filepath.Separator))
	for i := len(components) - 1; i >= 0; i-- {
		if components[i] == "vendor" {
			if i == 0 {
				// "vendor" is at the root of the path
				return
			}

			var abs = append([]string{"/"}, components[:i]...)
			fmt.Print(filepath.Join(abs...))
		}
	}
}

func main() {
	verbose := flag.Bool("v", false, "URL of ratesheet server")
	flag.Parse()

	log.SetFlags(log.Llongfile)
	if !*verbose {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	log.Printf("Version: %s build time: %s\n", Version, BuildTime)

	path, err := getGoPath()
	if err != nil {
		log.Fatal(err)
	}

	// If no path supplied then change directory to GOPATH.
	if flag.NArg() == 0 {
		fmt.Print(path)
		return
	}

	if flag.Arg(0) == VendorToken {
		tryGoToVendorParent()
		return
	}
	find(path)
}

func find(path string) {
	w := PkgFinder{
		gopath: path,
	}

	matches := w.Find(flag.Arg(0))
	if len(matches) == 0 {
		fmt.Println("No matching package found")
	}
	if len(matches) == 1 {
		fmt.Println(matches[0].Target)
		return
	}

	if flag.NArg() > 1 {
		i, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			log.Fatalf("cannot parse requested index %s: %s", flag.Arg(1), err)
		}

		if i > len(matches) {
			log.Fatalf("%d is an invalid index (max %d)", i, len(matches))
		}

		fmt.Println(matches[i].Target)
		return
	}

	for i, m := range matches {
		rel, _ := filepath.Rel(path, m.Target)
		fmt.Printf("  %d %s\n", i, rel)
	}
}
