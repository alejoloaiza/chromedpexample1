// search.
package main

import (
	"fmt"
	"os"
)

var startdate, enddate, origin, destination, logpath string

func main() {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			_, ok = r.(error)
			if !ok {
				fmt.Errorf("pkg: %v", r)
			}
		}
	}()
	if len(os.Args) > 5 {
		startdate = os.Args[1]
		enddate = os.Args[2]
		origin = os.Args[3]
		destination = os.Args[4]
		logpath = os.Args[5]
	} else {
		startdate = "190123"
		enddate = "190228"
		origin = "mdea"
		destination = "syda"
		logpath = "results.csv"
	}

	InitCaptureOneWay(logpath, startdate, enddate, origin, destination)
}
