package archiver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/slack-go/slack"
)

//
// Slack
//

func (c *Config) channelChecks(channel slack.Channel) bool {

	// Get channel creation date from now in days
	creationDays := c.getCreationInDays(channel)
	if creationDays < creationDaysLimit {
		c.Log.Debug().Msgf("%s is too new of a channel - skipping it", channel.Name)
		return false
	}

	// Check exceptions.yml blacklist for partial channel names to check by regex and channels for hard coded channel names to skip
	if c.checkRegexes(channel.Name, c.Exceptions.BlackList) || checkSliceForKey(channel.Name, c.Exceptions.Channels) {
		c.Log.Debug().Msgf("%s is in channel blacklist or exceptions - skipping it", channel.Name)
		return false
	}

	return true
}

func (c *Config) getCreationInDays(channel slack.Channel) float64 {

	c.Log.Debug().Msgf("Channel created time for %s - %s is %v", channel.Name, channel.ID, channel.Created.Time())
	durationFromNow := getDurationDiff(channel.Created.Time())
	c.Log.Debug().Msgf("Channel %s - %s was created %s ago", channel.Name, channel.ID, durationFromNow)

	createdDays := convertDurationToDays(durationFromNow)
	c.Log.Debug().Msgf("Duration from now - %v - converted to days is - %v", durationFromNow, createdDays)
	return createdDays
}

func (c *Config) getLastMessageInDays(message slack.Msg) (float64, error) {

	// Convert the last message timestamp to UTC
	c.Log.Debug().Msgf("Last message timestamp is - %v - attempting to convert it to UnixTime", message.Timestamp)
	timeUnix, err := c.getUnixTime(message.Timestamp)
	if err != nil {
		return 0, err
	}

	durationFromNow := getDurationDiff(timeUnix)
	c.Log.Debug().Msgf("Last message duration from now is - %s", durationFromNow)

	messagedDays := convertDurationToDays(durationFromNow)
	c.Log.Debug().Msgf("Last message duration converted to days is - %v", messagedDays)

	return messagedDays, nil
}

//
// Time
//

// Gets the duration difference between a unixtime and now
func getDurationDiff(t time.Time) time.Duration {
	return time.Now().Sub(t)
}

// Converts the duration from hours into days
func convertDurationToDays(t time.Duration) float64 {
	return t.Hours() / 24
}

func (c *Config) getUnixTime(ts string) (time.Time, error) {

	// Split string
	split := strings.Split(ts, ".")
	splitTS := split[0]
	c.Log.Debug().Msgf("Split %s to %s", ts, splitTS)

	// Convert string to int64
	intTS, err := strconv.ParseInt(splitTS, 10, 64)
	if err != nil {
		c.Log.Error().Err(err).Msgf("Could not convert %s to a int64", splitTS)
		return time.Time{}, err
	}
	c.Log.Debug().Msgf("Converted string %s to int %v", split, intTS)

	// Get Unix Time
	timeUnix := c.convertInt64ToUnixTime(intTS)
	c.Log.Debug().Msgf("Converted int %v to UnixTime %v", intTS, timeUnix)
	return timeUnix, nil
}

func (c *Config) convertInt64ToUnixTime(i int64) time.Time {

	timeUnix := time.Unix(i, 0)
	c.Log.Debug().Msgf("Unix Time is %s", timeUnix)
	return timeUnix
}

//
// Regex
//

// Checks to see if the item string is located anywhere in the slice of strings
func (c *Config) checkRegexes(item string, slice []string) bool {

	for _, str := range slice {

		// Generate the regex expression you want to check for
		expression := fmt.Sprintf("(?i)%s", str)
		// Parse regex expression
		regEx := regexp.MustCompile(expression)
		// Check regex expression

		if regEx.MatchString(item) {
			c.Log.Debug().Msgf("%s was a regex match for %s", item, str)
			return true
		}
	}
	return false
}

//
// Strings
//

func checkSliceForKey(item string, slice []string) bool {
	for _, s := range slice {
		if strings.Compare(s, item) == 0 {
			return true
		}
	}
	return false
}
