// This program rename numbered files, so that they can be alphabetically sorted.
// Ex.: before renaming:
// File10
// File1
// File2
// File3
// File4
// File5
// File6
// File7
// File8
// File9
//
// after renaming:
// File01
// File02
// File03
// File04
// File05
// File06
// File07
// File08
// File09
// File10

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"unicode"
)

// splits receives a file name and return it broken into pieces.
func split(s string) (prefix, suffix, digits string) {
	begin := -1
	end := -1
	for i, r := range s {
		if unicode.IsDigit(r) {
			if begin < 0 {
				begin = i
			}
			end = i
		} else {
			if begin >= 0 {
				break
			}
		}
	}
	if begin < 0 {
		return "", "", ""
	}
	return s[:begin-1], s[end+1:], s[begin : end+1]
}

// extractDigits obtains the maximum lenght of the numeric part
// receives a slice of FileInfo and returns the calculated map.
func extractDigits(files []os.FileInfo) map[string]int {
	m := make(map[string]int)
	for _, fi := range files {
		prefix, suffix, digits := split(fi.Name())
		mask := prefix + "|" + suffix
		if len(digits) > m[mask] {
			m[mask] = len(digits)
		}
	}
	return m
}

// createMapping produces the list of files to be renamed
// each entry in the map contains the old and the new name.
func createMapping(files []os.FileInfo) map[string]string {
	result := make(map[string]string)

	m := extractDigits(files)
	for _, fi := range files {
		prefix, suffix, digits := split(fi.Name())
		if len(digits) > 0 {
			mask := prefix + "|" + suffix
			max := m[mask]
			for len(digits) < max {
				digits = "0" + digits
			}
			newName = prefix + digits + suffix
			if fi.Name() != newName {
				result[fi.Name()] = newName
			}
		}
	}

	return result
}

// show exibits the files to be renamed
func show(mapping map[string]string) {
	for key, value := range mapping {
		fmt.Println(key, " -> ", value)
	}
}

// rename actually change the file names accordingly to the map.
func rename(mapping map[string]string) {
	for key, value := range mapping {
		err := os.Rename(key, value)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}
	mapping := createMapping(files)
	if len(mapping) == 0 {
		fmt.Fprintf(os.Stderr, "No files found")
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "/r" {
		rename(mapping)
	} else {
		show(mapping)
		fmt.Println("\nPlease use /r to rename")
	}
}
