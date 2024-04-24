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
	"gopkg.in/telebot.v3"
)

func helpFunc(ctx telebot.Context) error {
	response := `/dice <offset> <comment>: Throw 4 fate dice and add the given offset.
/dnd <NdM specification>: Throw DnD-style N dice with M faces each (e. g. 2d6 to throw 🎲🎲)
	`
	if err := ctx.Send(response); err != nil {
		return fmt.Errorf("error sending %s: %w", response, err)
	}

	return nil
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

	bot.Handle("/dice", func(ctx telebot.Context) error {
		words := ctx.Args()
		n := defaultN

		offset := 0
		if len(words) > 0 && len(words[0]) > 0 {
			o, err := strconv.Atoi(words[0])
			if err != nil {
				return fmt.Errorf("invalid offset after dice: %s", words[0])
			}
			offset = o
		}

		var comment string
		if len(words) > 1 {
			comment = " " + strings.Join(words[1:], " ")
		}

		resultString, result := fate.Fate(n)
		response := fmt.Sprintf("%s, total: %d%s", resultString, result+offset, comment)
		if err := ctx.Send(response); err != nil {
			return fmt.Errorf("error sending %sv: %w", response, err)
		}

		return nil
	})

	bot.Handle("/dnd", func(ctx telebot.Context) error {
		text := strings.TrimSpace(ctx.Message().Payload)

		if text == "" {
			if err := ctx.Send(telebot.Cube,
				&telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2}); err != nil {
				return fmt.Errorf("error sending %v: %w", telebot.Cube, err)
			}
			return nil
		}

		t, err := dnd.Parse(text)
		if err != nil {
			return err
		}

		dice := text
		if t.Faces == 6 && t.Dice < 10 {
			dice = t.Emoji()
		}

		response := fmt.Sprintf("%s: %d", dice, t.Throw())
		if err := ctx.Send(response,
			&telebot.SendOptions{ParseMode: telebot.ModeMarkdownV2}); err != nil {
			return fmt.Errorf("error sending %s: %w", response, err)
		}

		return nil
	})

	bot.Handle("/help", helpFunc)
	bot.Handle("/start", helpFunc)

	log.Println("Running.")
	bot.Start()
}
