package commands

import (
	"fmt"
	"giftDetester/db"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
)

var timeoutDurationCommand = &discordgo.ApplicationCommandOption{
	Name:        "duration",
	Type:        discordgo.ApplicationCommandOptionSubCommand,
	Description: "Configure the timeout duration for compromised users",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "duration",
			Description: "The timeout duration in minutes",
			Required:    true,
		},
	},
}

func timeoutDurationHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Member.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {

		data := i.ApplicationCommandData()

		switch data.Options[0].Name {
		case "duration":
			data := data.Options[0]
			duration := data.Options[0].IntValue()

			// Can't send min_ and max_value yet, so have to check on our own
			if duration < 1 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						TTS: false,
						Embeds: []*discordgo.MessageEmbed{
							&discordgo.MessageEmbed{
								Type:        "rich",
								Title:       "Error",
								Description: "Cannot set the timeout duration to a negative number",
								Color:       0xff0000,
								Author: &discordgo.MessageEmbedAuthor{
									Name:    i.Member.User.Username,
									IconURL: i.Member.User.AvatarURL("24px"),
								},
							},
						},
					},
				})
				return
			}

			if err := db.SetServerOption(i.GuildID, "timeoutDuration", strconv.FormatInt(duration, 10)); err != nil {
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
							Title:       "timeout duration configured successfully",
							Description: fmt.Sprintf("You have changed the timeout duration to %d minutes", duration),
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
