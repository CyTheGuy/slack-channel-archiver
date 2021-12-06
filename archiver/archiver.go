package archiver

import (
	"os"
)

func (c *Config) Archiver() {

	// Get all public and unarchived channels
	allChannels, err := c.getConvos()
	if err != nil {
		os.Exit(0)
	}
	c.Log.Info().Msgf("A total of %v channels were pulled", len(allChannels))

	// Initate counter for # of channels messaged and # of channels archived so you can know if the targetLimit is reached
	messagedCount := 0
	archivedCount := 0

	for _, channel := range allChannels {

		// Check to make sure the messagedCount and archivedCount aren't past targetLimit as a safety net
		if messagedCount >= targetLimit && archivedCount >= targetLimit {
			c.Log.Info().Msgf("Archived & messaged count reached the target limit of %v - exiting", targetLimit)
			return
		}

		// Make sure the channel is at least creationDaysLimit days old and that it does not meet the criteria in the exceptions list or blacklist
		if !c.channelChecks(channel) {
			continue
		}

		c.Log.Debug().Msgf("Targeting %s - %s", channel.Name, channel.ID)
		// Get message history for the channel - filter out channel_leave/channel_join
		latestMessage, hasMessages := c.getConvoHistoryFiltered(channel)

		// If there is no latest message archive the channel
		if !hasMessages {
			c.Log.Info().Msgf("%s - %s has no message history", channel.Name, channel.ID)
			c.messageChannel(channel, nohistoryMessage)
			c.archiveChannel(channel)
			continue
		}

		// See when last channel message occured in days
		days, err := c.getLastMessageInDays(latestMessage.Msg)
		if err != nil {
			// If you could not convert lastMessage to days then log it, skip channel and continue
			c.Log.Error().Err(err).Msgf("Could not convert lastMessage to days for %s - %s - %s", channel.Name, channel.ID, err)
			continue
		}

		// If Over archiveInactiveDaysLimit Days and latestMessage.Text == archive_message and username == botName
		if days > archiveInactiveDaysLimit && latestMessage.Username == botName && archivedCount <= targetLimit {
			c.Log.Info().Msgf("%s has been inactive for over %v days after archiveMessage was sent - %v", channel.Name, archiveInactiveDaysLimit, days)
			c.archiveChannel(channel)
			archivedCount++
		} else if days > messageInactiveDaysLimit && messagedCount <= targetLimit {
			c.Log.Info().Msgf("%s has been inactive for over %v days - %v", channel.Name, messageInactiveDaysLimit, days)
			c.messageChannel(channel, archiveMessage)
			messagedCount++
		}
	}
}
