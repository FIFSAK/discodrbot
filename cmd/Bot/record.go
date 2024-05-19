package bot

import (
	"discordbot/cmd/S3"
	"fmt"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

func RecordAndUpload(v *discordgo.VoiceConnection, duration time.Duration) (string, error) {
	receivedAudio := make(chan *discordgo.Packet, 2)
	go dgvoice.ReceivePCM(v, receivedAudio)

	var audioBuffer [][]int16

	// Start recording
	err := v.Speaking(true)
	if err != nil {
		return "", err
	}
	defer v.Speaking(false)

	// Record for the specified duration
	recording := true
	go func() {
		time.Sleep(duration)
		recording = false
	}()

	for recording && voiceChannelMemberCount > 1 {
		p, ok := <-receivedAudio
		if !ok {
			return "", fmt.Errorf("failed to receive audio")
		}
		audioBuffer = append(audioBuffer, p.PCM)
	}

	leaveVoiceChannel()

	// Сохранить аудиобуфер в файл
	id := strconv.Itoa(len(S3.GetAllBucketObjects()) + 1)
	pcmFilename := "recorded_audio" + id + ".pcm"
	err = SaveAudioToFile(audioBuffer, pcmFilename)
	if err != nil {
		fmt.Printf("Error saving audio to file: %v\n", err)
		return "", err
	}

	// Конвертировать PCM в MP3
	mp3Filename := "recorded_audio" + id + ".mp3"
	err = ConvertPCMToMP3(pcmFilename, mp3Filename)

	if err != nil {
		fmt.Printf("Error converting PCM to MP3: %v\n", err)
		return "", err
	}

	// Загрузить MP3 файл в S3
	err = S3.UploadAudioFile(mp3Filename)
	if err != nil {
		fmt.Printf("Error uploading audio file to S3: %v\n", err)
		return "", err
	}

	fmt.Println("Audio recorded and uploaded successfully")
	err = DeleteFile(pcmFilename)
	if err != nil {
		fmt.Println("Error deleting MP3 file: %v\n", err)
		return "", err
	}
	err = DeleteFile(mp3Filename)
	if err != nil {
		fmt.Println("Error deleting MP3 file: %v\n", err)
		return "", err
	}
	return mp3Filename, nil
}
