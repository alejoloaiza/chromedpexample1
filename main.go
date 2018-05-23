// search.
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf), chromedp.WithRunnerOptions(
		//runner.Flag("headless", true),
		runner.Flag("disable-gpu", true),
		runner.Flag("no-first-run", true),
		runner.Flag("no-default-browser-check", true),
		runner.Port(9222),
	))
	if err != nil {
		log.Fatal(err)
	}
	var best, cheaper, fastest string
	mydate := "180801"
	var url = getNewUrl(&mydate)
	err = c.Run(ctxt, skyscannerSearch(url, &best, &cheaper, &fastest))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("EXTRACTION 1: " + best)
	fmt.Println("EXTRACTION 2: " + cheaper)
	fmt.Println("EXTRACTION 3: " + fastest)

	url = getNewUrl(&mydate)
	err = c.Run(ctxt, skyscannerSearch(url, &best, &cheaper, &fastest))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("EXTRACTION 1: " + best)
	fmt.Println("EXTRACTION 2: " + cheaper)
	fmt.Println("EXTRACTION 3: " + fastest)

	// sutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
func getNewUrl(basenum *string) string {
	var base int
	base, _ = strconv.Atoi(*basenum)
	base++
	*basenum = strconv.Itoa(base)

	url := fmt.Sprintf(`https://www.skyscanner.net/transport/flights/mdea/syda/%s/?adults=1&children=0&adultsv2=1&childrenv2=&infants=0&cabinclass=economy&rtn=0&preferdirects=false&outboundaltsenabled=false&inboundaltsenabled=false&ref=home#results`, *basenum)

	return url

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

func googleSearch(q, text string, site, res *string) chromedp.Tasks {
	var buf []byte
	sel := fmt.Sprintf(`//a[text()[contains(., '%s')]]`, text)
	return chromedp.Tasks{
		chromedp.Navigate(`https://www.google.com`),
		chromedp.WaitVisible(`#hplogo`, chromedp.ByID),
		chromedp.SendKeys(`#lst-ib`, q+"\n", chromedp.ByID),
		chromedp.WaitVisible(`#res`, chromedp.ByID),
		chromedp.Text(sel, res),
		chromedp.Click(sel),
		chromedp.Screenshot(`h1.fade-up`, &buf, chromedp.ByQuery),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile("screenshot.png", buf, 0644)
		}),
	}
}
func skyscannerSearch(url string, best *string, cheapest *string, fastest *string) chromedp.Tasks {
	//var buf []byte
	return chromedp.Tasks{
		chromedp.Navigate(url),
		//chromedp.WaitVisible(`td.tab.active`),
		chromedp.Sleep(15 * time.Second),
		//chromedp.WaitVisible(`#fqs-tabs > table > tbody > tr > td.tab.active`),
		//chromedp.WaitNotVisible(`span.progress-spinner.hot-spinner`),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			log.Printf(">>>>>>>>>>>>>>>>>>>> BOX1 IS VISIBLE")
			return nil
		}),
		chromedp.Text(`#fqs-tabs > table > tbody > tr > td.tab.active`, best, chromedp.NodeVisible, chromedp.ByID),
		chromedp.Text(`#fqs-tabs > table > tbody > tr > td:nth-child(2)`, cheapest, chromedp.NodeVisible, chromedp.ByID),
		chromedp.Text(`#fqs-tabs > table > tbody > tr > td:nth-child(3)`, fastest, chromedp.NodeVisible, chromedp.ByID),
		/*chromedp.Screenshot(`td.tab.active`, &buf, chromedp.ByQuery),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile("screenshot.png", buf, 0644)
		}),*/
	}
}

/*
func main() {
	var err error

	// Define a chrome instance with remote debugging enabled.
	browser := chrome.New(
		// See https://developers.google.com/web/updates/2017/04/headless-chrome#cli
		// for details about startup flags
		&chrome.Flags{
			"addr":               "localhost",
			"disable-extensions": nil,
			"disable-gpu":        nil,
			"headless":           nil,
			"hide-scrollbars":    nil,
			"no-first-run":       nil,
			"no-sandbox":         nil,
			"port":               9222,
			"remote-debugging-address": "0.0.0.0",
			"remote-debugging-port":    9222,
		},
		"/usr/bin/google-chrome", // Path to Chromeium binary
		"/tmp",      // Set the Chromium working directory
		"/dev/null", // Ignore internal Chromium output, set to empty string for os.Stdout
		"/dev/null", // Ignore internal Chromium errors, set to empty string for os.Stderr
	)

	// Start the chrome process.
	if err := browser.Launch(); nil != err {
		panic(err)
	}

	// Open a tab and navigate to the URL you want to screenshot.
	tab, err := browser.NewTab("https://www.google.com")
	if nil != err {
		panic(err)
	}

	// Enable Page events for this tab.
	if enableResult := <-tab.Page().Enable(); nil != enableResult.Err {
		panic(enableResult.Err)
	}

	// Create a channel to receive the screenshot data generated by the
	// event handler.
	results := make(chan *page.CaptureScreenshotResult)

	// Create an event handler that will execute when the page load event
	// fires with a closure that will capture the screenshot and write
	// the result to the results channel. This can also be done manually
	// via:
	//	eventHandler := socket.NewEventHandler(
	//		"Page.loadEventFired",
	//		func(response *socket.Response) {
	//			...
	//			results <- screenshotResult
	//		},
	//	)
	//	tab.AddEventHandler(eventHandler)
	tab.Page().OnLoadEventFired(
		// This function will generate a screenshot and write the data
		// to the results channel.
		func(event *page.LoadEventFiredEvent) {

			// Set the device emulation parameters.
			overrideResult := <-tab.Emulation().SetDeviceMetricsOverride(
				&emulation.SetDeviceMetricsOverrideParams{
					Width:  1440,
					Height: 1440,
					ScreenOrientation: &emulation.ScreenOrientation{
						Type:  emulation.OrientationType.PortraitPrimary,
						Angle: 90,
					},
				},
			)
			if nil != overrideResult.Err {
				panic(overrideResult.Err)
			}

			// Capture a screenshot of the current state of the current
			// page.
			screenshotResult := <-tab.Page().CaptureScreenshot(
				&page.CaptureScreenshotParams{
					Format:  page.Format.Jpeg,
					Quality: 50,
				},
			)
			if nil != screenshotResult.Err {
				panic(screenshotResult.Err)
			}

			results <- screenshotResult
		},
	)

	// Wait for the handler to fire
	result := <-results

	// Decode the base64 encoded image data
	data, err := base64.StdEncoding.DecodeString(result.Data)
	if nil != err {
		panic(err)
	}

	// Write the generated image to a file
	err = ioutil.WriteFile("example.jpg", data, 0644)
	if nil != err {
		panic(err)
	}

	fmt.Println("Finished rendering example.jpg")
}
*/
