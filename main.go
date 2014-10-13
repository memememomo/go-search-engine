package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type SuffixText struct {
	Index  int
	String string
}
type SuffixTexts []SuffixText

func (t SuffixTexts) Len() int {
	return len(t)
}

func (t SuffixTexts) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t SuffixTexts) Less(i, j int) bool {
	return t[i].String < t[j].String
}

type Page struct {
	Path string
	Size int
}

func search_sub_string(str string, substr string) []int {
	suffix_array := create_suffix_array(str)
	p := binary_search(suffix_array, substr, 0, len(suffix_array))
	if p == -1 {
		return make([]int, 0)
	} else {
		var start, end = p, p

		for i := p - 1; i > 0; i-- {
			if strings.Index(suffix_array[i].String, substr) != 0 {
				start = i + 1
				break
			}
		}
		for i := p + 1; i < len(suffix_array); i++ {
			if strings.Index(suffix_array[i].String, substr) != 0 {
				end = i - 1
				break
			}
		}

		l := []int{}
		for i := start; i <= end; i++ {
			l = append(l, suffix_array[i].Index)
		}

		return l
	}
}

func create_suffix_array(str string) SuffixTexts {
	s := strings.Split(str, "")
	len := len(s)
	suffix := make(SuffixTexts, len)
	for i := 0; i < len; i++ {
		suffix[i] = SuffixText{Index: i, String: strings.Join(s[i:len], "")}
	}
	sort.Sort(suffix)
	return suffix
}

func binary_search(suffix_array SuffixTexts, substr string, start int, end int) int {
	p := (start + end) / 2
	str := suffix_array[p].String

	if strings.Index(str, substr) == 0 {
		return p
	} else if start == p && p == end {
		return -1
	}

	if str < substr {
		return binary_search(suffix_array, substr, p+1, end)
	} else {
		return binary_search(suffix_array, substr, 0, p)
	}
}

func read_file(path string) (string, error) {
	var fp *os.File
	var err error
	var html string

	fp, err = os.Open(path)
	if err != nil {
		return "", err
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		html += scanner.Text() + "\n"
	}

	return html, nil
}

func load_content(files []string) ([]Page, string, error) {
	var content = ""
	var info []Page

	for i := 0; i < len(files); i++ {
		file := files[i]
		html, _ := read_file(file)
		content += html + "\000"
		info = append(info, Page{Path: file, Size: len(html)})
	}

	return info, content, nil
}

func main() {
	var files []string
	var info []Page
	var pos []int
	var search_text = "search"

	if len(os.Args) < 2 {
		panic("Error!")
	}

	for i := 1; i < len(os.Args); i++ {
		files = append(files, os.Args[i])
	}

	info, content, _ := load_content(files)

	pos = search_sub_string(content, search_text)
	for i := 0; i < len(pos); i++ {
		var cur = 0
		for j := 0; j < len(info); j++ {
			if pos[i] < cur+info[j].Size {
				fmt.Printf("%s:%d\n", info[j].Path, pos[i])
				break
			}
			cur += info[j].Size
		}
	}
}
