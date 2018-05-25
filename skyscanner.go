package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

const shortForm = "060102"

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
	WriteResults(logpath, "Time;Iteration;Date;Best;best Time;Cheapest;Cheapest Time;Fastest;Fastest Time")
	for initdate != enddate {
		cycle++
		var url = getNewUrl(&initdate)
		err = c.Run(ctxt, skyscannerSearch(url, &best, &cheaper, &fastest))
		if err != nil {
			log.Fatal(err)
		}
		price1, duration1 := SplitPrices(best)
		price2, duration2 := SplitPrices(cheaper)
		price3, duration3 := SplitPrices(fastest)
		newline := fmt.Sprintf("%s;%d;%s;%s;%s;%s;%s;%s;%s", time.Now().Format("02-01-2006 15:04:05"), cycle, initdate, price1, duration1, price2, duration2, price3, duration3)
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

	sel := fmt.Sprintf(`//span[text()[contains(., '%s')]]`, "results sorted by")

	return chromedp.Tasks{
		chromedp.Navigate(url),

		//chromedp.WaitVisible(`#fqs-tabs > table > tbody > tr > td.tab.active`),

		chromedp.WaitVisible(sel),
		//chromedp.Text(`#header-list-count > div > span`, &results, nil, chromedp.ByID),
		//chromedp.Sleep(3 * time.Second),//*[@id="header-list-count"]/div/span/strong/span//*[@id="header-list-count"]/div/span/strong/span
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

	t, _ := time.Parse(shortForm, *basenum)
	t = t.Add(time.Hour * 24)
	*basenum = t.Format(shortForm)

	//url := fmt.Sprintf(`https://www.skyscanner.net/transport/flights/mdea/syda/%s/?currency=USD&adults=1&children=0&adultsv2=1&childrenv2=&infants=0&cabinclass=economy&rtn=0&preferdirects=false&outboundaltsenabled=false&inboundaltsenabled=false&ref=home#results`, *basenum)
	url := fmt.Sprintf(`https://www.skyscanner.net/transport/flights/mdea/ytoa/%s/?currency=USD&adults=1&children=0&adultsv2=1&childrenv2=&infants=0&cabinclass=economy&rtn=0&preferdirects=false&outboundaltsenabled=false&inboundaltsenabled=false&ref=home#results`, *basenum)

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
