package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/r-anime/ZeroTsu/commands"
	"github.com/r-anime/ZeroTsu/config"
	"github.com/r-anime/ZeroTsu/misc"
)

var (
	BotID string
	goBot *discordgo.Session
)

//Starts Bot and its Handlers
func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {

		fmt.Println(err.Error())
		return
	}
	u, err := goBot.User("@me")

	if err != nil {

		fmt.Println(err.Error())
	}

	//Saves bot ID
	BotID = u.ID

	//Reads spoiler roles database at bot start
	misc.SpoilerRolesRead()

	// Reads filters.json from storage at bot start
	misc.FiltersRead()

	// Reads memberInfo.json from storage at bot start
	misc.MemberInfoRead()

	// Reads bannedUsers.json from storage at bot start
	misc.BannedUsersRead()

	// Reads ongoing votes from VoteInfo.json
	commands.VoteInfoRead()

	//Updates Playing Status
	goBot.AddHandler(misc.StatusReady)

	//Listens for a role deletion
	goBot.AddHandler(misc.ListenForDeletedRoleHandler)

	//Sorts spoiler roles alphabetically between opt-in dummy roles
	goBot.AddHandler(commands.SortRolesHandler)

	//Phrase Filter
	goBot.AddHandler(commands.FilterHandler)

	//React Filter
	goBot.AddHandler(commands.FilterReactsHandler)

	//Deletes non-whitelisted attachments
	goBot.AddHandler(commands.MessageAttachmentsHandler)

	// Abstraction of a command handler
	goBot.AddHandler(commands.HandleCommand)

	//MemberInfo
	//goBot.AddHandler(misc.OnMemberJoinGuild)
	//goBot.AddHandler(misc.OnMemberUpdate)

	//Whois Command
	//goBot.AddHandler(commands.WhoisHandler)

	//Unban Command
	//goBot.AddHandler(commands.UnbanHandler)

	//React Channel Set Command
	goBot.AddHandler(commands.SetReactChannelHandler)

	//React Channel Join Command
	goBot.AddHandler(commands.ReactJoinHandler)

	//React Channel Remove Command
	goBot.AddHandler(commands.ReactRemoveHandler)

	//React Channel Join View Command
	goBot.AddHandler(commands.ViewSetReactJoinsHandler)

	//React Channel Join Remove Command
	goBot.AddHandler(commands.RemoveReactJoinHandler)

	//RSS Parse Command
	goBot.AddHandler(commands.RSSHandler)

	//RSS Thread Check
	goBot.AddHandler(misc.RssThreadReady)

	//Channel Vote Timer
	goBot.AddHandler(commands.ChannelVoteTimer)

	//Verified Role and Cookie Map Expiry Deletion Handler
	//goBot.AddHandler(verification.VerifiedRoleAdd)
	//goBot.AddHandler(verification.VerifiedAlready)

	err = goBot.Open()

	if err != nil {

		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")
}
