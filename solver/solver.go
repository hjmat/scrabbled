/*
 * Copyright (c) 2015, Henrik Mattsson
 * All rights reserved. See LICENSE.
 */

package solver

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

// Map from key -> [words matching key]
var words = map[string][]string{}

// Alphabetic string sorting implementation
type sortChars []rune

func (s sortChars) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortChars) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortChars) Len() int {
	return len(s)
}

// Sorts a string by character, e.g. cat -> act
func sortStringByChar(s string) string {
	sChars := []rune(s)
	sort.Sort(sortChars(sChars))
	return string(sChars)
}

// Finds dictionary words that match a scrabble hand
func Solve(hand string) []string {
	return words[sortStringByChar(hand)]
}

// Preprocesses the word list
func Populate(corpuspath string) error {
	file, err := os.Open(corpuspath)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to open corpus '%s': %s", corpuspath, err))
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		key := sortStringByChar(word)
		words[key] = append(words[key], word)
	}

	if scanner.Err() != nil {
		return errors.New(fmt.Sprintf("Unable to tokenize corpus: %s", scanner.Err()))
	}

	return nil
}
