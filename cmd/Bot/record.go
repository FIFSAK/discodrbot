package bot

import (
	"fmt"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"time"
)

func recordAndPlay(v *discordgo.VoiceConnection, duration time.Duration) ([][]int16, error) {
	receivedAudio := make(chan *discordgo.Packet, 2)
	go dgvoice.ReceivePCM(v, receivedAudio)

	var audioBuffer [][]int16

	// Start recording
	err := v.Speaking(true)
	if err != nil {
		return nil, err
	}
	defer v.Speaking(false)

	// Record for the specified duration
	recording := true
	go func() {
		time.Sleep(duration * time.Second)
		recording = false
	}()

	for recording {
		p, ok := <-receivedAudio
		if !ok {
			fmt.Println("record error")
			return nil, nil
		}
		audioBuffer = append(audioBuffer, p.PCM)
	}

	return audioBuffer, nil
}
