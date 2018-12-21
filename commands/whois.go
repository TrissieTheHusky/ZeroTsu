package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/r-anime/ZeroTsu/config"
	"github.com/r-anime/ZeroTsu/misc"
)

// Sends memberInfo user information to channel
func whoisCommand(s *discordgo.Session, m *discordgo.Message) {

	var (
		pastUsernames 		string
		pastNicknames 		string
		warnings      		string
		kicks         		string
		bans          		string
		unbanDate     		string
		splitMessage 		[]string
		isInsideGuild = 	true
		altIsInsideGuild =	true
	)

	messageLowercase := strings.ToLower(m.Content)
	commandStrings := strings.SplitN(messageLowercase, " ", 2)

	if len(commandStrings) < 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `"+config.BotPrefix+"whois [@user, userID, or username#discrim]`\n\n" +
			"Note: this command supports username#discrim where username contains spaces.")
		if err != nil {
			_, err = s.ChannelMessageSend(config.BotLogID, err.Error() + "\n" + misc.ErrorLocation(err))
			if err != nil {
				return
			}
			return
		}
		return
	}

	userID, err := misc.GetUserID(s, m, commandStrings)
	if err != nil {
		misc.CommandErrorHandler(s, m, err)
		return
	}

	// Fetches user from server if possible
	mem, err := s.State.Member(config.ServerID, userID)
	if err != nil {
		mem, err = s.GuildMember(config.ServerID, userID)
		if err != nil {
			isInsideGuild = false
		}
	}
	// Checks if user is in MemberInfo and assigns to user variable. Else initializes user.
	misc.MapMutex.Lock()
	user, ok := misc.MemberInfoMap[userID]
	if !ok {
		if mem == nil {
			_, err = s.ChannelMessageSend(m.ChannelID, "Error: User not found in memberInfo. Cannot whois until user joins the server.")
			if err != nil {
				_, err = s.ChannelMessageSend(config.BotLogID, err.Error() + "\n" + misc.ErrorLocation(err))
				if err != nil {
					misc.MapMutex.Unlock()
					return
				}
				misc.MapMutex.Unlock()
				return
			}
			misc.MapMutex.Unlock()
			return
		}

		// Initializes user if he doesn't exist and is in server
		misc.InitializeUser(mem)
		user = misc.MemberInfoMap[userID]
		misc.MemberInfoWrite(misc.MemberInfoMap)
	}
	misc.MapMutex.Unlock()

	// Puts past usernames into a string
	if len(user.PastUsernames) != 0 {
		for i := 0; i < len(user.PastUsernames); i++ {
			if len(pastUsernames) == 0 {
				pastUsernames = user.PastUsernames[i]
			} else {
				pastUsernames = pastUsernames + ", " + user.PastUsernames[i]
			}
		}
	} else {
		pastUsernames = "None"
	}

	// Puts past nicknames into a string
	if len(user.PastNicknames) != 0 {
		for i := 0; i < len(user.PastNicknames); i++ {
			if len(pastNicknames) == 0 {
				pastNicknames = user.PastNicknames[i]
			} else {
				pastNicknames = pastNicknames + ", " + user.PastNicknames[i]
			}
		}
	} else {
		pastNicknames = "None"
	}

	// Puts warnings into a slice
	if len(user.Warnings) != 0 {
		for i := 0; i < len(user.Warnings); i++ {
			if len(warnings) == 0 {
				// Converts index to string and appends warning
				iStr := strconv.Itoa(i + 1)
				warnings = user.Warnings[i] + " [" + iStr + "]"
			} else {
				// Converts index to string and appends new warning to old ones
				iStr := strconv.Itoa(i + 1)
				warnings = warnings + ", " + user.Warnings[i] + " [" + iStr + "]"

			}
		}
	} else {
		warnings = "None"
	}

	// Puts kicks into a slice
	if len(user.Kicks) != 0 {
		for i := 0; i < len(user.Kicks); i++ {
			if len(kicks) == 0 {
				// Converts index to string and appends kick
				iStr := strconv.Itoa(i + 1)
				kicks = user.Kicks[i] + " [" + iStr + "]"
			} else {
				// Converts index to string and appends new kick to old ones
				iStr := strconv.Itoa(i + 1)
				kicks = kicks + ", " + user.Kicks[i] + " [" + iStr + "]"
			}
		}
	} else {
		kicks = "None"
	}

	// Puts bans into a slice
	if len(user.Bans) != 0 {
		for i := 0; i < len(user.Bans); i++ {
			if len(bans) == 0 {
				// Converts index to string and appends ban
				iStr := strconv.Itoa(i + 1)
				bans = user.Bans[i] + " [" + iStr + "]"
			} else {
				// Converts index to string and appends new ban to old ones
				iStr := strconv.Itoa(i + 1)
				bans = bans + ", " + user.Bans[i] + " [" + iStr + "]"
			}
		}
	} else {
		bans = "None"
	}

	// Puts unban Date into a separate string variable
	unbanDate = user.UnbanDate
	if unbanDate == "" {
		unbanDate = "No Ban"
	}

	// Sets whois message
	message := "**User:** " + user.Username + "#" + user.Discrim + " | **ID:** " + user.ID +
		"\n\n**Past Usernames:** " + pastUsernames +
		"\n\n**Past Nicknames:** " + pastNicknames + "\n\n**Warnings:** " + warnings +
		"\n\n**Kicks:** " + kicks + "\n\n**Bans:** " + bans +
		"\n\n**Join Date:** " + user.JoinDate + "\n\n**Verification Date:** " +
		user.VerifiedDate

	// Sets reddit Username if it exists
	if user.RedditUsername != "" {
		message = message + "\n\n**Reddit Account:** " + "<https://reddit.com/u/" + user.RedditUsername + ">"
	} else {
		message += "\n\n**Reddit Account:** " + "None"
	}

	// Sets unban date if it exists
	if user.UnbanDate != "" {
		message += "\n\n**Unban Date:** " + user.UnbanDate
	}

	if !isInsideGuild {
		message += "\n\n**_User is not in the server._**"
	}

	// Alt check
	misc.MapMutex.Lock()
	alts := CheckAltAccountWhois(userID)

	// If there's more than one account with the same reddit username print a message
	if len(alts) > 1 {

		// Forms the alts string
		success := "\n\n**Alts:**\n"
		for _, altID := range alts {
			alt, err := s.State.Member(config.ServerID, altID)
			if err != nil {
				alt, err = s.GuildMember(config.ServerID, altID)
				if err != nil {
					altIsInsideGuild = false
				}
			}

			if altIsInsideGuild {
				success += alt.User.Mention() + "\n"
			} else {
				success += fmt.Sprintf("%v#%v | %v\n", misc.MemberInfoMap[altID].Username, misc.MemberInfoMap[altID].Discrim, altID)
				// Reset bool for future iterations
				altIsInsideGuild = true
			}
		}

		// Adds the alts to the whois message
		message += success
		alts = nil
	}
	misc.MapMutex.Unlock()

	// Checks if the message contains a mention and finds the actual name instead of ID
	message = misc.MentionParser(s, message)

	// Splits the message if it's over 1950 characters
	if len(message) > 1950 {
		splitMessage = misc.SplitLongMessage(message)
	}

	// Prints split or unsplit whois
	if splitMessage == nil {
		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			_, err = s.ChannelMessageSend(config.BotLogID, err.Error() + "\n" + misc.ErrorLocation(err))
			if err != nil {
				return
			}
			return
		}
		return
	}
	for i := 0; i < len(splitMessage); i++ {
		_, err := s.ChannelMessageSend(m.ChannelID, splitMessage[i])
		if err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, "Error: cannot send whois message.")
			if err != nil {
				_, err = s.ChannelMessageSend(config.BotLogID, err.Error() + "\n" + misc.ErrorLocation(err))
				if err != nil {
					return
				}
				return
			}
		}
	}
}

