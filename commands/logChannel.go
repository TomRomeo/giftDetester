package commands

import (
	"fmt"
	"giftDetester/db"
	"github.com/bwmarrin/discordgo"
	"log"
)

var logChannelCommand = &discordgo.ApplicationCommandOption{
	Name:        "logs",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Description: "Configure a log channel for Gift De(tester)",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type: discordgo.ApplicationCommandOptionChannel,
			ChannelTypes: []discordgo.ChannelType{
				discordgo.ChannelTypeGuildText,
			},
			Name:        "logchannel",
			Description: "The logchannel",
			Required:    true,
		},
	},
}

func logChannelHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {

		data := i.ApplicationCommandData()

		switch data.Options[0].Name {
		case "logs":
			data := data.Options[0]
			choice := data.Options[0].ChannelValue(s)
			if err := db.SetServerOption(i.GuildID, "logChannel", choice.ID); err != nil {
				log.Printf("Got Error when trying to set server option:\n%s", err)
				return
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					TTS: false,
					Embeds: []*discordgo.MessageEmbed{
						&discordgo.MessageEmbed{
							Type:        "rich",
							Title:       "phishing log channel configured successfully",
							Description: fmt.Sprintf("You have changed the log channel for Gift De(tester) to %s", choice.Name),
							Color:       0x00ff00,
							Author: &discordgo.MessageEmbedAuthor{
								Name:    i.Member.User.Username,
								IconURL: i.Member.User.AvatarURL("24px"),
							},
						},
					},
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
