package main

import (
	"log"
	"net/http"
	"time"

	nc "github.com/gavinwade12/nicole-christmas"
	"github.com/gavinwade12/nicole-christmas/twilio"
)

func main() {
	client := &http.Client{Timeout: time.Second * 60}
	tc, err := twilio.NewClient(client, "../../"+nc.TwilioConfigLocation)
	if err != nil {
		log.Fatal(err)
	}

	err = nc.SendQuote(tc)
	if err != nil {
		log.Fatal(err)
	}
}
