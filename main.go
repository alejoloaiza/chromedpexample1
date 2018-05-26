// search.
package main

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
	//SplitPrices("Best$1.21033h 44")
	var startdate, enddate, origin, destination string
	startdate = "181001"
	enddate = "181231"
	origin = "mdea"
	destination = "syda"
	InitCapture("results.csv", startdate, enddate, origin, destination)
}

func screenshot(urlstr, sel string) chromedp.Tasks {
	var buf []byte
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(5 * time.Second),
		chromedp.WaitVisible(sel, chromedp.ByID),
		chromedp.Screenshot(sel, &buf, chromedp.ByID),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile("screenshot.png", buf, 0644)
		}),
	}
}
