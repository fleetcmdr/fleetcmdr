package main

import "github.com/kardianos/service"

func getPlatformAgentConfig() *service.Config {
	return &service.Config{
		Name:        "FleetCmdrAgent",
		DisplayName: "FleetCmdr Agent",
		Description: "IT Fleet Command Platform",
		Executable:  "/usr/bin/local/fleetcmdr/fc_agent",
	}
}

func getPlatformInstallerConfig() *service.Config {
	return &service.Config{
		Name:        "FleetCmdrAgent",
		DisplayName: "FleetCmdr Agent",
		Description: "IT Fleet Command Platform",
		Executable:  "/usr/bin/local/fleetcmdr/fc_installer",
	}
}
