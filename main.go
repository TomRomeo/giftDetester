package main

import (
	"database/sql"
	"fmt"
	"giftDetester/commands"
	"giftDetester/db"
	"giftDetester/logging"
	"giftDetester/util"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Could not load .env file")
	}

	// intialize db
	db.InitializeDb()
	log.Println("Connected to database")

	dg, err := discordgo.New("Bot " + os.Getenv("BOT_KEY"))
	if err != nil {
		log.Fatalf("An error occurred while trying to create the bot:\n%s", err)
	}
	if err = dg.Open(); err != nil {
		log.Fatalf("Could not establish a connection with Discord:\n%s", err)
	}
	log.Println("Successfully established a discord ws connection..")

	// register commands
	commands.RegisterCommands(dg, "")
	log.Println("Registered all Commands successfully...")

	dg.AddHandler(guildCreate)
	dg.AddHandler(guildDelete)
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageUpdate)
	log.Println("On the lookout for fake gift messages..")

	// graceful exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-c
	log.Println("Shutting down...")
	if err = dg.Close(); err != nil {
		log.Println("Failed to close Discord connection properly")
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// ignore own messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	var links []string

	// extract links from message
	links = append(links, util.ExtractLinks(m.Content)...)

	// also gather embed links
	for _, e := range m.Embeds {
		links = append(links, e.URL)
		links = append(links, util.ExtractLinks(e.Description)...)

		for _, f := range e.Fields {
			links = append(links, util.ExtractLinks(f.Value)...)
		}
	}

	// check each link
	for _, l := range links {
		if checkFakeGiftLink(l) {
			handleFakeGiftMessage(s, m.Message, l)
			break
		}
	}
}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {

	// ignore own messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	var links []string

	// extract links from message
	links = append(links, util.ExtractLinks(m.Message.Content)...)

	// also gather embed links
	for _, e := range m.Embeds {
		links = append(links, e.URL)
		links = append(links, util.ExtractLinks(e.Description)...)

		for _, f := range e.Fields {
			links = append(links, util.ExtractLinks(f.Value)...)
		}
	}

	// check each link
	for _, l := range links {
		if checkFakeGiftLink(l) {
			handleFakeGiftMessage(s, m.Message, l)
			break
		}
	}
}

func guildCreate(s *discordgo.Session, c *discordgo.GuildCreate) {

	if len(c.Roles) != 0 {

		var role *discordgo.Role
		// check each role
		for _, r := range c.Roles {
			if r.Name == "Gift (De)tester" {
				role = r
				break
			}
		}

		// check if has permissions
		p := role.Permissions
		pManageMessages := p & discordgo.PermissionManageMessages
		pCreateInstantInvite := p & discordgo.PermissionCreateInstantInvite
		pKickMembers := p & discordgo.PermissionKickMembers
		pModerateMembers := p & discordgo.PermissionModerateMembers
		pEmbedLinks := p & discordgo.PermissionEmbedLinks

		if !(pManageMessages == discordgo.PermissionManageMessages &&
			pCreateInstantInvite == discordgo.PermissionCreateInstantInvite &&
			pKickMembers == discordgo.PermissionKickMembers &&
			pModerateMembers == discordgo.PermissionModerateMembers &&
			pEmbedLinks == discordgo.PermissionEmbedLinks) {

			s.ChannelMessageSendEmbed(c.SystemChannelID, &discordgo.MessageEmbed{
				Type:  "rich",
				Title: "Did not receive all relevant permissions",
				Description: `I need all permissions in order to work correctly.
						Please add me again with the needed permissions`,
				Color: 0xff0000,
				Fields: []*discordgo.MessageEmbedField{
					&discordgo.MessageEmbedField{
						Name:   "Manage Messages",
						Value:  "To delete phishing messages",
						Inline: false,
					},
					&discordgo.MessageEmbedField{
						Name:   "Create Instant Invites",
						Value:  "To send members a rejoin link if they get kicked",
						Inline: false,
					},
					&discordgo.MessageEmbedField{
						Name:   "Kick members",
						Value:  "To kick members",
						Inline: false,
					},
					&discordgo.MessageEmbedField{
						Name:   "Moderate members",
						Value:  "To timeout members",
						Inline: false,
					},
					&discordgo.MessageEmbedField{
						Name:   "Embed links",
						Value:  "To Embed links",
						Inline: false,
					},
				},
			})
			s.GuildLeave(c.ID)

		} else {

			s.ChannelMessageSendEmbed(c.SystemChannelID, &discordgo.MessageEmbed{
				Type:        "rich",
				Title:       "Thanks for inviting me!",
				Description: `You can configure me with /phishing`,
				Color:       0x00ff00,
			})
		}
	}

}
func guildDelete(s *discordgo.Session, c *discordgo.GuildDelete) {
	if err := db.RemoveServer(c.ID); err != nil {
		log.Printf("Error occurred when trying to delete guild from db:\n%s", err)
	}
}

