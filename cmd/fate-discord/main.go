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

	"github.com/Xjs/fate-dice/dnd"
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

const debug = true

// makeMessageCreateCallback creates a MessageCreate callback function.
func makeMessageCreateCallback(defaultN int) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		if debug {
			log.Printf("Received message %q from %s", m.Message.Content, m.Author.Username)
		}

		if strings.HasPrefix(m.Content, "dice") {
			arg := strings.TrimSpace(strings.TrimPrefix(m.Content, "dice"))

			words := strings.Split(arg, " ")
			n := defaultN
			if numDice := words[0]; len(numDice) >= 3 && numDice[0] == '(' && numDice[len(numDice)-1] == ')' {
				arg := numDice[1 : len(numDice)-1]
				num, err := strconv.Atoi(arg)
				if err != nil {
					log.Printf("invalid argument to dice(): %s\n", arg)
					return
				}
				n = num
				words = words[1:]
			}

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
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, total: %d%s", resultString, result+offset, comment))
		} else if t, err := dnd.Parse(m.Content); err == nil {
			dice := m.Content
			if t.Faces == 6 && t.Dice < 10 {
				dice = t.Emoji()
			}
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s: %d", dice, t.Throw()))
		}
	}
}
