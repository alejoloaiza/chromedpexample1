package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

func InitCapture(logpath string, initdate string, enddate string) {
	var err error
	var cycle int
	var best, cheaper, fastest string
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
	for initdate != enddate {
		cycle++
		var url = getNewUrl(&initdate)
		err = c.Run(ctxt, skyscannerSearch(url, &best, &cheaper, &fastest))
		if err != nil {
			log.Fatal(err)
		}

		newline := fmt.Sprintf("EXTRACTION %d: BEST %s CHEAP %s FAST %s \n", cycle, best, cheaper, fastest)
		WriteResults(logpath, newline)

	}

	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
func skyscannerSearch(url string, best *string, cheapest *string, fastest *string) chromedp.Tasks {
	//var buf []byte
	return chromedp.Tasks{
		chromedp.Navigate(url),
		//chromedp.WaitVisible(`td.tab.active`),
		//chromedp.Sleep(15 * time.Second),
		chromedp.WaitVisible(`#fqs-tabs > table > tbody > tr > td.tab.active`),
		//chromedp.WaitNotVisible(`span.progress-spinner.hot-spinner`),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			log.Printf(">>>>>>>>>>>>>>>>>>>> BOX1 IS VISIBLE")
			return nil
		}),
		//chromedp.Sleep(3 * time.Second),
		chromedp.Text(`#fqs-tabs > table > tbody > tr > td.tab.active`, best, chromedp.NodeVisible, chromedp.ByID),
		chromedp.Text(`#fqs-tabs > table > tbody > tr > td:nth-child(2)`, cheapest, chromedp.NodeVisible, chromedp.ByID),
		chromedp.Text(`#fqs-tabs > table > tbody > tr > td:nth-child(3)`, fastest, chromedp.NodeVisible, chromedp.ByID),
		/*chromedp.Screenshot(`td.tab.active`, &buf, chromedp.ByQuery),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile("screenshot.png", buf, 0644)
		}),*/
	}
}
func getNewUrl(basenum *string) string {
	const shortForm = "060102"
	t, _ := time.Parse(shortForm, *basenum)
	t = t.Add(time.Hour * 24)
	*basenum = t.Format(shortForm)

	url := fmt.Sprintf(`https://www.skyscanner.net/transport/flights/mdea/syda/%s/?adults=1&children=0&adultsv2=1&childrenv2=&infants=0&cabinclass=economy&rtn=0&preferdirects=false&outboundaltsenabled=false&inboundaltsenabled=false&ref=home#results`, *basenum)

	return url
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
