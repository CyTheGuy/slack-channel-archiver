package archiver

import (
	"time"
	"github.com/slack-go/slack"
)

//
// Slack
//

func (c *Config) channelChecks(channel slack.Channel) bool {

	// Check creation date and filter out channels created less than 60 days ago since they are too new to archive
	if c.getCreationInDays(channel) < creationDaysLimit {
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

	duration := getDurationDiff(channel.Created.Time())
	c.Log.Debug().Msgf("Created duration from now for %s is - %s", channel.Name, duration)

	createdDays := convertDurationToDays(duration)
	c.Log.Debug().Msgf("Created duration converted to days is - %v", createdDays)

	return createdDays
}

func (c *Config) getLastMessageInDays(message slack.Msg) (float64, error) {

	// Convert TS to UTC
	timeUnix, err := c.getUnixTime(message.Timestamp)
	if err != nil {
		return 0, err
	}

	duration := getDurationDiff(timeUnix)
	c.Log.Debug().Msgf("Last message duration from now is - %s", duration)

	messagedDays := convertDurationToDays(duration)
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

// Converts the duration into days
func convertDurationToDays(t time.Duration) (float64, error)

	c.Log.Debug().Msgf("Duration passed convertDurationToDays is - %v", t)
	d, err := time.ParseDuration(t)
	if err != nil {
		c.Log.Error().Err(err).Msgf("Could not convert %s to a int64", splitTS)
		return time.Time{}, err
	}

	days := d.Hours()/24
	c.Log.Debug().Msgf("Duration converted to days is - %v", days)

	return days, nil
}

func (c *Config) getUnixTime(ts string) (time.Time, error) {

	// Split string
	splitTS := strings.Split(ts, ".", 0)
	c.Log.Debug().Msgf("Split %s to %s", ts, splitTS)

	// Convert string to int64
	intTS, err := strconv.ParseInt(splitTS, 10, 64)
	c.Log.Debug().Msgf("Converted string %s to int %v", splitTS, intTS)

	if err != nil {
		c.Log.Error().Err(err).Msgf("Could not convert %s to a int64", splitTS)
		return time.Time{}, err
	}

	// Get Unix Time
	timeUnix := convertInt64ToUnixTime(intTS)
	c.Log.Debug().Msgf("Converted int %v to %v", intTS, timeUnix)

	// Return
	c.Log.Debug().Msgf("Unix Time is %s", timeUnix)
	return timeUnix, nil
}

func (c *Config) convertInt64ToUnixTime(i int64) (time.Time, error) {

	unitTime = time.Unix(i, 0)
	return 
}

//
// Regex
//

// Checks to see if the item string is located anywhere in the slice of strings
func (c *Config) checkRegexes(item string, slice []string) bool {
	for _, s := range slice {
		// Generate the regex expression you want to check for
		expression := fmt.Sprintf("(?i)%s", s)
		// Parse regex expression
		regEx := regexp.MustCompile(expression)
		// Check regex expression
		matched, err := regEx.MatchString(item)

		// If there was an error checking regexes then exit
		if err != nil {
			c.Log.Debug().Msgf("There was an error checking %s for %s using regex %s - %s" s, item, expression, err)
			os.Exit(0)
		}
	}
	return matched
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
