package main

import (
	"fmt"
	"log"
	"os"
	"replacebot/internal/replacebot"

	twitch "github.com/gempir/go-twitch-irc/v4"
	"gopkg.in/yaml.v3"
)

func main() {
	var cfg replacebot.Config
	c, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("failed to read YAML file: %v", err)
	}
	if err := yaml.Unmarshal(c, &cfg); err != nil {
		log.Fatalf("failed to unmarshal YAML: %v", err)
	}

	bot := replacebot.NewReplaceBot(&cfg)
	client := twitch.NewClient(cfg.UserName, cfg.OAuth)

	client.OnPrivateMessage(func(m twitch.PrivateMessage) {
		r := bot.Respond(m.Message)
		if r != nil {
			client.Say(m.Channel, *r)
		}
	})

	client.OnConnect(func() {
		fmt.Println("Connected to Twitch chat.")
		for _, ch := range cfg.Channels {
			client.Join(ch)
		}
	})

	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}
}
