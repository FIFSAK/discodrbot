package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"time"
)

var discord *discordgo.Session
var voiceConnection *discordgo.VoiceConnection

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	var BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	var err error

	// create a session
	discord, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal("Error creating session")
	}

	// add event handlers
	discord.AddHandler(channelCreate)
	discord.AddHandler(voiceStateUpdate)
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

	// Close voice connection when the bot is terminated
	if voiceConnection != nil {
		voiceConnection.Close()
	}
}

func channelCreate(s *discordgo.Session, c *discordgo.ChannelCreate) {
	// Check if the new channel is a voice channel
	if c.Type == discordgo.ChannelTypeGuildVoice {
		// Attempt to join the new voice channel
		err := joinVoiceChannel(s, c.GuildID, c.ID)
		if err != nil {
			fmt.Println("Error joining voice channel:", err)
		}
	}
}

func voiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if voiceConnection == nil {
		return
	}

	// Check if there are no users left in the voice channel
	channel, err := s.State.Channel(voiceConnection.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel:", err)
		return
	}

	// Получаем информацию о гильдии (сервере), к которой принадлежит канал
	guild, err := s.State.Guild(channel.GuildID)
	if err != nil {
		fmt.Println("Error getting guild:", err)
		return
	}
	voiceChannelMemberCount := 0
	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == voiceConnection.ChannelID {
			voiceChannelMemberCount++
		}
	}
	fmt.Println("Количество пользователей в голосовом канале:", voiceChannelMemberCount)
	recordStart := false
	if voiceChannelMemberCount > 0 {
		recordStart = true
		go func() {
			time.Sleep(7 * time.Second)
			voiceConnection.Close()
			err := voiceConnection.Disconnect()
			if err != nil {
				fmt.Println("Error disconnecting from voice channel:", err)
				return
			}
			recordStart = false
			fmt.Println("Отключен от голосового канала:", voiceConnection.ChannelID)
		}()
	}
	if recordStart && voiceChannelMemberCount == 1 {
		voiceConnection.Close()
		err := voiceConnection.Disconnect()
		if err != nil {
			fmt.Println("Error disconnecting from voice channel:", err)
			return
		}
		recordStart = false
		fmt.Println("Отключен от голосового канала:", voiceConnection.ChannelID)
	}
}

func joinVoiceChannel(s *discordgo.Session, guildID, voiceChannelID string) error {
	var err error
	voiceConnection, err = s.ChannelVoiceJoin(guildID, voiceChannelID, false, false)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к голосовому каналу: %v", err)
	}

	fmt.Println("Подключен к голосовому каналу:", voiceChannelID)
	return nil
}
