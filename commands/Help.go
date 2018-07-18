package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/r-anime/ZeroTsu/config"
	"github.com/r-anime/ZeroTsu/misc"
)

//Prints two versions of help depending on whether the user is a mod or not
func HelpHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, config.BotPrefix) {
		mem, err := s.State.Member(config.ServerID, m.Author.ID)
		if err != nil {
			mem, err = s.GuildMember(config.ServerID, m.Author.ID)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}

		if m.Author.ID == config.BotID {
			return
		}

		//Puts the command to lowercase
		messageLowercase := strings.ToLower(m.Content)

		//Checks for mod
		if misc.HasPermissions(mem) {
			if messageLowercase == config.BotPrefix+"help" {

				//Help message if user is a mod
				successMod := "`" + config.BotPrefix + "about` | Shows information about me \n " +
					"`" + config.BotPrefix + "filters` | Shows all current filters. \n " +
					"`" + config.BotPrefix + "addfilter [filter]` | Adds a normal or regex word to the filter. \n " +
					"`" + config.BotPrefix + "removefilter [filter]` | Removes a word from the filter. \n " +
					"`" + config.BotPrefix + "avatar [@mention or user ID]` | Returns user avatar URL and image embed. \n " +
					//"`" + config.BotPrefix + "convert [number][measurement] to [measurement]` | Converts measurement units (celsius, fahrenheit, kelvin, mm, cm, inch, foot, meter, km and mile.) \n " +
					"`" + config.BotPrefix + "create [name] [airing or general; defaults to opt-in] [category ID] [description; must have at least one other non-name parameter]` | Creates a channel and role of the same name. Do not start name with hyphens. \n " +
					"`" + config.BotPrefix + "help` | Lists commands and their usage. \n " +
					"`" + config.BotPrefix + "joinchannel [channel name]` | Joins an opt-in channel. `" + config.BotPrefix + "join` works too. \n " +
					"`" + config.BotPrefix + "leavechannel [channel name]` | Leaves an opt-in channel. `" + config.BotPrefix + "leave` works too. \n " +
					"`" + config.BotPrefix + "lock` | Locks a non-mod channel. Takes a few seconds only if the channel has no custom mod permissions set. \n " +
					"`" + config.BotPrefix + "unlock` | Unlocks a non-mod channel. \n " +
					"`" + config.BotPrefix + "ping` | Returns Pong message. \n " +
					"`" + config.BotPrefix + "say [message]` | Sends a message from the bot. \n "

				_, err = s.ChannelMessageSend(m.ChannelID, successMod)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				//Help message 2 if user is a mod
				successMod = "`" + config.BotPrefix + "setreactjoin [messageID] [role] [emote]` | Sets a specific message's emote to give those reacted a role. \n " +
					"`" + config.BotPrefix + "removereactjoin [messageID] OPTIONAL[emote]` | Removes the set react emote join from an entire message or only a specific emote of that message. \n " +
					"`" + config.BotPrefix + "viewreacts` | Prints out all current set message react emote joins. \n " +
					"`" + config.BotPrefix + "viewrss` | Prints out all currently set rss thread post. \n " +
					"`" + config.BotPrefix + "setrss OPTIONAL[/u/author] [thread name]` | Set a thread name which it'll look for in /new by the author (default /u/AutoLovepon) and then post that thread in the channel this command was executed in. \n " +
					"`" + config.BotPrefix + "removerss` | Remove a thread name from a previously set rss command. \n " +
					"`" + config.BotPrefix + "sortcategory [category name or ID]` | Sorts all channels within given category alphabetically. \n " +
					"`" + config.BotPrefix + "sortroles` | Sorts spoiler roles created with the create command between opt-in dummy roles alphabetically. Freezes server for a few seconds. Use preferably with large batches.\n" +
					"`" + config.BotPrefix + "startvote OPTIONAL[required votes] OPTIONAL[type] [name] OPTIONAL[categoryID] + OPTIONAL[description]` | Starts a reaction vote in the channel the command is in. " +
					"Creates and sorts the channel if successful. Required votes are how many non-bot reacts are needed for channel creation(default 7). Types are airing,general and optin(default)." +
					"CategoryID is what category to put the channel in and sort alphabetically. Description is the channel description but NEEDS a categoryID to work.\n"

				_, err = s.ChannelMessageSend(m.ChannelID, successMod)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		} else {

			if messageLowercase == config.BotPrefix+"help" {

				//Help message if user is not a mod
				successUser := "`" + config.BotPrefix + "about` | Shows information about me \n " +
					"`" + config.BotPrefix + "avatar [@mention or user ID]` | Returns user avatar URL and image embed. \n " +
					//"`" + config.BotPrefix + "convert [number][measurement] to [measurement]` | Converts measurement units (celsius, fahrenheit, kelvin, mm, cm, inch, foot, meter, km and mile.) \n " +
					"`" + config.BotPrefix + "help` | Lists commands and their usage. \n " +
					"`" + config.BotPrefix + "joinchannel [channel name]` | Joins an opt-in channel. `" + config.BotPrefix + "join` works too. \n " +
					"`" + config.BotPrefix + "leavechannel [channel name]` | Leaves an opt-in channel. `" + config.BotPrefix + "leave` works too. \n "

				_, err = s.ChannelMessageSend(m.ChannelID, successUser)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		}
	}
}
