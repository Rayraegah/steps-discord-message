package main

// Config is populated with the values from the environment variables and user inputs
type Config struct {
	//Debug bool `env:"is_debug_mode,opt[yes,no]""`

	WebhookURL string `env:"webhook_url_input,required"`
	Message    string `env:"webhook_message_input"`
}
