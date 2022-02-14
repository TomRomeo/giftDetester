package commands

import (
	"fmt"
	"giftDetester/db"
	"github.com/bwmarrin/discordgo"
	"log"
)

var phishingActionCommand = &discordgo.ApplicationCommandOption{
	Name:        "action",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Description: "Configure how to handle compromised users",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "action",
			Description: "The action to perform on a compromised user",
			Required:    true,
			Choices: []*discordgo.ApplicationCommandOptionChoice{
				&discordgo.ApplicationCommandOptionChoice{
					Name:  "kick",
					Value: "kick",
				},
				&discordgo.ApplicationCommandOptionChoice{
					Name:  "timeout",
					Value: "timout",
				},
			},
		},
	},
}

func phishingActionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Member.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {

		data := i.ApplicationCommandData()

		switch data.Options[0].Name {
		case "action":
			data := data.Options[0]
			choice := data.Options[0].StringValue()

			if err := db.SetServerOption(i.GuildID, "action", choice); err != nil {
				log.Printf("Got Error when trying to set server option:\n%s", err)
				return
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					TTS: false,
					Embeds: []*discordgo.MessageEmbed{
						&discordgo.MessageEmbed{
							Type:        "rich",
							Title:       "phishing action configured successfully",
							Description: fmt.Sprintf("You have changed the phishing action to %s", choice),
							Color:       0x00ff00,
							Author: &discordgo.MessageEmbedAuthor{
								Name:    i.Member.User.Username,
								IconURL: i.Member.User.AvatarURL("24px"),
							},
						},
					},
				},
			})
		}
	}
}