func checkFakeGiftLink(l string) bool {
	// we only care about the domain, not the path after
	u, _ := url.Parse(l)

	// firstly, return no if message is definitely from a Discord owned domain
	dcDomains := map[string]bool{
		"discord.com":       true,
		"discord.gg":        true,
		"discord.media":     true,
		"discordapp.com":    true,
		"discordapp.net":    true,
		"discordstatus.com": true,
		"discord.gift":      true,
	}
	if dcDomains[u.Host] {
		return false
	}

	// secondly, check if message contains similar spellings of discord
	spellings := []string{
		"discord",
		"dlscord",
		"d1scord",
		"discod",
		"disc0rd",
		"Diskord",
		"Hiscord",
		"Dhscord",
		"Dihcord",
		"Dishord",
		"Dischrd",
		"Discohd",
		"Discorh",
		"Discod1",
		"Iscord1",
		"Discor1",
		"Eiscord1",
		"Descord1",
		"Diecord1",
		"Diseord1",
		"Discerd1",
		"Discore1",
		"Diccord1",
		"Dyscord2",
		"Dscord2",
		"Dicord2",
		"Ddiscord2",
		"Discordd2",
		"Dizcord2",
		"Daiscord2",
		"Discorda2",
		"Discored2",
		"Deiscord2",
		"Discorde3",
		"Dissord3",
		"Discorrd3",
		"Disord3",
		"Discrd3",
		"Disscord3",
		"Dascord3",
		"Discorid3",
		"Diskaord3",
		"Niscord3",
		"Dnscord4",
		"Dincord4",
		"Disnord4",
		"Discnrd4",
		"Discond4",
		"Discorn4",
		"Iiscord4",
		"Diicord4",
		"Disiord4",
		"Discird4",
		"Discori",
	}
	for _, s := range spellings {
		if strings.Contains(u.Host, s) {
			return true
		}
	}

	// lastly, check if the link is similar to the official discord.gifts link
	similarity := util.CompareTwoLinks("discord.gift", u.Host)

	if similarity > 0.4 {
		return true
	}
	return false
}
func handleFakeGiftMessage(s *discordgo.Session, m *discordgo.Message, l string) {

	// gather what action to take
	var action string
	err := db.GetServerOption(m.GuildID, "action", &action)
	if err == sql.ErrNoRows || action == "" {
		action = "kick"
	}

	// firstly, notify user that they have been hacked
	logging.NotifyUser(s, m, action)

	// delete message
	if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
		logging.SendError(s, m, "Could not delete message, missing permissions?", err)
	} else {
		logging.LogAction(s, m, "Deleted Message")
	}

	if action == "kick" {

		if err := s.GuildMemberDeleteWithReason(m.GuildID, m.Author.ID, fmt.Sprintf("Fake gift link send: %s", l)); err != nil {
			logging.SendError(s, m, "Could not kick user, missing permissions?", err)
		} else {
			logging.LogAction(s, m, "Kicked User")
		}

	} else {

		// default is 24h
		duration := 1440
		db.GetServerOption(m.GuildID, "timeoutDuration", &duration)

		timeout := time.Now().Add(time.Duration(duration) * time.Minute)

		if err := s.GuildMemberTimeout(m.GuildID, m.Author.ID, &timeout); err != nil {
			logging.SendError(s, m, "Could not timeout user, missing permissions?", err)
		} else {
			logging.LogAction(s, m, "Timeouted User")
		}

	}

}
