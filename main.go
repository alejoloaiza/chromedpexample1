// search.
package main

func main() {
	var startdate, enddate, origin, destination string
	startdate = "181231"
	enddate = "190228"
	origin = "mdea"
	destination = "syda"
	InitCapture("results.csv", startdate, enddate, origin, destination)
}
