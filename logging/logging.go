package logging

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func getLogChannel(s *discordgo.Session, guildID string) *discordgo.Channel {
	var logChannel *discordgo.Channel

	channels, _ := s.GuildChannels(guildID)
	for _, channel := range channels {
		if channel.Name == "gift-detester-log" {
			logChannel = channel
			break
		}
	}
	return logChannel
}
func LogAction(s *discordgo.Session, m *discordgo.Message, action string) {
	e := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       action,
		Description: "",
		Color:       0x7289da,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL:    m.Author.AvatarURL("24px"),
			Width:  24,
			Height: 24,
		},
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Message",
				Value:  m.Content,
				Inline: false,
			},
		},
	}
	logChannel := getLogChannel(s, m.GuildID)
	if logChannel != nil {
		_, _ = s.ChannelMessageSendEmbed(logChannel.ID, e)
	}
}

func SendError(s *discordgo.Session, m *discordgo.MessageCreate, description string, err error) {
	e := &discordgo.MessageEmbed{
		Type:        "rich",
		Title:       "Woops, an error occurred!",
		Description: fmt.Sprintf("%s:\n%s", description, err.Error()),
		Color:       0xff0000,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Is this not supposed to happen?\nYou can report bugs at",
		},
	}

	logChannel := getLogChannel(s, m.GuildID)
	if logChannel != nil {
		if _, err := s.ChannelMessageSendEmbed(logChannel.ID, e); err != nil {
			log.Println("Could not send error message:", err)
		}
	}
}
func NotifyUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	ch, _ := s.UserChannelCreate(m.Author.ID)

	guild, _ := s.Guild(m.GuildID)
	i, err := s.ChannelInviteCreate(guild.SystemChannelID, discordgo.Invite{
		Guild:     guild,
		MaxUses:   1,
		MaxAge:    86400,
		Temporary: false,
	})
	var maybeLink string
	if err != nil {
		SendError(s, m, "Could not create invite, missing permissions?", err)
	} else {
		maybeLink = fmt.Sprintf(" via [this link](https://discord.gg/%s)", i.Code)
	}
	e := []*discordgo.MessageEmbed{
		&discordgo.MessageEmbed{
			Type:        "rich",
			Title:       fmt.Sprintf("You have been kicked from %s", guild.Name),
			Description: "Your account has probably been compromised and sent multiple phishing links on this server.",
			Color:       0xff0000,
			Footer:      nil,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    s.State.User.Username,
				IconURL: s.State.User.AvatarURL("24px"),
			},

			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name: "What happened?",
					Value: `There is a popular scam going around of people getting access to someones discord and using it to send malicious links
						Most of the time this happens without the user even realizing it.`,
					Inline: false,
				},
			},
		},
		&discordgo.MessageEmbed{
			Type:  "rich",
			Title: "What can I do?",
			Description: `Change your Discord password. This will log you and any potential hackers out on all devices.
						Additionally, if you use the same password for other services like Spotify or Netflix, you should change them there as well.`,
			Color: 0x7289da,
		},
		&discordgo.MessageEmbed{
			Type:  "rich",
			Title: `"I want to join the server again!"`,
			Description: `Before rejoining the server, please reset your password. Sending phishing links multiple times might lead to a permanent ban.
						If you have changed your password, feel free to join the server again` + maybeLink,
			Color: 0x7289da,
		},
	}
	s.ChannelMessageSendEmbeds(ch.ID, e)
}
