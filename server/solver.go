/*
 * Copyright (c) 2015, Henrik Mattsson
 * All rights reserved. See LICENSE.
 */

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
)

func NewSolver() *Solver {
	s := Solver{}
	s.words = map[string][]string{}
	return &s
}

// Map from key -> [words matching key]
type Solver struct {
	words map[string][]string
}

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
func (s *Solver) Solve(hand string) []string {
	return s.words[sortStringByChar(hand)]
}

// Preprocesses the word list
func (s *Solver) Populate(corpuspath string) error {
	s.words = map[string][]string{}

	file, err := os.Open(corpuspath)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to open corpus '%s': %s", corpuspath, err))
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		word := scanner.Text()
		key := sortStringByChar(word)
		s.words[key] = append(s.words[key], word)
	}

	if scanner.Err() != nil {
		return errors.New(fmt.Sprintf("Unable to tokenize corpus: %s", scanner.Err()))
	}

	return nil
}
