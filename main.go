package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-steputils/stepconf"
)

// An object represting the Discord WebHook
type DiscordWebhook struct {
	Content string `json:"content"`
}

func newDiscordWebhook(c Config) DiscordWebhook {
	webhook := DiscordWebhook{
		Content: c.Message,
	}
	return webhook
}

func postMessage(webhookURL string, webhook DiscordWebhook) error {
	b, err := json.Marshal(webhook)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("Failed to send the request: %s", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); err == nil {
			err = cerr
		}
	}()

	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Server error: %s, failed to read response: %s", resp.Status, err)
		}
		return fmt.Errorf("Server error: %s, response: %s", resp.Status, body)
	}

	return nil
}

func main() {
	var conf Config
	if err := stepconf.Parse(&conf); err != nil {
		log.Errorf("Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(conf)

	webhook := newDiscordWebhook(conf)
	if err := postMessage(conf.WebhookURL, webhook); err != nil {
		log.Errorf("Error: %s", err)
		os.Exit(1)
	}

	log.Donef("\nSuccess: Discord message was sent!\n")
}
