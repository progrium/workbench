//#!/usr/bin/env osascript -l JavaScript
ObjC.import("Cocoa");

function run(argv) {
    app = Application.currentApplication();
    app.includeStandardAdditions = true

    var menu = $.NSMenu.new.autorelease;
    var menuItem = $.NSMenuItem.alloc.initWithTitleActionKeyEquivalent("Click meeee", "clicked:", "c").autorelease;
    menu.addItem(menuItem);

    ObjC.registerSubclass({
        name: "AppDelegate",
        methods: {
            "applicationDidFinishLaunching:": {
                types: ["void", ["id"]],
                implementation: function (notification) {
                    console.log("Hello")
                }
            },
            "setDockTile:": {
                types: ["void", ["id"]],
                implementation: function (tile) {
                    
                }
            },
            "applicationDockMenu:": {
                types: ["id", ["id"]],
                implementation: function (tile) {
                    return menu;
                }
            },
            "dockMenu:": {
                types: ["id", []],
                implementation: function () {
                    return menu;
                }
            },
            "clicked:": {
                types: ["void", ["id"]],
                implementation: function (notification) {
                    app.displayAlert('wowzer');
                }
            }
        }
    });

    var appDelegate = $.AppDelegate.alloc.init;
    var cocoaApp = $.NSApplication.sharedApplication;
    cocoaApp.setActivationPolicy($.NSApplicationActivationPolicyRegular);
    cocoaApp.delegate = appDelegate;
    cocoaApp.run;


}