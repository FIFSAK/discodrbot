package bot

import (
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

func echo(v *discordgo.VoiceConnection) {

	receivedAudio := make(chan *discordgo.Packet, 2)
	go dgvoice.ReceivePCM(v, receivedAudio)

	send := make(chan []int16, 2)
	go dgvoice.SendPCM(v, send)

	err := v.Speaking(true)
	if err != nil {
		return
	}
	defer v.Speaking(false)

	for {

		p, ok := <-receivedAudio
		if !ok {
			return
		}

		send <- p.PCM
	}
}
