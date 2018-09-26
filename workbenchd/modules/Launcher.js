#!/usr/bin/env osascript -l JavaScript
ObjC.import("Cocoa");
 
function run(argv) {
    var a = Application.currentApplication();
    a.includeStandardAdditions = true;
    a.displayNotification('The file has been converted',{ withTitle: ''})
    console.log("Running Launcher...");
    var app = $.NSApplication.sharedApplication;
    var statusBarItem = $.NSStatusBar.systemStatusBar.statusItemWithLength($.NSVariableStatusItemLength);
    statusBarItem.button.title = "🖥️";

    var menu = $.NSMenu.new.autorelease;
    var quitMenuItem = $.NSMenuItem.alloc.initWithTitleActionKeyEquivalent("Quit", "terminate:", "q").autorelease;
    menu.addItem(quitMenuItem);
    statusBarItem.menu = menu; 

    ObjC.registerSubclass({
        name: "ProgramDelegate",
        methods: {
            "interval:": {
                types: ["void", ["id"]],
                implementation: function () {
                    
                }
            }
        }
    });

    //$.NSTimer.scheduledTimerWithTimeIntervalTargetSelectorUserInfoRepeats(0.5, $.ProgramDelegate.alloc.init, "interval:", null, true);
    app.run;
}
