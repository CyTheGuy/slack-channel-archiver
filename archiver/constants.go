package archiver

const (
	targetLimit              = 1
	messageInactiveDaysLimit = 90
	archiveInactiveDaysLimit = 30
	creationDaysLimit        = 0
	botName                  = "ArchiveBot"
	secretsPath              = "./secrets/destination.yml"
	exceptionsPath           = "./settings/exceptions.yml"
	archiveMessage           = "Hello,\n\n This is an inactive channel. `It will be archived in 30 days if there are no new messages.`\n - All existing comments in the channel will be retained for easy browsing and everything can be read like any other channel.\n - If a need for this channel arises again, it can be unarchived by clicking `Channel Settings --> Unarchive`.\n\n Reach out to `#help-collabtools` if you have questions or want this channel added to the exceptions list."
	nohistoryMessage         = "There is no message history in this channel. So we are archiving it.\n - All existing comments in the channel will be retained for easy browsing and everything can be read like any other channel.\n - If a need for this channel arises again, it can be unarchived by clicking `Channel Settings --> Unarchive`."
)
