// Just golang excercise for
// rg tags= content/posts/ | sed  's/.*"\([^"]*\)".*/\1/' | sort |uniq -c | sort -r
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"

func FindAllInText(re *regexp.Regexp, s string) []string {
	return re.FindAllString(s, -1)
}

func CreateProcessor(lineProc TagExtractor) filepath.WalkFunc {

	proc := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".md") {

			file, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				t := scanner.Text()
				if lineProc.Extract(t) {
					// only parse first tags= line
					break
				}
			}
		}
		return nil
	}
	return proc
}

func GetPath() string {
	path := flag.String("p", ".", "directory for starting a recursive search for .md files")
	flag.Parse()
	absp := *path
	absp, _ = filepath.Abs(*path)
	fmt.Println(Green + "Searching path: " + Yellow + absp + Reset)
	return *path
}

type TagExtractor interface {
	// Extracts tags from s and returns true
	// if there was something to extract in s
	Extract(s string) bool
}

type Tag struct {
	name  string
	count int
}

// Keeps all tags
type Tags struct {
	nameToTag map[string]*Tag
}

func NewTags() *Tags {
	return &Tags{nameToTag: map[string]*Tag{}}
}

// Extracts all tags from string s
func (t *Tags) Extract(s string) bool {
	tagsGroupRe := regexp.MustCompile(`tags\s*=`)
	tagsExtractRe := regexp.MustCompile(`"[^"]*"`)
	if !tagsGroupRe.MatchString(s) {
		return false
	}
	ms := tagsExtractRe.FindAllString(s, -1)
	if ms != nil && len(ms) > 0 {
		for _, m := range ms {
			t.Update(strings.Trim(m, "\""))
		}
	} else {
		return false
	}
	return true
}

func (t *Tag) Print() {
	indent := len(Red + "1234" + Reset)
	fmt.Printf("%*s entries with tag  %s\n", indent, Red+strconv.Itoa(t.count)+Reset, Green+t.name+Reset)
}

func (ts *Tags) Update(v string) {
	el, exists := ts.nameToTag[v]
	if !exists {
		el = &Tag{v, 0}
	}
	el.count += 1
	ts.nameToTag[v] = el
}

func (ts Tags) Len() int {
	return len(ts.nameToTag)
}

type TagA []Tag

func (ta TagA) Len() int {
	return len(ta)
}

func (ta TagA) Less(i, j int) bool {
	a := ta[i]
	b := ta[j]

	if a.count > b.count {
		return true
	}
	if a.count == b.count && strings.ToLower(a.name) < strings.ToLower(b.name) {
		return true
	}
	return false
}
func (ta TagA) Swap(i, j int) {
	ta[i], ta[j] = ta[j], ta[i]
}

// Sorts by count (highest first) then by lowercase name
func (ts Tags) Sorted() TagA {
	tags := make(TagA, 0)
	for _, tg := range ts.nameToTag {
		tags = append(tags, *tg)
	}
	sort.Sort(tags)
	return tags
}

func (ts Tags) Print() {
	sorted := ts.Sorted()
	for _, t := range sorted {
		indent := len(Red + "1234" + Reset)
		fmt.Printf("%*s:  %s\n", indent, Red+strconv.Itoa(t.count)+Reset, Green+t.name+Reset)
	}
}

func main() {
	tags := NewTags()
	err := filepath.Walk(GetPath(), CreateProcessor(tags))
	if err != nil {
		log.Printf("Walk exited wiht err %v", err)
	}
	tags.Print()

}
