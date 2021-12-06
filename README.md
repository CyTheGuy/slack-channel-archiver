## How It Works:

Each time this runs it will message & archive up to the specified [targetLimit](https://github.com/CyTheGuy/slack-channel-archiver/blob/main/archiver/constants.go#L4) set. This makes sure if things go haywire that it doesn't archive beyond the limit given.

1. It gets all channels created over [creationDaysLimit](https://github.com/CyTheGuy/slack-channel-archiver/blob/main/archiver/constants.go#L7) days ago
2. Filters out channels that contain any of the black listed words or the name exactly matches a channel in this [exceptions list](https://github.com/CyTheGuy/slack-channel-archiver/blob/main/settings/exceptions.yml)
3. If the channel has no activity at all, it messages the channel with the [nohistoryMessage](https://github.com/CyTheGuy/slack-channel-archiver/blob/main/archiver/constants.go#L12) and archives the channel
4. If the channel has been inactive for the [messageInactiveDaysLimit](https://github.com/CyTheGuy/slack-channel-archiver/blob/main/archiver/constants.go#L5) it messages the channel with the [archiveMessage](hhttps://github.com/CyTheGuy/slack-channel-archiver/blob/main/archiver/constants.go#L11)
5. If there is no new activity (excluding channel leave/joins) for [archiveInactiveDaysLimit](https://github.com/CyTheGuy/slack-channel-archiver/blob/main/archiver/constants.go#L6) days after the archive warning was sent. It will archive the channel.

## Adding Specific Channels To The Exceptions List:

If you need to add a specific channel add it under "channels" in the [exceptions list](https://github.com/CyTheGuy/slack-channel-archiver/blob/main/settings/exceptions.yml)

## Adding Keywords To The Exceptions List:

If you need to target multiple channels with the same keyword (EX: outage). Add the keyword to the blacklist section in [exceptions list](https://github.com/CyTheGuy/slack-channel-archiver/blob/main/settings/exceptions.yml). If the channels includes the word anywhere in the name then it will not archive the channel.
