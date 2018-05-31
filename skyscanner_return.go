package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

//const baseUrlReturn = `https://www.skyscanner.net/transport/flights/%s/%s/%s/?currency=%s&adults=%s&children=%s&adultsv2=%s&childrenv2=&infants=0&cabinclass=%s&rtn=0&preferdirects=false&outboundaltsenabled=false&inboundaltsenabled=false&ref=home#results`
const baseUrlReturn = `https://www.skyscanner.net/transport/flights/%s/%s/%s/%s/?currency=%s&adults=%s&children=%s&adultsv2=%s&childrenv2=&infants=0&cabinclass=%s&rtn=1&preferdirects=false&outboundaltsenabled=false&inboundaltsenabled=false&ref=home#results`
const daystomove = 5

var oscillator int
var movearray [daystomove]int

func InitCaptureReturn(logpath string, initdate string, enddate string, duration int, origin string, dest string) {
	var err error
	var cycle int
	var best, cheaper, fastest string
	oscillator = 0
	initArrayToMove()
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf), chromedp.WithRunnerOptions(
		//runner.Flag("headless", true),
		//	runner.UserAgent("Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36"),
		runner.Flag("disable-gpu", true),
		runner.Flag("no-first-run", true),
		runner.Flag("no-default-browser-check", true),
		runner.Port(9222),
	))
	if err != nil {
		log.Fatal(err)
	}
	WriteResults(logpath, "Time;Iteration;Initial Date;Return Date;Best;Best Time;Cheapest;Cheapest Time;Fastest;Fastest Time")
	for initdate != enddate {
		cycle++
		url, returndate := getNewUrlReturn(&initdate, &enddate, origin, dest, duration)
		fmt.Println(url)

		err = c.Run(ctxt, skyscannerSearchReturn(url, &best, &cheaper, &fastest))
		if err != nil {
			log.Fatal(err)
		}
		price1, duration1 := SplitPrices(best)
		price2, duration2 := SplitPrices(cheaper)
		price3, duration3 := SplitPrices(fastest)
		newline := fmt.Sprintf("%s;%d;%s;%s;%s;%s;%s;%s;%s", time.Now().Format("02-01-2006 15:04:05"), cycle, initdate, returndate, price1, duration1, price2, duration2, price3, duration3)
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

func skyscannerSearchReturn(url string, best *string, cheapest *string, fastest *string) chromedp.Tasks {

	sel := fmt.Sprintf(`//span[text()[contains(., '%s')]]`, "results sorted by")
	//var buf []byte
	return chromedp.Tasks{
		chromedp.Navigate(url),

		//chromedp.WaitVisible(`#fqs-tabs > table > tbody > tr > td.tab.active`),
		/*
			chromedp.Sleep(10 * time.Second), //*[@id="header-list-count"]/div/span/strong/span//*[@id="header-list-count"]/div/span/strong/span
			chromedp.Screenshot(`:root`, &buf, chromedp.ByQuery),
			chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
				return ioutil.WriteFile("screenshot.png", buf, 0644)
			}),*/
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
func getNewUrlReturn(datefrom *string, dateto *string, origin string, dest string, duration int) (string, string) {
	var t1, t2 time.Time

FinalDate:
	if oscillator < daystomove {
		t2, _ = time.Parse(shortForm, *datefrom)
		t2 = t2.Add(time.Hour * time.Duration(24*duration))
		t2 = t2.Add(time.Hour * time.Duration(24*movearray[oscillator]))
		//fmt.Println(duration, movearray[oscillator], t2.Format(shortForm))
		//*dateto = t2.Format(shortForm)
		oscillator++
	} else {
		oscillator = 0
		t1, _ = time.Parse(shortForm, *datefrom)
		t1 = t1.Add(time.Hour * 24)
		*datefrom = t1.Format(shortForm)

		goto FinalDate
	}
	url := fmt.Sprintf(baseUrlReturn, origin, dest, *datefrom, t2.Format(shortForm), currency, adults, children, adults, cabinclass)

	return url, t2.Format(shortForm)
}
func initArrayToMove() {

	firstval := int(daystomove/2) * -1
	for i := 0; i < len(movearray); i++ {
		movearray[i] = firstval + i
	}

}
