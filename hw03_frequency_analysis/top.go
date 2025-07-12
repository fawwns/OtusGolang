package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str string) []string {
	type Pair struct {
		Key   string
		Value int
	}

	freq := []Pair{}
	mapStr := make(map[string]int)
	strSlic := strings.Fields(str)
	slicStr := []string{}

	for _, word := range strSlic {
		mapStr[word]++
	}

	for k, v := range mapStr {
		freq = append(freq, Pair{k, v})
	}

	sort.Slice(freq, func(i, j int) bool {
		if freq[i].Value != freq[j].Value {
			return freq[i].Value > freq[j].Value
		}
		return freq[i].Key < freq[j].Key
	})

	for i := 0; i < len(freq); i++ {
		slicStr = append(slicStr, freq[i].Key)
		if i == 9 {
			return slicStr
		}
	}
	return slicStr
}
