package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gavinwade12/nicole-christmas/twilio"
)

const (
	quoteURL         = "https://www.romanticlovemessages.com/random/random.php"
	twilioConfigFile = "twilio.config"
)

func main() {
	cfgFile, err := os.Open(twilioConfigFile)
	if err != nil {
		log.Fatalf("error opening twilio config file: %v", err)
	}
	defer cfgFile.Close()

	cfg := twilio.Config{}
	err = json.NewDecoder(cfgFile).Decode(&cfg)
	if err != nil {
		log.Fatalf("error decoding twilio config json: %v", err)
	}

	client := &http.Client{Timeout: time.Second * 60}
	smsURL := "https://api.twilio.com/2010-04-01/Accounts/" + cfg.AccountSID + "/Messages.json"
	tc := twilio.NewClient(cfg, client, smsURL)

	quote, err := getQuote()
	if err != nil {
		err = fmt.Errorf("error getting quote: %v", err)
		logErr := tc.Log(err.Error())
		if logErr != nil {
			log.Fatalf("error sending log message as sms: %v: %v", logErr, err)
		}
		log.Fatal(err)
	}

	err = tc.Send(quote, true)
	if err != nil {
		log.Fatal(err)
	}
}

func okResponseCode(code int) bool {
	return code >= 200 && code <= 299
}

func getQuote() (string, error) {
	resp, err := http.Post(quoteURL, "", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if !okResponseCode(resp.StatusCode) {
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
