package archiver

import (
	"os"
	"time"

	"github.com/slack-go/slack"
)

//
// Channels
//

// Paginates through all public channels to return their information
func (c *Config) getConvos() ([]slack.Channel, error) {

	convos, cursor, err := c.Slack.GetConversations(&slack.GetConversationsParameters{
		Types:           []string{"public_channel"},
		Limit:           1000,
		ExcludeArchived: true})

	if err != nil {
		c.Log.Error().Msgf("Failed to get slack conversations - %s", err)
		return []slack.Channel{}, err
	}

	// Preform pagination if cursor determines there are additional channels to append
	for cursor != "" {
		temp, next, err := c.Slack.GetConversations(&slack.GetConversationsParameters{
			Types:           []string{"public_channel"},
			Limit:           1000,
			Cursor:          cursor,
			ExcludeArchived: true,
		})

		if err != nil {
			if rateLimitErr, ok := err.(*slack.RateLimitedError); ok {
				c.Log.Debug().Msg("Rate Limited.... Sleeping")
				time.Sleep(rateLimitErr.RetryAfter)
				// If retry time does not appear to be long enough you can add more sleep
				//time.Sleep(3000000000)
				continue
			} else {
				c.Log.Error().Msgf("Failed to get slack conversations - %s", err)
				return []slack.Channel{}, err
			}
		}

		cursor = next
		convos = append(convos, temp...)
	}

	return convos, nil
}

// Get a specific channels info by the channelID
func (c *Config) getConvoByID(channelID string) (*slack.Channel, error) {

	channel, err := c.Slack.GetConversationInfo(channelID, false)

	if err != nil {
		if rateLimitErr, ok := err.(*slack.RateLimitedError); ok {
			c.Log.Debug().Msg("Rate Limited.... Sleeping")
			time.Sleep(rateLimitErr.RetryAfter)
			return c.getConvoByID(channelID)
		}
		c.Log.Error().Msgf("Could not get channel info for %v - %s", channelID, err)
		return &slack.Channel{}, err

	}

	c.Log.Debug().Msgf("Got chanel info for %s", channelID)
	return channel, nil
}

//
// History
//

func (c *Config) getConvoHistoryFiltered(channel slack.Channel) (slack.Message, bool) {

	history, err := c.Slack.GetConversationHistory(
		&slack.GetConversationHistoryParameters{
			Limit:     200,
			ChannelID: channel.ID})

	if err != nil {
		if rateLimitErr, ok := err.(*slack.RateLimitedError); ok {
			c.Log.Debug().Msg("Rate Limited.... Sleeping")
			time.Sleep(rateLimitErr.RetryAfter)
			return c.getConvoHistoryFiltered(channel)
		}
		c.Log.Error().Msgf("Could not get filtered channel history for %v - %s", channel.Name, err)
		// Exit since if it returns an empty message struct it will assume this channel has no history and it will try to archive it
		os.Exit(0)
	}

	for _, message := range history.Messages {
		if message.SubType == "channel_leave" || message.SubType == "channel_join" {
			c.Log.Debug().Msgf("Message subtype is %v for %s - skipping it", message.SubType, message.Timestamp)
			continue
		}
		// Found latest message based on filters
		c.Log.Debug().Msgf("Got latest message in filtered channel history for %s", channel.Name)
		return message, true
	}

	// Check other pages for relevant history
	cursor := history.ResponseMetaData.NextCursor

	for cursor != "" {
		history, err := c.Slack.GetConversationHistory(
			&slack.GetConversationHistoryParameters{
				Limit:     200,
				Cursor:    cursor,
				ChannelID: channel.ID})

		if err != nil {
			if rateLimitErr, ok := err.(*slack.RateLimitedError); ok {
				c.Log.Debug().Msg("Rate Limited.... Sleeping")
				time.Sleep(rateLimitErr.RetryAfter)
				// Retry time does not appear to be long enough
				//time.Sleep(3000000000)
				continue
			} else {
				c.Log.Error().Msgf("Could not get filtered channel history for %v - %s", channel.Name, err)
				// Exit since if it returns an empty message struct it will assume this channel has no history and try to archive it
				os.Exit(0)
			}
		}

		for _, message := range history.Messages {
			if message.SubType == "channel_leave" || message.SubType == "channel_join" {
				c.Log.Debug().Msgf("Message subtype is %v for %s - skipping it", message.SubType, message.Timestamp)
				continue
			}
			// Found latest message based on filters
			c.Log.Debug().Msgf("Got latest message in filtered channel history for %s", channel.Name)
			return message, true
		}
		cursor = history.ResponseMetaData.NextCursor
	}

	c.Log.Error().Msgf("Could not find any latest message in filtered channel history for %s", channel.Name)
	return slack.Message{}, false
}
