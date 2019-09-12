package commands

import (
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/r-anime/ZeroTsu/config"
	"github.com/r-anime/ZeroTsu/misc"
)

var (
	commandMap = make(map[string]*command)
	aliasMap   = make(map[string]string)
	l          = log.New(os.Stderr, "cmds: ", log.LstdFlags|log.Lshortfile)
)

type command struct {
	execute      func(*discordgo.Session, *discordgo.Message)
	trigger      string
	aliases      []string
	desc         string
	deleteAfter  bool
	elevated     bool
	admin        bool
	category     string
	DMAble		 bool
}

func add(c *command) {
	commandMap[c.trigger] = c
	for _, alias := range c.aliases {
		aliasMap[alias] = c.trigger
	}
	l.Printf("Added command %s | %d aliases | %v category", c.trigger, len(c.aliases), c.category)
}

// HandleCommand handles the incoming message
func HandleCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
			log.Println("Recovery in HandleCommand with message: " + m.Content)
		}
	}()

	if m == nil {
		return
	}
	if m.Author == nil {
		return
	}
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	if m.Message == nil {
		return
	}
	if m.Message.Content == "" {
		return
	}

	// Handle guild command if it's coming from a server
	if m.GuildID != "" {
		handleGuild(s, m)
		return
	}

	// Parse the command
	var guildPrefix = "."
	if m.Message.Content[0:len(guildPrefix)] != guildPrefix {
		return
	}
	cmdTrigger := strings.Split(m.Content, " ")[0][len(guildPrefix):]
	cmdTrigger = strings.ToLower(cmdTrigger)
	cmd, ok := commandMap[cmdTrigger]
	if !ok {
		cmd, ok = commandMap[aliasMap[cmdTrigger]]
		if !ok {
			return
		}
	}

	// Allow only normal DMable commands
	if !cmd.DMAble {
		return
	}

	// Execute the command
	cmd.execute(s, m.Message)
}

// Handles a command from a guild
func handleGuild(s *discordgo.Session, m *discordgo.MessageCreate) {
	misc.MapMutex.Lock()
	if _, ok := misc.GuildMap[m.GuildID]; !ok {
		misc.InitDB(m.GuildID)
		misc.LoadGuilds()
	}

	var (
		guildPrefix       string
		guildVoteModule   bool
		guildWaifuModule  bool
		guildReactsModule bool
	)

	if _, ok := misc.GuildMap[m.GuildID]; !ok {
		misc.MapMutex.Unlock()
		return
	}
	guildPrefix = misc.GuildMap[m.GuildID].GuildConfig.Prefix
	guildVoteModule = misc.GuildMap[m.GuildID].GuildConfig.VoteModule
	guildWaifuModule = misc.GuildMap[m.GuildID].GuildConfig.WaifuModule
	guildReactsModule = misc.GuildMap[m.GuildID].GuildConfig.ReactsModule
	misc.MapMutex.Unlock()

	if m.Message.Content[0:len(guildPrefix)] != guildPrefix {
		return
	}
	cmdTrigger := strings.Split(m.Content, " ")[0][len(guildPrefix):]
	cmdTrigger = strings.ToLower(cmdTrigger)
	cmd, ok := commandMap[cmdTrigger]
	if !ok {
		cmd, ok = commandMap[aliasMap[cmdTrigger]]
		if !ok {
			return
		}
	}
	if cmd.trigger == "votecategory" ||
		cmd.trigger == "startvote" {
		if !guildVoteModule {
			return
		}
	}
	if cmd.category == "waifus" {
		if !guildWaifuModule {
			return
		}
	}
	if cmd.category == "reacts" {
		if !guildReactsModule {
			return
		}
	}
	if cmd.elevated {
		misc.MapMutex.Lock()
		if !HasElevatedPermissions(s, m.Author.ID, m.GuildID) {
			misc.MapMutex.Unlock()
			return
		}
		misc.MapMutex.Unlock()
	}
	if cmd.admin && m.Author.ID != config.OwnerID {
		mem, err := s.State.Member(m.GuildID, m.Author.ID)
		if err != nil {
			mem, err = s.GuildMember(m.GuildID, m.Author.ID)
			if err != nil {
				return
			}
		}

		admin, _ := MemberIsAdmin(s, m.GuildID, mem, discordgo.PermissionAdministrator)
		if !admin {
			return
		}
	}
	cmd.execute(s, m.Message)
	if cmd.deleteAfter {
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			return
		}
	}
}

// Handles a command from DMs
func handleDM(s *discordgo.Session, m *discordgo.MessageCreate) {

}

// Checks if a user has the admin permissions or is a privileged role
func HasElevatedPermissions(s *discordgo.Session, userID string, guildID string) bool {
	if userID == config.OwnerID {
		return true
	}

	mem, err := s.State.Member(guildID, userID)
	if err != nil {
		mem, err = s.GuildMember(guildID, userID)
		if err != nil {
			log.Println(err)
			return false
		}
	}

	isAdmin, err := MemberIsAdmin(s, guildID, mem, discordgo.PermissionAdministrator)
	if err != nil {
		log.Println(err.Error())
	}
	if isAdmin != false {
		return true
	}

	return HasPermissions(mem, guildID)
}

// Checks if member has admin permissions
func MemberIsAdmin(s *discordgo.Session, guildID string, mem *discordgo.Member, permission int) (bool, error) {
	// Iterate through the role IDs stored in member.Roles
	// to check permissions
	for _, roleID := range mem.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}

	return false, nil
}

// Checks if a user has a privileged role in a given server
func HasPermissions(m *discordgo.Member, guildID string) bool {
	for _, role := range m.Roles {
		for _, goodRole := range misc.GuildMap[guildID].GuildConfig.CommandRoles {
			if role == goodRole.ID {
				return true
			}
		}
	}
	return false
}
