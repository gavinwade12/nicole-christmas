package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	nc "github.com/gavinwade12/nicole-christmas"
	"github.com/gavinwade12/nicole-christmas/twilio"
)

var srv *http.Server

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
		"sweets":      "I smell great, but I'm nothing like momma's cookin'.",
		"merlot":      "There are only a few people that live here, but we have so much stuff!",
		"gingerbread": "If you wear me, where would you be?",
		"pink":        "It's might cold traveling in part of my elephant. No need to worry; I can keep you warm.",
		"warm christmas": "I hope you had fun! There's one last present, and this is the longest-lasting one. " +
			"Expect a love quote every day from now until the end of time. Enjoy and I love you!",
	}

	srv = &http.Server{Addr: ":8080"}

	http.HandleFunc("/christmas", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if err := r.ParseForm(); err != nil {
			shutdown(err)
		}
		data := r.FormValue("Body")

		keyword := strings.ToLower(string(data))
		if keyword == "start" {
			err := tc.Send(greeting, true)
			if err != nil {
				shutdown(err)
			}
			return
		}

		hint := keywords[keyword]
		if hint == "" {
			err := tc.Send("Sorry, I don't know that keyword! Try again. :-)", true)
			if err != nil {
				shutdown(err)
			}
			return
		}

		err = tc.Send(hint, true)
		if err != nil {
			shutdown(err)
		}

		if keyword == "warm christmas" {
			err := nc.SendQuote(tc)
			if err != nil {
				shutdown(err)
			}
			shutdown(nil)
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

func shutdown(err error) {
	sherr := srv.Shutdown(context.Background())
	if sherr != nil {
		if err != nil {
			log.Fatalf("failed shutting down: %v --- reason for shutdown: %v", sherr, err)
		}
		log.Fatal(sherr)
	}
	if err != nil {
		log.Fatal(err)
	}
}
