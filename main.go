package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"math/rand"
	"time"
)

const (
	Token         = "YOUR_BOT_TOKEN"
	NekoAPIURL    = "https://nekos.life/api/v2/img/"
	DefaultColor  = 0x00ff00 
	ErrorMessage  = "Oops! Something went wrong."
	MaxRandomSeed = 100000
)

var (
	commandList = []string{
		"neko",
	}
)

func main() {
	rand.Seed(time.Now().UnixNano())

	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %s", err)
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		log.Fatalf("Failed to open Discord connection: %s", err)
	}

	defer discord.Close()

	log.Println("Bot is now running. Press CTRL-C to exit.")

	<-make(chan struct{})
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, "!") {
		command := strings.TrimPrefix(m.Content, "!")
		if command == "help" {
			sendHelpMessage(s, m.ChannelID)
		} else if containsString(commandList, command) {
			fetchAndSendImage(s, m.ChannelID, command)
		}
	}
}

func sendHelpMessage(s *discordgo.Session, channelID string) {
	commands := "Available commands are : " + strings.Join(commandList, ", ")
	_, err := s.ChannelMessageSend(channelID, commands)
	if err != nil {
		log.Printf("Failed to send help message: %s", err)
	}
}

func fetchAndSendImage(s *discordgo.Session, channelID, command string) {
	response, err := http.Get(NekoAPIURL + command)
	if err != nil {
		log.Printf("Failed to fetch image for command %s: %s", command, err)
		sendErrorMessage(s, channelID)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Failed to read response body for command %s: %s", command, err)
		sendErrorMessage(s, channelID)
		return
	}

	imageURL := gjson.GetBytes(body, "url").String()

	embed := &discordgo.MessageEmbed{
		Title:  strings.Title(command) + " Image",
		Image:  &discordgo.MessageEmbedImage{URL: imageURL},
		Color:  DefaultColor,
		Footer: &discordgo.MessageEmbedFooter{Text: "Powered by nekos.life"},
	}

	_, err = s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		log.Printf("Failed to send embed for command %s: %s", command, err)
		sendErrorMessage(s, channelID)
		return
	}
}

func sendErrorMessage(s *discordgo.Session, channelID string) {
	_, err := s.ChannelMessageSend(channelID, ErrorMessage)
	if err != nil {
		log.Printf("Failed to send error message: %s", err)
	}
}

func containsString(list []string, val string) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}