// Function that iterates through memberInfo.json and checks for any alt accounts for that ID. Whois version
func CheckAltAccountWhois(id string) []string {

	var alts []string

	// Stops func if target reddit username is nil
	if misc.MemberInfoMap[id].RedditUsername == "" {
		return nil
	}

	// Iterates through all users in memberInfo.json
	for _, user := range misc.MemberInfoMap {
		// Skips iteration if iteration reddit username is nil
		if user.RedditUsername == "" {
			continue
		}
		// Checks if the current user has the same reddit username as the entry parameter and adds to alts string slice if so
		if user.RedditUsername == misc.MemberInfoMap[id].RedditUsername {
			alts = append(alts, user.ID)
		}
	}
	if len(alts) > 1 {
		return alts
	} else {
		return nil
	}
}

// Displays all punishments for that user with timestamps and type of punishment
func showTimestampsCommand(s *discordgo.Session, m *discordgo.Message) {

	var message string

	messageLowercase := strings.ToLower(m.Content)
	commandStrings := strings.Split(messageLowercase, " ")

	if len(commandStrings) != 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Usage: `"+config.BotPrefix+"timestamps [@user, userID, or username#discrim]`")
		if err != nil {
			_, err = s.ChannelMessageSend(config.BotLogID, err.Error() + "\n" + misc.ErrorLocation(err))
			if err != nil {
				return
			}
			return
		}
		return
	}

	userID, err := misc.GetUserID(s, m, commandStrings)
	if err != nil {
		misc.CommandErrorHandler(s, m, err)
		return
	}

	// Fetches user from server if possible
	mem, err := s.State.Member(config.ServerID, userID)
	if err != nil {
		mem, err = s.GuildMember(config.ServerID, userID)
		if err != nil {
		}
	}
	// Checks if user is in MemberInfo and assigns to user variable. Else initializes user.
	misc.MapMutex.Lock()
	user, ok := misc.MemberInfoMap[userID]
	if !ok {
		if mem == nil {
			_, err = s.ChannelMessageSend(m.ChannelID, "Error: User not found in memberInfo. Cannot timestamp until they rejoin server.")
			if err != nil {
				_, err = s.ChannelMessageSend(config.BotLogID, err.Error() + "\n" + misc.ErrorLocation(err))
				if err != nil {
					misc.MapMutex.Unlock()
					return
				}
				misc.MapMutex.Unlock()
				return
			}
			misc.MapMutex.Unlock()
			return
		}

		// Initializes user if he doesn't exist and is in server
		misc.InitializeUser(mem)
		user = misc.MemberInfoMap[userID]
		misc.MemberInfoWrite(misc.MemberInfoMap)
	}

	// Check if timestamps exist
	if len(user.Timestamps) == 0 {
		_, err = s.ChannelMessageSend(m.ChannelID, "Error: No saved timestamps for that user.")
		if err != nil {
			_, err = s.ChannelMessageSend(config.BotLogID, err.Error() + "\n" + misc.ErrorLocation(err))
			if err != nil {
				misc.MapMutex.Unlock()
				return
			}
			misc.MapMutex.Unlock()
			return
		}
		misc.MapMutex.Unlock()
		return
	}

	// Formats message
	for _, timestamp := range user.Timestamps {
		timezone, displacement := timestamp.Timestamp.Zone()
		message += fmt.Sprintf("**%v:** `%v` - _%v %v %v, %v:%v:%v %v+%v_\n", timestamp.Type, timestamp.Punishment, timestamp.Timestamp.Day(),
			timestamp.Timestamp.Month(), timestamp.Timestamp.Year(), timestamp.Timestamp.Hour(), timestamp.Timestamp.Minute(), timestamp.Timestamp.Second(), timezone, displacement)
	}

	// Splits messsage if too long
	msgs := misc.SplitLongMessage(message)

	// Prints timestamps
	for index := range msgs {
		_, err = s.ChannelMessageSend(m.ChannelID, msgs[index])
		if err != nil {
			_, err = s.ChannelMessageSend(config.BotLogID, err.Error() + "\n" + misc.ErrorLocation(err))
			if err != nil {
				misc.MapMutex.Unlock()
				return
			}
			misc.MapMutex.Unlock()
			return
		}
	}


	misc.MapMutex.Unlock()
}

func init() {
	add(&command{
		execute:  whoisCommand,
		trigger:  "whois",
		desc:     "Pulls mod information about a user.",
		elevated: true,
		category: "misc",
	})
	add(&command{
		execute:  showTimestampsCommand,
		trigger:  "timestamp",
		aliases:  []string{"timestamps"},
		desc:     "Shows all punishments for a user and their timestamps.",
		elevated: true,
		category: "misc",
	})
}