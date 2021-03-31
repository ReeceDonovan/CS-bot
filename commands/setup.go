package commands

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

type commandFunc func(session *discordgo.Session, message *discordgo.MessageCreate)

var commandsMap = make(map[string]func(*discordgo.Session, *discordgo.MessageCreate))
var helpMsgs = make(map[string]string)

func command(name string, helpMessage string, function commandFunc) {
	helpMsgs[name] = helpMessage
	commandsMap[name] = function
}

func RegisterCommands(s *discordgo.Session) {
	command("link", "Link server", Link)
	s.AddHandler(messageCreate)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	if !strings.HasPrefix(m.Content, viper.GetString("discord.prefix")) {
		return
	}
	callCommand(s, m)
}

func callCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	commandStr, _ := extractCommand(m.Content)
	if command, ok := commandsMap[commandStr]; ok {
		log.Println("Command Triggered")
		command(s, m)
		return
	}
}

func extractCommand(c string) (commandStr string, body string) {
	body = strings.TrimPrefix(c, viper.GetString("discord.prefix"))
	commandStr = strings.Fields(body)[0]
	return
}