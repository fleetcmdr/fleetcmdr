package main

import (
	"github.com/kardianos/service"
)

func logLocation() string {
	return "/Library/Application Support/FleetCmdr/fc_updater.log"
}

func getPlatformAgentConfig() *service.Config {

	return &service.Config{
		Name:             "FleetCmdrAgent",
		DisplayName:      "FleetCmdr Agent",
		Description:      "IT Fleet Command Platform",
		Executable:       "/Applications/FleetCmdr/fc_agent",
		WorkingDirectory: "/Applications/FleetCmdr",
	}
}