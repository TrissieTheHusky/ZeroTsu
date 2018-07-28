package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/r-anime/ZeroTsu/config"
)

// Returns a message on "ping" to see if bot is alive
func pingCommand(s *discordgo.Session, m *discordgo.Message) {
	_, err := s.ChannelMessageSend(m.ChannelID, "Hmm? Do you want some honey, darling? Open wide~~")
	if err != nil {

		_, err = s.ChannelMessageSend(config.BotLogID, err.Error())
		if err != nil {

			return
		}
		return
	}
}

func init() {
	add(&command{
		execute:  pingCommand,
		trigger:  "ping",
		aliases:  []string{"pingme"},
		desc:     "Am I alive?",
		elevated: true,
	})
}