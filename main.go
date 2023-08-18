package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	wordleHeaderRegex = regexp.MustCompile(`Wordle \w* [0-6X]\/[0-6]`)
	potatoRegex       = regexp.MustCompile(`🟩`)
	sweetPotatoRegex  = regexp.MustCompile(`🟨`)
	RemoveCommands    = false
)

func init() {
	flag.BoolVar(&RemoveCommands, "rmcmd", false, "Remove all commands after shutdowning or not")
	flag.Parse()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env: %s", err.Error())
	}

	dg, err := discordgo.New(fmt.Sprintf("Bot %s", os.Getenv("BOT_TOKEN")))
	if err != nil {
		fmt.Printf("error creating session: %s\n", err.Error())
		return
	}

	err = dg.Open()
	if err != nil {
		fmt.Printf("error opening websocket: %s\n", err.Error())
	}

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	fmt.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, os.Getenv("GUILD_ID"), v)
		if err != nil {
			fmt.Printf("Cannot create '%v' command: %v", v.Name, err)
			return
		}
		registeredCommands[i] = cmd
	}

	defer dg.Close()

	// Wait here until CTRL-C or other term signal is received.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Press Ctrl+C to exit")
	<-stop

	if RemoveCommands {
		for _, v := range registeredCommands {
			err := dg.ApplicationCommandDelete(dg.State.User.ID, os.Getenv("GUILD_ID"), v.ID)
			if err != nil {
				fmt.Printf("Cannot delete '%v' command: %v", v.Name, err)
				return
			}
		}
	}
}

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "potatodle",
		Description: "adds potatoes to your Wordle",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "wordle-result",
				Description: "Results of the Wordle you'd like to potatofy",
				Required:    true,
			},
		},
	},
}

var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"potatodle": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		fmt.Println(strconv.Quote(i.ApplicationCommandData().Options[0].StringValue()))
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: PotatofyWordle(i.ApplicationCommandData().Options[0].StringValue()),
			},
		})
		if err != nil {
			fmt.Printf("error responding to interaction: %s", err.Error())
		}
	},
}

func PotatofyWordle(message string) string {
	potatoMsg := potatoRegex.ReplaceAllString(message, `🥔`)
	sweetPotatoMsg := sweetPotatoRegex.ReplaceAllString(potatoMsg, `🍠`)
	wordle := wordleHeaderRegex.Split(sweetPotatoMsg, -1)
	header := wordleHeaderRegex.FindString(sweetPotatoMsg)
	return header + strings.ReplaceAll(wordle[1], " ", "\n")
}
