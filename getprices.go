package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// HistoMinuteData is the type to save response from the API
type HistoMinuteData struct {
	Response   string `json:"Response"`
	Message    string `json:"Message"`
	Type       int    `json:"Type"`
	Aggregated bool   `json:"Aggregated"`
	Data       []struct {
		Time       int     `json:"time"`
		Close      float64 `json:"close"`
		High       float64 `json:"high"`
		Low        float64 `json:"low"`
		Open       float64 `json:"open"`
		Volumefrom float64 `json:"volumefrom"`
		Volumeto   float64 `json:"volumeto"`
	} `json:"Data"`
	TimeTo            int  `json:"TimeTo"`
	TimeFrom          int  `json:"TimeFrom"`
	FirstValueInArray bool `json:"FirstValueInArray"`
	ConversionType    struct {
		Type             string `json:"type"`
		ConversionSymbol string `json:"conversionSymbol"`
	} `json:"ConversionType"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	MMDD := os.Args[1]
	layout := "2006-01-02T15:04:05.000Z"
	str := fmt.Sprintf("2017-%sT00:00:00.000Z", MMDD)
	t, err := time.Parse(layout, str)
	check(err)

	safeTs := url.QueryEscape(strconv.FormatInt(t.Unix(), 10))

	coins := [...]string{"BTC", "ETH", "BCH", "ETC", "XUC", "VERI", "XRP", "DASH", "LTC", "BTG", "IOT", "EOS", "ADA", "XMR", "ETC", "BCCOIN", "NEO", "XEM", "EVR", "XLM", "OMG", "QTUM"}

	for _, coin := range coins {
		url := fmt.Sprintf("https://min-api.cryptocompare.com/data/histominute?fsym=%s&tsym=EUR&toTs=%s", coin, safeTs)

		// Build the request
		req, err := http.NewRequest("GET", url, nil)
		check(err)

		// For control over HTTP client headers,
		// redirect policy, and other settings,
		// create a Client
		// A Client is an HTTP client
		client := &http.Client{}

		// Send the request via a client
		// Do sends an HTTP request and
		// returns an HTTP response
		resp, err := client.Do(req)
		check(err)

		// Callers should close resp.Body
		// when done reading from it
		// Defer the closing of the body
		defer resp.Body.Close()

		// Fill the record with the data from the JSON
		var record HistoMinuteData

		// Use json.Decode for reading streams of JSON data
		if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
			log.Println(err)
		}

		fmt.Println("Response = ", record.Response)
		fmt.Println("Message = ", record.Message)

		fileName := fmt.Sprintf("%s.%s.txt", MMDD, coin)
		f, err := os.Create(fileName)
		check(err)

		defer f.Close()

		w := bufio.NewWriter(f)

		for _, element := range record.Data {
			//tm := time.Unix(int64(element.Time), 0)
			//msg := fmt.Sprintf("%s: %f", tm, element.Close)
			msg := fmt.Sprintf("%d %f\r\n", element.Time, element.Close)
			fmt.Print(msg)

			_, err := w.WriteString(msg)
			check(err)
		}

		w.Flush()

		time.Sleep(100 * time.Millisecond)
	}
}
