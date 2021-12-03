package main

import (
	"github.com/CyTheGuy/slack-channel-archiver/archiver"
)

func main() {

	config := archiver.ReadConfigs()
	config.Log.Info().Msgf("Running in - %v", config.App.Stack)
	config.Archiver()
	config.Log.Info().Msgf("Da End")
}
