package commands

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

var commands = []*discordgo.ApplicationCommand{

	&discordgo.ApplicationCommand{
		Name:        "phishing",
		Description: "configure phishing mitigation",
		Options: []*discordgo.ApplicationCommandOption{
			phishingActionCommand,
			logChannelCommand,
		},
	},
}

func RegisterCommands(dg *discordgo.Session) {
	for _, v := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}
	dg.AddHandler(commandHandler)
}

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	data := i.ApplicationCommandData()

	switch data.Name {
	case "phishing":
		data := i.ApplicationCommandData()

		switch data.Options[0].Name {
		case "action":
			phishingActionHandler(s, i)
		case "logs":
			logChannelHandler(s, i)
		}
	}

}
