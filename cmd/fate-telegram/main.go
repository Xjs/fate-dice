package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Xjs/fate-dice/dnd"
	"github.com/Xjs/fate-dice/fate"
	"gopkg.in/tucnak/telebot.v2"
)

func helpFunc(bot *telebot.Bot) func(m *telebot.Message) {
	return func(m *telebot.Message) {
		response := `/dice <offset> <comment>: Throw 4 fate dice and add the given offset.
/dnd <NdM specification>: Throw DnD-style N dice with M faces each (e. g. 2d6 to throw 🎲🎲)
	`
		if _, err := bot.Send(m.Chat, response,
			&telebot.SendOptions{ParseMode: telebot.ModeMarkdown}); err != nil {
			log.Printf("Error sending %s to %v: %v\n", response, m.Chat, err)
		}
	}
}

func main() {
	var token string
	var apiURL string = "https://api.telegram.org"
	var timeout time.Duration
	var defaultN = 4

	flag.IntVar(&defaultN, "number", defaultN, "Number of dice to throw when not specified in message")
	flag.StringVar(&token, "token", token, "Telegram API token")
	flag.StringVar(&apiURL, "api", apiURL, "Telegram API URL")
	flag.DurationVar(&timeout, "timeout", timeout, "Poller timeout")

	flag.Parse()

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  token,
		URL:    apiURL,
		Poller: &telebot.LongPoller{Timeout: timeout},
	})
	if err != nil {
		log.Fatalf("Error initialising bot: %v\n", err)
	}

	bot.Handle("/dice", func(m *telebot.Message) {
		words := strings.Split(m.Text, " ")[1:]
		n := defaultN

		offset := 0
		if len(words) > 0 && len(words[0]) > 0 {
			o, err := strconv.Atoi(words[0])
			if err != nil {
				log.Printf("invalid offset after dice: %s\n", words[0])
				return
			}
			offset = o
		}

		var comment string
		if len(words) > 1 {
			comment = " " + strings.Join(words[1:], " ")
		}

		resultString, result := fate.Fate(n)
		response := fmt.Sprintf("%s, total: %d%s", resultString, result+offset, comment)
		if _, err := bot.Send(m.Chat, response,
			&telebot.SendOptions{ParseMode: telebot.ModeMarkdown}); err != nil {
			log.Printf("Error sending %s to %v: %v\n", response, m.Chat, err)
		}
	})

	bot.Handle("/dnd", func(m *telebot.Message) {
		text := strings.TrimSpace(m.Text[len("/dnd"):])

		if text == "" {
			if _, err := bot.Send(m.Chat, telebot.Cube,
				&telebot.SendOptions{ParseMode: telebot.ModeMarkdown}); err != nil {
				log.Printf("Error sending %s to %v: %v\n", telebot.Cube, m.Chat, err)
			}
			return
		}

		t, err := dnd.Parse(text)
		if err != nil {
			log.Println(err)
			return
		}

		dice := text
		if t.Faces == 6 && t.Dice < 10 {
			dice = t.Emoji()
		}

		response := fmt.Sprintf("%s: %d", dice, t.Throw())
		if _, err := bot.Send(m.Chat, response,
			&telebot.SendOptions{ParseMode: telebot.ModeMarkdown}); err != nil {
			log.Printf("Error sending %s to %v: %v\n", response, m.Chat, err)
		}
	})

	bot.Handle("/help", helpFunc(bot))
	bot.Handle("/start", helpFunc(bot))

	log.Println("Running.")
	bot.Start()
}
