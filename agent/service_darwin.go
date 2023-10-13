package main

import (
	"github.com/kardianos/service"
)

func getPlatformAgentConfig() *service.Config {

	return &service.Config{
		Name:             "FleetCmdrAgent",
		DisplayName:      "FleetCmdr Agent",
		Description:      "IT Fleet Command Platform",
		Executable:       "/Applications/FleetCmdr/fc_agent",
		WorkingDirectory: "/Applications/FleetCmdr",
	}
}

func getPlatformInstallerConfig() *service.Config {
	return &service.Config{
		Name:             "FleetCmdrAgent",
		DisplayName:      "FleetCmdr Agent",
		Description:      "IT Fleet Command Platform",
		Executable:       "/Applications/FleetCmdr/fc_installer",
		WorkingDirectory: "/Applications/FleetCmdr",
	}
}
