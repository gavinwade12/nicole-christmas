package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

	greeting := "Hello! I'm BenjiBot, and I'm going to help you find your christmas presents. " +
		"Go ahead and start with the one under the tree. Just send me the keyword from each present, " +
		"and I'll send you the hint for your next present. Good luck!"

	keywords := map[string]string{
		"sweets":         "I smell great, but I'm nothing like momma's cookin'.",
		"merlot":         "There are only a few people that live here, but we have so much stuff!",
		"gingerbread":    "If you wear me, where would you be?",
		"pink":           "It's might cold traveling in part of my elephant. No need to worry; I can keep you warm.",
		"warm christmas": "I hope you had fun! There's one last present, and this is the longest-lasting one. Enjoy!",
	}

	srv := &http.Server{Addr: ":8080"}

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		err := tc.Send(greeting, true)
		if err != nil {
			log.Fatal(err)
		}
	})

	http.HandleFunc("/christmas", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		defer r.Body.Close()
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			tc.Log(err.Error())
			return
		}

		keyword := strings.ToLower(string(data))
		hint := keywords[keyword]
		if hint == "" {
			err := tc.Send("Sorry, I don't know that keyword! Try again. :-)", true)
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		err = tc.Send(hint, true)
		if err != nil {
			log.Fatal(err)
		}

		if keyword == "warm christmas" {
			srv.Shutdown(nil)
			err := nc.SendQuote(tc)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello!"))
	})

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
