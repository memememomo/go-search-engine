package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func search_sub_string(str string, substr string) []int {
	suffix_array := create_suffix_array(str)
	p := binary_search(suffix_array, substr, 0, len(suffix_array))
	if p == -1 {
		return make([]int, 0)
	} else {
		var start, end = p, p
		for i := p - 1; i > 0; i-- {
			if strings.Index(suffix_array[i], substr) != 0 {
				start = i + 1
				break
			}
		}
		for i := p + 1; i < len(suffix_array); i++ {
			if strings.Index(suffix_array[i], substr) != 0 {
				end = i - 1
				break
			}
		}
		l := []int{}
		if start != end {
			for i := start; i <= end; i++ {
				l = append(l, i)
			}
		} else {
			l[0] = p
		}
		return l
	}
}

func create_suffix_array(str string) []string {
	s := strings.Split(str, "")
	len := len(s)
	suffix := make([]string, len)
	for i := 0; i < len; i++ {
		suffix[i] = strings.Join(s[i:len], "")
	}
	sort.Strings(suffix)
	return suffix
}

func binary_search(suffix_array []string, substr string, start int, end int) int {
	p := (start + end) / 2
	str := suffix_array[p]

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

func main() {
	var fp *os.File
	var err error
	var html string

	if len(os.Args) < 2 {
		fp = os.Stdin
	} else {
		fp, err = os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		defer fp.Close()
	}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		html += scanner.Text() + "\n"
	}

	fmt.Println(search_sub_string(html, "a"))
}
