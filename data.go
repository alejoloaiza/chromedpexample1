package main

import (
	"regexp"
	"strings"
)

func SplitPrices(input string) (string, string) {
	priceextractor := regexp.MustCompile(`\$[0-9]{1,3}(?:\,[0-9]{1,3})?`)
	match := priceextractor.FindStringSubmatch(input)
	if len(match) > 0 {
		time := input[strings.Index(input, match[0])+len(match[0]):]
		time = strings.TrimSpace(time)
		return match[0], time
	}
	return "", ""
}
