
build:
	gcc -framework Cocoa -x objective-c -o Minimal minimal.m
	osacompile -l JavaScript -o Demo.app dock-menu.js
	mv Demo.app/Contents/MacOS/applet Demo.app/Contents/MacOS/applet2
	mv Minimal Demo.app/Contents/MacOS/applet
	cp Info.plist Demo.app/Contents/Info.plist

plugin:
	gcc -dynamiclib -framework Cocoa -x objective-c -o DockTileDemo plugin.m
	mv DockTileDemo Demo.app/Contents/PlugIns/DockTileDemo.docktileplugin/Contents/MacOS/DockTileDemo