package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(input string) []string {
	words := strings.Fields(input)
	if len(words) == 0 {
		return words
	}

	frequencies := map[string]int{}
	for _, word := range words {
		frequencies[word]++
	}

	sortedWords := make([]string, len(frequencies))
	idx := 0
	for key := range frequencies {
		sortedWords[idx] = key
		idx++
	}

	sort.Slice(sortedWords, func(i, j int) bool {
		if frequencies[sortedWords[i]] == frequencies[sortedWords[j]] {
			return sortedWords[i] < sortedWords[j]
		} else {
			return frequencies[sortedWords[i]] > frequencies[sortedWords[j]]
		}
	})

	wordsCount := len(frequencies)
	if wordsCount > 10 {
		wordsCount = 10
	}

	return sortedWords[:wordsCount:wordsCount]
}
