package main

import (
	"fmt"
	"go/build"
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
	log.SetFlags(0)

	path, err := getGoPath()
	if err != nil {
		log.Fatal(err)
	}

	// If no path supplied then change directory to GOPATH.
	if len(os.Args) == 1 {
		fmt.Print(path)
		return
	}

	if os.Args[1] == "-h" {
		fmt.Printf("go-cd version: %s build time: %s\n", Version, BuildTime)
		return
	}

	if os.Args[1] == VendorToken {
		tryGoToVendorParent()
		return
	}

	w := PkgFinder{
		gopath: path,
	}

	matches := w.Find(os.Args[1])

	if matches == nil {
		log.Fatal("No matching package found")
	}
	m := *matches
	if len(m) == 1 {
		fmt.Println(m[0].Target)
		return
	}

	if len(os.Args) > 2 {
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("cannot parse requested index %s: %s", os.Args[2], err)
		}

		if i > len(m) {
			log.Fatalf("%d is an invalid index (max %d)", i, len(m))
		}

		fmt.Println(m[i].Target)
		return
	}

	for i, m := range m {
		rel, _ := filepath.Rel(path, m.Target)
		log.Printf("  %d %s\n", i, rel)
	}
	os.Exit(1)
}
