package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var discord *discordgo.Session
var voiceConnection *discordgo.VoiceConnection
var isRecording bool
var voiceChannelMemberCount int

func Run() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}
	var BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	var err error

	// Create a session
	discord, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		fmt.Println("Error creating session:", err)
		return
	}

	// Add event handlers
	discord.AddHandler(channelCreate)
	discord.AddHandler(voiceStateUpdate)

	// Open session
	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening session:", err)
		return
	}
	defer discord.Close() // Close session after function termination

	// Keep bot running until there is no os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Close voice connection when the bot is terminated
	if voiceConnection != nil {
		voiceConnection.Close()
		fmt.Println("Voice connection closed during shutdown")
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

	// Получаем информацию о текущем голосовом канале, к которому подключен бот
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

	voiceChannelMemberCount = 0
	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == voiceConnection.ChannelID {
			voiceChannelMemberCount++
		}
	}
	fmt.Println("Количество пользователей в голосовом канале:", voiceChannelMemberCount)

	if voiceChannelMemberCount > 1 && !isRecording {
		isRecording = true
		go RecordAndUpload(voiceConnection, 10*time.Second, voiceChannelMemberCount) // Установите время записи в 10 секунд
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

func leaveVoiceChannel() {

	if voiceConnection != nil {
		voiceConnection.Close()
		fmt.Println("Отключен от голосового канала:", voiceConnection.ChannelID)
		isRecording = false
		err := voiceConnection.Disconnect()
		if err != nil {
			fmt.Printf("Disconnect error: %v\n", err)
			return
		}
	}
}
