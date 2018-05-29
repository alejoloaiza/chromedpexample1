// search.
package main

import (
	"fmt"
	"os"
	"strconv"
)

var startdate, enddate, origin, destination, logpath string
var aproxduration int

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
	if os.Args[1] == "oneway" {
		if len(os.Args) > 6 {
			startdate = os.Args[2]
			enddate = os.Args[3]
			origin = os.Args[4]
			destination = os.Args[5]
			logpath = os.Args[6]
		} else {
			startdate = "190123"
			enddate = "190228"
			origin = "mdea"
			destination = "syda"
			logpath = "results.csv"
		}
		InitCaptureOneWay(logpath, startdate, enddate, origin, destination)

	} else if os.Args[1] == "return" {
		if len(os.Args) > 7 {
			startdate = os.Args[2]
			enddate = os.Args[3]
			aproxduration, _ = strconv.Atoi(os.Args[4])
			origin = os.Args[5]
			destination = os.Args[6]
			logpath = os.Args[7]
		} else {
			startdate = "190123"
			enddate = "190228"
			aproxduration = 20
			origin = "mdea"
			destination = "syda"
			logpath = "results.csv"
		}
		InitCaptureReturn(logpath, startdate, enddate, aproxduration, origin, destination)

	}

}
