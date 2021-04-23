package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/renstrom/fuzzysearch/fuzzy"
)

type PkgFinder struct {
	gopath    string
	cache     map[string]time.Time
	cachePath string
}

// Find a package by the given key
func (w *PkgFinder) Find(find string) OrderedRanks {
	w.loadCache()
	defer w.storeCache()

	filepath.Walk(w.gopath, w.walker())

	if r := w.findExact(find); r.Len() > 0 {
		return r
	}

	return w.findFuzzy(find)
}

func (w *PkgFinder) loadCache() {
	w.cache = make(map[string]time.Time)
	usr, _ := user.Current()
	w.cachePath = filepath.Join(usr.HomeDir, ".gocd")

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
		if !i.IsDir() {
			return nil
		}

		// Skip if path prefixed with .,_ or vendor
		if strings.HasPrefix(i.Name(), ".") || strings.HasPrefix(i.Name(), "_") || strings.Contains(path, "vendor") {
			return nil
		}

		pkg, _ := filepath.Rel(w.gopath, filepath.Dir(path))
		ppp := strings.Split(pkg, string(filepath.Separator))
		if len(ppp) > 2 {
			return filepath.SkipDir
		}
		pkgName := filepath.Join(pkg, i.Name())
		if cached, ok := w.cache[pkgName]; ok {
			if len(ppp) > 1 && cached.Sub(i.ModTime()) <= 0 {
				log.Println("Skip dir: ", path, i.ModTime(), cached)
				return filepath.SkipDir
			}
		}
		w.cache[pkgName] = i.ModTime()
		return nil
	}
}

func (w *PkgFinder) findExact(find string) OrderedRanks {
	found := OrderedRanks{}
	for pkg := range w.cache {
		components := strings.Split(pkg, string(filepath.Separator))
		for x := len(components) - 1; x >= 0; x-- {
			if find == components[x] {
				found = append(found, fuzzy.Rank{
					Target: filepath.Join(w.gopath, pkg),
				})
			}
		}
	}
	return found
}

func (w *PkgFinder) findFuzzy(find string) OrderedRanks {
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
	var first fuzzy.Rank
	allSame := true
	for i, pkg := range found {
		if i == 0 {
			first = pkg
		} else if !strings.HasPrefix(pkg.Target, first.Target) {
			allSame = false
		}
	}
	if allSame {
		return OrderedRanks{first}
	}
	return found
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
