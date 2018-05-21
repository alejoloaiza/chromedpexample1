// search.
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func main() {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var site, res string
	//err = c.Run(ctxt, googleSearch("site:brank.as", "Home", &site, &res))
	err = c.Run(ctxt, skyscannerSearch())
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("saved screenshot from search result listing `%s` (%s)", res, site)
}

func googleSearch(q, text string, site, res *string) chromedp.Tasks {
	var buf []byte
	sel := fmt.Sprintf(`//a[text()[contains(., '%s')]]`, text)
	fmt.Println(sel)
	fmt.Println("=======================")
	return chromedp.Tasks{
		chromedp.Navigate(`https://www.google.com`),
		chromedp.WaitVisible(`span.fqs-price`, chromedp.ByID),
		chromedp.SendKeys(`#lst-ib`, q+"\n", chromedp.ByID),
		chromedp.WaitVisible(`#res`, chromedp.ByID),
		chromedp.Text(sel, res),
		chromedp.Click(sel),
		chromedp.WaitNotVisible(`.preloader-content`, chromedp.ByQuery),
		chromedp.WaitVisible(`a[href*="twitter"]`, chromedp.ByQuery),
		chromedp.Location(site),
		chromedp.ScrollIntoView(`.banner-section.third-section`, chromedp.ByQuery),
		chromedp.Sleep(2 * time.Second), // wait for animation to finish
		chromedp.Screenshot(`.banner-section.third-section`, &buf, chromedp.ByQuery),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile("screenshot.png", buf, 0644)
		}),
	}
}
func skyscannerSearch() chromedp.Tasks {
	var buf []byte
	return chromedp.Tasks{
		chromedp.Navigate(`https://www.skyscanner.net/transport/flights/mdea/syda/180801/?adults=1&children=0&adultsv2=1&childrenv2=&infants=0&cabinclass=economy&rtn=0&preferdirects=false&outboundaltsenabled=false&inboundaltsenabled=false&ref=home#results`),
		chromedp.WaitVisible(`span.fqs-price`),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			log.Printf(">>>>>>>>>>>>>>>>>>>> BOX1 IS VISIBLE")
			return nil
		}),
		chromedp.Screenshot(`span.fqs-price`, &buf, chromedp.ByQuery),
		chromedp.ActionFunc(func(context.Context, cdp.Executor) error {
			return ioutil.WriteFile("screenshot.png", buf, 0644)
		}),
	}
}
