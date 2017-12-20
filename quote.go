package christmas

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gavinwade12/nicole-christmas/twilio"
)

const (
	// QuoteURL is the URL used to retrieve quotes.
	QuoteURL = "https://www.romanticlovemessages.com/random/random.php"
	// TwilioConfigLocation is the location of the twilio config file.
	TwilioConfigLocation = "twilio.config"
)

// SendQuote will get a quote from the QuoteURL and send it via
// the provided twilio client.
func SendQuote(tc *twilio.Client) error {
	quote, err := getQuote()
	if err != nil {
		return err
	}

	return tc.Send(quote, true)
}

func getQuote() (string, error) {
	resp, err := http.Post(QuoteURL, "", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if err != nil {
			return "", fmt.Errorf("error reading body of request with status code: %d", resp.StatusCode)
		}
		return "", fmt.Errorf("got status code: %d\nbody: %s", resp.StatusCode, data)
	}

	bodyString := string(data)
	garbageIndex := strings.LastIndex(bodyString, "<br>")
	quote := bodyString[:garbageIndex]
	quote = strings.Replace(quote, "<br>", "\n", -1)
	return quote, nil
}

// Heated blanket - trunk: "It's mighty cold traveling in part of my elephant. No need to worry; I can keep you warm." - keyword: "warm christmas" - next: quote
// Candy - under tree - keyword: "sweets" - next: candle
// Gingerbread village - old bedroom: "There are only a few people that live here, but we have so much stuff!" - keyword: "gingerbread" - next: shirt
// Candle - oven: "I smell great, but I'm nothing like momma's cookin'." - keyword: "merlot" - next: gingerbread
// Shirt - closet: "If you wear me, where would you be?" - keyword: "pink" - next: blanket
