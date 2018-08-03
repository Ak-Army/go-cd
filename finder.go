package main

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"encoding/json"
	"fmt"
	"github.com/renstrom/fuzzysearch/fuzzy"
	"io/ioutil"
	"os/user"
	"time"
)

type PkgFinder struct {
	gopath    string
	cache     map[string]time.Time
	cachePath string
}

// Find a package by the given key
func (w *PkgFinder) Find(find string) *OrderedRanks {
	w.loadCache()
	defer w.storeCache()

	filepath.Walk(w.gopath, w.walker())

	if r := w.findExact(find); r!= nil {
		return r
	}

	return w.findFuzzy(find)
}

func (w *PkgFinder) loadCache() {
	w.cache = make(map[string]time.Time)
	usr, _ := user.Current()
	w.cachePath = fmt.Sprintf("%s/.gocd", usr.HomeDir)

	cache, _ := ioutil.ReadFile(w.cachePath)
	json.Unmarshal(cache, &w.cache)
}

func (w *PkgFinder) storeCache() {
	b, _ := json.Marshal(w.cache)
	ioutil.WriteFile(w.cachePath, b, 0644)
}

func (w *PkgFinder) walker() filepath.WalkFunc {
	return func(path string, i os.FileInfo, err error) (e error) {
		// Skip GOPATH/src
		if path == w.gopath {
			return nil
		}
		// Skip if path contains .,_ or vendor
		if i.IsDir() && (strings.HasPrefix(i.Name(), ".") || strings.HasPrefix(i.Name(), "_") || strings.Contains(path, "vendor")) {
			return filepath.SkipDir
		}
		// Ignore if path is a directory or is not a go file.
		if i.IsDir() || !strings.HasSuffix(i.Name(), "go") {
			return nil
		}

		// Scan every component of the relative path until we find a direct match.
		pkg, _ := filepath.Rel(w.gopath, filepath.Dir(path))
		ppp := strings.Split(pkg, string(filepath.Separator))
		if len(ppp) > 3 {
			return filepath.SkipDir
		}
		if cached, ok := w.cache[pkg]; ok {
			if cached.Sub(i.ModTime()) < 0 {
				return filepath.SkipDir
			}
		}
		w.cache[pkg] = i.ModTime()
		return nil
	}
}

func (w *PkgFinder) findExact(find string) *OrderedRanks{
	for pkg := range w.cache {
		components := strings.Split(pkg, string(filepath.Separator))
		for x := len(components) - 1; x >= 0; x-- {
			if find == filepath.Join(components[x:]...) {
				return &OrderedRanks{
					{
						Target: filepath.Join(w.gopath, pkg),
					},
				}
			}
		}
	}
	return nil
}

func (w *PkgFinder) findFuzzy(find string) *OrderedRanks {
	// Find possible matches from list of seen packages
	found := OrderedRanks{}
	for pkg := range w.cache {
		path := strings.Split(pkg, string(filepath.Separator))
		ranks := fuzzy.RankFindFold(find, path)

		for _, r := range ranks {
			if r.Distance > 10 {
				continue
			}
			r.Target = filepath.Join(w.gopath, pkg)
			found = append(found, r)
		}
	}

	sort.Sort(found)

	return &found
}

type OrderedRanks []fuzzy.Rank

func (r OrderedRanks) Len() int {
	return len(r)
}

func (r OrderedRanks) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r OrderedRanks) Less(i, j int) bool {
	if r[i].Distance < r[j].Distance {
		return true
	}
	if r[i].Distance > r[j].Distance {
		return false
	}
	return r[i].Target < r[j].Target
}