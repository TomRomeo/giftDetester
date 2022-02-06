package commands

import (
	"fmt"
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

	if i.Member.Permissions&discordgo.PermissionAdministrator == 0 {

		data := i.ApplicationCommandData()

		switch data.Options[0].Name {
		case "action":
			data := data.Options[0]
			choice := data.Options[0].StringValue()
			switch choice {
			case "kick":
				log.Println("kick action")
			case "timeout":
				log.Println("timeout action")
			}
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				TTS: false,
				Embeds: []*discordgo.MessageEmbed{
					&discordgo.MessageEmbed{
						Type:        "rich",
						Title:       "phishing action configured successfully",
						Description: fmt.Sprintf("You have changed the phishing action to %s", i.ApplicationCommandData().Options[0].Name),
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
