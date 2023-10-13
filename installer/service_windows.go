package main

import "github.com/kardianos/service"

func getPlatformAgentConfig() *service.Config {
	return &service.Config{
		Name:             "FleetCmdrAgent",
		DisplayName:      "FleetCmdr Agent",
		Description:      "IT Fleet Command Agent Service",
		Executable:       "C:\\ProgramData\\FleetCmdr\\fc_agent.exe",
		WorkingDirectory: "C:\\Windows\\System32",
	}
}

func getPlatformInstallerConfig() *service.Config {
	return &service.Config{
		Name:             "FleetCmdrInstaller",
		DisplayName:      "FleetCmdr Installer",
		Description:      "IT Fleet Command Installer Service",
		Executable:       "C:\\ProgramData\\FleetCmdr\\fc_installer.exe",
		WorkingDirectory: "C:\\Windows\\System32",
	}
}
