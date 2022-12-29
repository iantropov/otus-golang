package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

func isLetter(c rune) bool {
	return !unicode.IsLetter(c)
}

func Top10(input string) []string {
	words := strings.Fields(input)
	if len(words) == 0 {
		return words
	}

	frequencies := map[string]int{}
	for _, word := range words {
		normalizedValue := strings.ToLower(strings.TrimFunc(word, isLetter))
		if len(normalizedValue) == 0 {
			normalizedValue = word
		}
		if normalizedValue != "-" {
			frequencies[normalizedValue]++
		}
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
		}
		return frequencies[sortedWords[i]] > frequencies[sortedWords[j]]
	})

	wordsCount := len(frequencies)
	if wordsCount > 10 {
		wordsCount = 10
	}

	result := make([]string, wordsCount)
	copy(result, sortedWords)
	return result
}
