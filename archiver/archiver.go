package archiver

import (
	"os"
)

// Archiver - starts the archiving of inactive channels workflow
func (c *Config) Archiver() {

	// Get all public and unarchived channels
	allChannels, err := c.getConvos()
	if err != nil {
		os.Exit(0)
	}

	// Initate counter for # of messaged and archived so you can know if they hit the targetLimit
	messagedCount := 0
	archivedCount := 0

	for _, channel := range allChannels {

		// Check to make sure the messagedCount and archivedCount aren't past targetLimit as a safety net
		if messagedCount >= targetLimit && archivedCount >= targetLimit {
			c.Log.Info().Msgf("Archived & messaged count reached the target limit of %v - exiting", targetLimit)
			return
		}

		// Make sure the channel is 60+ days old and that it does not meet the criteria in the exceptions list
		if !c.channelChecks(channel) {
			continue
		}

		c.Log.Debug().Msgf("Targeting %s", channel.Name)
		// Get message history for the channel - filter out channel_leave/channel_join
		lastestMessage, hasMessages := c.getConvoHistoryFiltered(channel)

		// If there is no latest message archive the channel
		if !hasMessages {
			c.Log.Info().Msgf("%s has no message history", channel.Name)
			c.messageChannel(channel, nohistoryMessage)
			c.archiveChannel(channel)
			continue
		}

		// See when last activity occured in days
		days, err := c.getLastMessageInDays(lastestMessage.Msg)
		if err != nil {
			c.Log.Error().Msgf("Could not get diff in days for %s - %s", lastestMessage.Msg.Timestamp, err)
			continue
		}

		// If Over 30 Days and lastestMessage.Text == archive_message and username == "ArchiveBot"
		if days > 30 && lastestMessage.Username == "ArchiveBot" && archivedCount <= targetLimit {
			c.Log.Info().Msgf("%s has been inactive for over 30 days after archiveMessage was sent - %v", channel.Name, days)
			c.archiveChannel(channel)
			archivedCount++
		} else if days > inactiveDaysLimit && messagedCount <= targetLimit {
			c.Log.Info().Msgf("%s has been inactive for over %v days - %v", channel.Name, inactiveDaysLimit, days)
			c.messageChannel(channel, archiveMessage)
			messagedCount++
		}
	}
}
