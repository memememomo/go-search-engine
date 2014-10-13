package main

import (
	"bufio"
	"fmt"
	"github.com/PuerkitoBio/goquery"
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
	Url  string
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

func load_content(docs map[string]*goquery.Document) ([]Page, string, error) {
	var content = ""
	var info []Page

	for url, doc := range docs {
		str, _ := doc.Html()
		content += str + "\000"
		info = append(info, Page{Url: url, Size: len(str)})
	}

	return info, content, nil
}

func get_url(doc *goquery.Document) []string {
	var urls []string
	uniq := make(map[string]bool)

	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		if url != "" {
			if strings.Index(url, "http") != 0 && strings.Index(url, "#") == -1 {
				uniq[url] = true
			}
		}
	})

	for key, _ := range uniq {
		urls = append(urls, key)
	}

	return urls
}

func scraping(base_url string, path string, cache map[string]*goquery.Document) map[string]*goquery.Document {
	url := base_url + "/" + path
	doc, _ := goquery.NewDocument(url)
	urls := get_url(doc)
	cache[url] = doc

	for i := 0; i < len(urls); i++ {
		if cache[base_url+"/"+urls[i]] == nil {
			cache = scraping(base_url, urls[i], cache)
		}
	}

	return cache
}

func main() {
	var info []Page
	var pos []int
	var cache = make(map[string]*goquery.Document)
	var search_text = "Java"
	base_url := "http://www.osss.cs.tsukuba.ac.jp/kato/codeconv"
	start_path := "CodeConvTOC.doc.html"

	cache = scraping(base_url, start_path, cache)

	info, content, _ := load_content(cache)

	pos = search_sub_string(content, search_text)
	for i := 0; i < len(pos); i++ {
		var cur = 0
		for j := 0; j < len(info); j++ {
			if pos[i] < cur+info[j].Size {
				fmt.Printf("%s:%d\n", info[j].Url, pos[i])
				break
			}
			cur += info[j].Size
		}
	}
}
