package util

import (
	"fmt"
	"regexp"
)

func CompareTwoLinks(linkOne, linkTwo string) float32 {

	// if identical, return 1 instantly
	if linkOne == linkTwo {
		return 1
	}

	firstBigrams := make(map[string]int)
	for i := 0; i < len(linkOne)-1; i++ {
		a := fmt.Sprintf("%c", linkOne[i])
		b := fmt.Sprintf("%c", linkOne[i+1])

		bigram := a + b

		var count int

		if value, ok := firstBigrams[bigram]; ok {
			count = value + 1
		} else {
			count = 1
		}

		firstBigrams[bigram] = count
	}

	var intersectionSize float32
	intersectionSize = 0

	for i := 0; i < len(linkTwo)-1; i++ {
		a := fmt.Sprintf("%c", linkTwo[i])
		b := fmt.Sprintf("%c", linkTwo[i+1])

		bigram := a + b

		var count int

		if value, ok := firstBigrams[bigram]; ok {
			count = value
		} else {
			count = 0
		}

		if count > 0 {
			firstBigrams[bigram] = count - 1
			intersectionSize = intersectionSize + 1
		}
	}

	return (2.0 * intersectionSize) / (float32(len(linkOne)) + float32(len(linkTwo)) - 2)
}

func ExtractLinks(s string) []string {

	links := []string{}
	// https?:\/\/((?:[\w,-,_,~]+\.)*[\w,-,_]+\.\w+)
	r := regexp.MustCompile("https?:\\/\\/((?:[\\w,-,_,~]+\\.)*[\\w,-,_]+\\.\\w+)")

	matches := r.FindAllStringSubmatch(s, -1)
	for _, m := range matches {
		links = append(links, m[1])
	}

	return links

}
