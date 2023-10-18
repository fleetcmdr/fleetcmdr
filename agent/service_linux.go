package main

import "github.com/kardianos/service"

func logLocation() string {
	return "/var/log/fleetcmdr/fc_agent.log"
}

func getPlatformAgentConfig() *service.Config {
	return &service.Config{
		Name:        "FleetCmdrAgent",
		DisplayName: "FleetCmdr Agent",
		Description: "IT Fleet Command Platform",
		Executable:  "/usr/bin/local/fleetcmdr/fc_agent",
	}
}

func getPlatformUpdaterConfig() *service.Config {
	return &service.Config{
		Name:        "FleetCmdrAgent",
		DisplayName: "FleetCmdr Agent",
		Description: "IT Fleet Command Platform",
		Executable:  "/usr/bin/local/fleetcmdr/fc_updater",
	}
}
