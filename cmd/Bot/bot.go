package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	var BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	fmt.Println("Bot Token: ", BotToken) // выводим токен
	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal("Error creating session")
	}

	// add event handlers
	discord.AddHandler(newMessage)

	// open session
	err = discord.Open()
	if err != nil {
		log.Fatalf("Error opening session: %v", err)
	}
	defer discord.Close() // close session, after function termination

	// keep bot running until there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	// prevent bot responding to its own message
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// debug: print out the message details
	fmt.Printf("Message ID: %s\n", message.ID)
	fmt.Printf("Channel ID: %s\n", message.ChannelID)
	fmt.Printf("Author ID: %s\n", message.Author.ID)
	fmt.Printf("Content: %s\n", message.Content)
	fmt.Printf("Timestamp: %s\n", message.Timestamp)
	fmt.Printf("Embeds: %+v\n", message.Embeds)
	fmt.Printf("Attachments: %+v\n", message.Attachments)

	// respond to user message if it contains `!help` or `!bye`
	if message.Content == "!help" {
		discord.ChannelMessageSend(message.ChannelID, "Here is some help!")
	} else if message.Content == "!bye" {
		discord.ChannelMessageSend(message.ChannelID, "Goodbye!")
	} else {
		discord.ChannelMessageSend(message.ChannelID, "You said: "+message.Content)
	}
}
