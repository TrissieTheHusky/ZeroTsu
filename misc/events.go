package misc

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"

	"github.com/r-anime/ZeroTsu/config"
)

// Periodic events such as Unbanning and RSS timer every 1 min
func StatusReady(s *discordgo.Session, e *discordgo.Ready) {
	err := s.UpdateStatus(0, "with her darling")
	if err != nil {
		_, err = s.ChannelMessageSend(config.BotLogID, err.Error())
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	for range time.NewTicker(1 * time.Minute).C {

		// Checks whether it has to post rss thread every 15 seconds
		RSSParser(s)

		////Goes through bannedUsers.json if it's not empty and unbans if needed
		//if BannedUsersSlice != nil {
		//	if len(BannedUsersSlice) != 0 {
		//
		//		t := time.Now()
		//
		//		for i := 0; i < len(BannedUsersSlice); i++ {
		//			difference := t.Sub(BannedUsersSlice[i].UnbanDate)
		//			if difference > 0 {
		//
		//				// Checks if user is in MemberInfo and assigns to user variable if true
		//				user, ok := MemberInfoMap[BannedUsersSlice[i].ID]
		//				if !ok {
		//					continue
		//				}
		//
		//				// Unbans user
		//				err := s.GuildBanDelete(config.ServerID, BannedUsersSlice[i].ID)
		//				if err != nil {
		//
		//					_, _ = s.ChannelMessageSend(config.BotLogID, err.Error())
		//				}
		//
		//				// Sends a message to bot-log
		//				_, _ = s.ChannelMessageSend(config.BotLogID, "User: " + user.Username + "#"+
		//					user.Discrim+ " has been unbanned.")
		//
		//				// Removes the user ban from bannedUsers.json
		//				BannedUsersSlice = append(BannedUsersSlice[:i], BannedUsersSlice[i+1:]...)
		//
		//				// Writes to bannedUsers.json
		//				BannedUsersWrite(BannedUsersSlice)
		//			}
		//		}
		//	}
		//}
	}
}

// Periodic 20min Events
func TwentyMinTimer(s *discordgo.Session, e *discordgo.Ready) {
	for range time.NewTicker(20 * time.Minute).C {

		// Writes emoji stats to disk
		_, err := EmojiStatsWrite(EmojiStats)
		if err != nil {
			_, err = s.ChannelMessageSend(config.BotLogID, err.Error())
			if err != nil {
				return
			}
			return
		}

		// Fetches all guild users
		guild, err := s.Guild(config.ServerID)
		if err != nil {
			_, err = s.ChannelMessageSend(config.BotLogID, err.Error())
			if err != nil {
				return
			}
			return
		}
		// Fetches all server roles
		roles, err := s.GuildRoles(config.ServerID)
		if err != nil {
			_, err = s.ChannelMessageSend(config.BotLogID, err.Error())
			if err != nil {
				return
			}
			return
		}
		// Updates optin role stat
		t := time.Now()
		MapMutex.Lock()
		for chas := range ChannelStats {
			if ChannelStats[chas].Optin {
				ChannelStats[chas].RoleCount[t.Format(DateFormat)] = GetRoleUserAmount(guild, roles, ChannelStats[chas].Name)
			}
		}
		MapMutex.Unlock()

		// Writes channel stats to disk
		_, err = ChannelStatsWrite(ChannelStats)
		if err != nil {
			_, err = s.ChannelMessageSend(config.BotLogID, err.Error())
			if err != nil {
				return
			}
			return
		}

		// Writes user gain stats to disk
		_, err = UserChangeStatsWrite(UserStats)
		if err != nil {
			_, err = s.ChannelMessageSend(config.BotLogID, err.Error())
			if err != nil {
				return
			}
			return
		}
	}
}

// Pulls the rss thread and prints it
func RSSParser(s *discordgo.Session) {

	var exists bool

	if len(ReadRssThreads) == 0 {
		return
	}

	// Pulls the feed from /r/anime and puts it in feed variable
	fp := gofeed.NewParser()
	fp.Client = &http.Client{Transport: &UserAgentTransport{http.DefaultTransport}, Timeout: time.Minute * 1}
	feed, err := fp.ParseURL("http://www.reddit.com/r/anime/new/.rss")
	if err != nil {
		return
	}
	fp.Client = &http.Client{}

	t := time.Now()
	hours := time.Hour * 16

	// Removes a thread if more than 16 hours have passed
	for p := 0; p < len(ReadRssThreadsCheck); p++ {
		// Calculates if it's time to remove
		dateRemoval := ReadRssThreadsCheck[p].Date.Add(hours)
		difference := t.Sub(dateRemoval)

		if difference > 0 {
			// Removes the fact that the thread had been posted already
			RssThreadsTimerRemove(ReadRssThreadsCheck[p].Thread, ReadRssThreadsCheck[p].Date)
		}
	}

	// Iterates through each feed item to see if it finds something from storage
	for i := 0; i < len(feed.Items); i++ {
		itemTitleLowercase := strings.ToLower(feed.Items[i].Title)
		itemAuthorLowercase := strings.ToLower(feed.Items[i].Author.Name)

		for j := 0; j < len(ReadRssThreads); j++ {
			exists = false
			storageAuthorLowercase := strings.ToLower(ReadRssThreads[j].Author)

			if strings.Contains(itemTitleLowercase, ReadRssThreads[j].Thread) &&
				strings.Contains(itemAuthorLowercase, storageAuthorLowercase) {

				for k := 0; k < len(ReadRssThreadsCheck); k++ {
					if ReadRssThreadsCheck[k].Thread == ReadRssThreads[j].Thread {
						exists = true
						break
					}
				}
				if !exists {
					// Writes to storage that the thread has been posted and posts rss in channel if no error via valid bool
					valid := RssThreadsTimerWrite(ReadRssThreads[j].Thread, t)
					if valid {
						_, err = s.ChannelMessageSend(ReadRssThreads[j].Channel, feed.Items[i].Link)
						if err != nil {
							_, _ = s.ChannelMessageSend(config.BotLogID, err.Error())
						}
					}
				}
			}
		}
	}
}

// Adds the voice role whenever a user joins the config voice chat
func VoiceRoleHandler(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {

	var roleIDString string

	// Saves program from panic and continues running normally without executing the command if it happens
	defer func() {
		if rec := recover(); rec != nil {
			_, err := s.ChannelMessageSend(config.BotLogID, rec.(string))
			if err != nil {
				return
			}
		}
	}()

	m, err := s.State.Member(v.GuildID, v.UserID)
	if err != nil {
		m, err = s.GuildMember(v.GuildID, v.UserID)
		if err != nil {
			_, err = s.ChannelMessageSend(config.BotLogID, err.Error())
			if err != nil {
				return
			}
			return
		}
		return
	}

	// Fetches role ID
	guildRoles, err := s.GuildRoles(config.ServerID)
	if err != nil {
		_, err = s.ChannelMessageSend(config.BotLogID, err.Error())
		if err != nil {
			return
		}
		return
	}
	for roleID := range guildRoles {
		if guildRoles[roleID].Name == "voice" {
			roleIDString = guildRoles[roleID].ID
		}
	}

	if v.ChannelID == config.VoiceChaID {
		// Adds role
		for _, role := range m.Roles {
			if role == roleIDString {
				return
			}
		}
		err = s.GuildMemberRoleAdd(v.GuildID, v.UserID, roleIDString)
		if err != nil {
			_, err = s.ChannelMessageSend(config.BotLogID, err.Error())
			if err != nil {
				return
			}
			return
		}
	} else {
		// Removes role
		for _, role := range m.Roles {
			if role == roleIDString {
				err := s.GuildMemberRoleRemove(v.GuildID, v.UserID, roleIDString)
				if err != nil {
					_, err = s.ChannelMessageSend(config.BotLogID, err.Error())
					if err != nil {
						return
					}
					return
				}
				break
			}
		}
	}
}