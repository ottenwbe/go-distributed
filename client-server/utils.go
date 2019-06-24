package main

import (
	log "github.com/sirupsen/logrus"
	"os"
)

// how the app handles errors; can be mocked during tests
var exitFunc = os.Exit

func exitOnError(err error, code int) {
	if err != nil {
		log.WithError(err).Errorf("%v: Stopping client server demo due to error", code)
		exitFunc(code)
	}
}

func reverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
