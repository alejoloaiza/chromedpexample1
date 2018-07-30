// search.
package main

import (
	"fmt"
	"os"
	"strconv"
)

var startdate, enddate, origin, destination, mid, logpath string
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
			startdate = "180801"
			enddate = "181031"
			origin = "sins"
			destination = "syda"
			logpath = "resultssin-syd.csv"
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
			startdate = "181120"
			enddate = "190110"
			aproxduration = 45
			origin = "syda"
			destination = "mdea"
			logpath = "resultsjuan.csv"
		}
		InitCaptureReturn(logpath, startdate, enddate, aproxduration, origin, destination)

	} else if os.Args[1] == "multiple" {
		if len(os.Args) > 6 {
			startdate = os.Args[2]
			enddate = os.Args[3]
			origin = os.Args[4]
			mid = os.Args[5]
			destination = os.Args[6]
			logpath = os.Args[7]
		} else {
			startdate = "2018-08-01"
			enddate = "2018-10-31"
			origin = "mdea"
			mid = "scla"
			destination = "syda"
			logpath = "resultsmed-scl-syd.csv"
		}
		InitCaptureMultiple(logpath, startdate, enddate, origin, mid, destination)

	}

}
