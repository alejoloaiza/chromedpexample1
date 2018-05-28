package main

import (
	"os"
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
func WriteResults(logpath string, line string) {
	f, err := os.OpenFile(logpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	newline := line + string('\n')
	_, err = f.WriteString(newline)
	if err != nil {
		panic(err)
	}
}
