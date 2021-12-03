package archiver

import (
	"github.com/slack-go/slack"
)

func (c *Config) messageChannel(channel slack.Channel, message string) {

	if c.App.Stack == "production" {
		_, _, _, err := c.Slack.SendMessage(channel.ID, slack.MsgOptionText(message, false))
		if err != nil {
			c.Log.Error().Msgf("Failed to message %v - %s", channel.Name, err)
			return
		}

		c.Log.Info().Msgf("Messaged %s - %s", channel.Name, channel.ID)
		return
	}

	c.Log.Info().Msgf("Would have messaged %s - %s", channel.Name, channel.ID)
	return
}

func (c *Config) archiveChannel(channel slack.Channel) {

	if c.App.Stack == "production" {
		err := c.Slack.ArchiveConversation(channel.ID)
		if err != nil {
			c.Log.Error().Msgf("Failed to archive %v - %s", channel.Name, err)
			return
		}

		c.Log.Info().Msgf("Archived %s - %s", channel.Name, channel.ID)
		return
	}

	c.Log.Info().Msgf("Would have archived %s - %s", channel.Name, channel.ID)
	return
}
