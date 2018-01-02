package main

import (
	"fmt"
	mhook "hookbot/mailhook"
	"log"
	"strconv"
	"time"

	DB "hookbot/db"

	"github.com/spf13/viper"
	tb "gopkg.in/tucnak/telebot.v2"
)

func waitForPdfs(m mhook.Mailhook, b *tb.Bot, db DB.DB) {
	for file := range m.GetPdfs() {
		func() {
			defer file.Destroy()
			fmt.Printf("New file received: %s\n", file.Path())

			tbFile := &tb.Document{File: tb.FromDisk(file.Path())}

			subs, err := db.GetSubscribers()
			if err != nil {
				fmt.Printf("Something went wrong while retriving subscribers %s\n", err.Error())
				return
			}

			for _, recipient := range subs {
				_, err := b.Send(recipient, tbFile)

				if err != nil {
					fmt.Printf("Something went wrong while sending attachment %s\n", err.Error())
				} else {
					fmt.Printf("Attachment sent successfully\n")
				}
			}
		}()
	}
}

func initizeConfigurations() {
	viper.SetDefault("Endpoint", "/mailhook")
	viper.SetDefault("Host", ":3000")
	viper.SetDefault("DownloadFolder", "/tmp")
	viper.SetDefault("BotApiToken", "")
	viper.SetDefault("DbPath", "./subscribers.db")
	viper.SetDefault("Commands", map[string]string{
		"subscribe":   "/subscribe",
		"unsubscribe": "/unsubscribe",
	})
	viper.SetDefault("Responses", map[string]string{
		"subscribe_success":     "You are now subscribed to Hookbot.",
		"subscribe_duplicate":   "You are already subscribed.",
		"unsubscribe_success":   "You are now unsubscribed to Hookbot.",
		"unsubscribe_duplicate": "You are not yet subscribed.",
	})
	viper.SetDefault("WebhookAuthKeys", []string{})
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/hookbot/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func getRecipient(m *tb.Message) DB.BotRecipient {
	var recipient DB.BotRecipient

	if !m.FromGroup() {
		recipient = DB.RecipientBuilder(m.Sender.Recipient(), false)

	} else {
		recipient = DB.RecipientBuilder(strconv.FormatInt(m.Chat.ID, 10), true)
	}

	return recipient
}

func main() {
	initizeConfigurations()

	db := DB.DbBuilder(viper.GetString("DbPath"))

	b, err := tb.NewBot(tb.Settings{
		Token:  viper.GetString("BotApiToken"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	commands := viper.GetStringMapString("Commands")
	responses := viper.GetStringMapString("Responses")

	b.Handle(commands["subscribe"], func(m *tb.Message) {
		recipient := getRecipient(m)

		if db.IsSubscribed(recipient.Recipient()) {
			b.Send(recipient, responses["subscribe_duplicate"])
			return
		}

		err := db.AddSubscriberRecipient(recipient)
		if err != nil {
			fmt.Printf("Something went wrong while subscribing user %s\n", err.Error())
		} else {
			b.Send(recipient, responses["subscribe_success"])
		}
	})

	b.Handle(commands["unsubscribe"], func(m *tb.Message) {
		recipient := getRecipient(m)

		if !db.IsSubscribed(recipient.Recipient()) {
			b.Send(recipient, responses["unsubscribe_duplicate"])
			return
		}

		err := db.RemoveSubscriber(recipient.Recipient())
		if err != nil {
			fmt.Printf("Something went wrong while unsubscribing user %s\n", err.Error())
		} else {
			b.Send(recipient, responses["unsubscribe_success"])
		}
	})

	go b.Start()

	var m = mhook.MailhookBuilder()
	go waitForPdfs(m, b, db)

	m.EnableWebHook(
		viper.GetString("Endpoint"),
		viper.GetString("Host"),
		viper.GetString("DownloadFolder"),
		viper.GetStringSlice("WebhookAuthKeys"),
	)
}
