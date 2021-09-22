package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(sample string) []string {
	wordsCount := map[string]int{}

	for _, word := range strings.Fields(sample) {
		wordsCount[word]++
	}

	words := getKeys(wordsCount)
	sort.Slice(words, func(i, j int) bool {
		if wordsCount[words[i]] == wordsCount[words[j]] {
			return words[i] < words[j]
		}
		return wordsCount[words[i]] > wordsCount[words[j]]
	})

	if len(words) < 10 {
		return words
	}

	return words[:10]
}

func getKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
