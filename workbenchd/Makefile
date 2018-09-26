
install-agent:
	cp workbenchd.plist ~/Library/LaunchAgents/workbenchd.plist
	@launchctl unload ~/Library/LaunchAgents/workbenchd.plist
	launchctl load ~/Library/LaunchAgents/workbenchd.plist

clear-logs:
	rm logs/*