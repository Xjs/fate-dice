package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/Xjs/fate-dice/fate"
)

func main() {
	var token string
	var defaultN = 4

	flag.StringVar(&token, "token", token, "Bot Token")
	flag.IntVar(&defaultN, "number", defaultN, "Number of dice to throw when not specified in message")
	flag.Parse()

	if token == "" {
		log.Fatalln("-token must be specified")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}
	defer dg.Close()

	callback := makeMessageCreateCallback(defaultN)

	removeCallback := dg.AddHandler(callback)
	defer removeCallback()

	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
		return
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func init() {
	fate.Seed()
}

// makeMessageCreateCallback creates a MessageCreate callback function.
func makeMessageCreateCallback(defaultN int) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, "dice") {
			arg := strings.TrimSpace(strings.TrimPrefix(m.Content, "dice"))
			n, err := strconv.Atoi(arg)
			if err != nil {
				n = defaultN
			}
			resultString, result := fate.Fate(n)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, total: %d", resultString, result))
		}
	}
}
