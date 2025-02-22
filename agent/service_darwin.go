package main

import (
	"github.com/kardianos/service"
)

func logLocation() string {
	return "/Library/Application Support/FleetCmdr/fc_agent.log"
}

func getPlatformAgentConfig() *service.Config {

	return &service.Config{
		Name:             "FleetCmdrAgent",
		DisplayName:      "FleetCmdr Agent",
		Description:      "IT Fleet Command Platform Agent",
		Executable:       "/Applications/FleetCmdr/fc_agent",
		WorkingDirectory: "/Applications/FleetCmdr",
	}
}

func getPlatformUpdaterConfig() *service.Config {
	return &service.Config{
		Name:             "FleetCmdrUpdater",
		DisplayName:      "FleetCmdr Updater",
		Description:      "IT Fleet Command Platform Updater",
		Executable:       "/Applications/FleetCmdr/fc_updater",
		WorkingDirectory: "/Applications/FleetCmdr",
	}
}
