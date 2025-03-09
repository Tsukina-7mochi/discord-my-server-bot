package discordchannelwriter

import "github.com/bwmarrin/discordgo"

type DiscordChannelWriter struct {
	session   *discordgo.Session
	channelID string
}

func NewDiscordChannelWriter(session *discordgo.Session, channelID string) *DiscordChannelWriter {
	return &DiscordChannelWriter{
		session:   session,
		channelID: channelID,
	}
}

func (w *DiscordChannelWriter) Write(p []byte) (n int, err error) {
	_, err = w.session.ChannelMessageSend(w.channelID, string(p))
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
