package main

import "github.com/kardianos/service"

func logLocation() string {
	return "C:\\ProgramData\\FleetCmdr\\fc_agent.log"
}

func getPlatformAgentConfig() *service.Config {
	return &service.Config{
		Name:             "FleetCmdrAgent",
		DisplayName:      "FleetCmdr Agent",
		Description:      "IT Fleet Command Agent Service",
		Executable:       "C:\\ProgramData\\FleetCmdr\\fc_agent.exe",
		WorkingDirectory: "C:\\Windows\\System32",
	}
}

func getPlatformUpdaterConfig() *service.Config {
	return &service.Config{
		Name:             "FleetCmdrInstaller",
		DisplayName:      "FleetCmdr Installer",
		Description:      "IT Fleet Command Installer Service",
		Executable:       "C:\\ProgramData\\FleetCmdr\\fc_updater.exe",
		WorkingDirectory: "C:\\Windows\\System32",
	}
}
