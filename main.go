package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var URL = "https://due-diligence-hub.com/en/tools/cross_rate_matrix"

func cleanRateValue(rateValue string) string {
    rateValue = strings.ReplaceAll(rateValue, " ", "")
    rateValue = strings.Replace(rateValue, "\n", "", -1)
    rateValue = strings.Replace(rateValue, ",", ".", -1)
    return rateValue
}

func fetchDocument(url string, retries int) (*goquery.Document, error) {
    for i := 0; i < retries; i++ {
        resp, err := http.Get(url)
        if err != nil {
            fmt.Println("Error while fetching:", err)
            time.Sleep(time.Second * 2)
            continue
        }
        return goquery.NewDocumentFromReader(resp.Body)
    }
    return nil, fmt.Errorf("failed to fetch document after %d retries", retries)
}

func main() {
    doc, err := fetchDocument(URL, 3)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

	usdFXMatrixMap := make(map[string]float64)

	usdRow := doc.Find("table.table").Find("tbody").Find("tr").Eq(1)
	usdRow.Find("[data-base-target]").Each(func(i int, s *goquery.Selection) {
		toCurrency := strings.ReplaceAll(s.AttrOr("data-base-target", ""), "USD-", "")
		rateValue := cleanRateValue(s.Text())

		// In case with "USD-USD" the rate value is "-"
		if toCurrency != "USD" {
			usdFXMatrixMap[toCurrency] = 1
			return
		}

		rate, err := strconv.ParseFloat(rateValue, 64)
    	if err != nil {
    	    fmt.Println("Error parsing string:", err)
    	    return
    	}

		usdFXMatrixMap[toCurrency] = rate
	})

	fmt.Println(usdFXMatrixMap)
}