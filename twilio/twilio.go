package twilio

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// Config holds the configuration options for the Twilio client.
type Config struct {
	AccountSID string `json:"account_sid"`
	AuthToken  string `json:"auth_token"`
	FromNumber string `json:"from_number"`
	ToNumber   string `json:"to_number"`
	LogNumber  string `json:"log_number"`
}

// Client is used to send sms messages via the Twilio REST api.
type Client struct {
	cfg    Config
	client *http.Client
	smsURL string
}

// NewClient attempts to return an instantiated client. It can fail trying to open
// or decode the config file.
func NewClient(client *http.Client, configLocation string) (*Client, error) {
	cfgFile, err := os.Open(configLocation)
	if err != nil {
		return nil, errors.Wrap(err, "error opening twilio config file")
	}
	defer cfgFile.Close()

	cfg := Config{}
	err = json.NewDecoder(cfgFile).Decode(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding twilio config json")
	}

	tc := &Client{
		cfg:    cfg,
		client: client,
		smsURL: "https://api.twilio.com/2010-04-01/Accounts/" + cfg.AccountSID + "/Messages.json",
	}

	return tc, nil
}

// Send will send the message to the ToNumber in the Config.
// If logging is true and there is an error sending the message,
// an attempt to send the error to the LogNumber in the Config
// will be made.
func (c *Client) Send(msg string, logError bool) error {
	err := c.send(msg, c.cfg.ToNumber)
	if !logError || err == nil {
		return err
	}

	err = fmt.Errorf("error sending msg: %v", err)
	logErr := c.Log(err.Error())
	if logErr != nil {
		return fmt.Errorf("error sending log message: %v: %v", logErr, err)
	}
	return err
}

// Log will send the message to the LogNumber in the Config.
func (c *Client) Log(msg string) error {
	return c.send(msg, c.cfg.LogNumber)
}

func (c *Client) send(msg, to string) error {
	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", c.cfg.FromNumber)
	msgData.Set("Body", msg)
	msgDataReader := strings.NewReader(msgData.Encode())

	req, err := http.NewRequest(http.MethodPost, c.smsURL, msgDataReader)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.cfg.AccountSID, c.cfg.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading body of request (status: %d): %v", resp.StatusCode, err)
		}
		return fmt.Errorf("got status code: %d\nbody: %s", resp.StatusCode, data)
	}
	return nil
}
