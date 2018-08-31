#!/usr/bin/env osascript -l JavaScript
ObjC.import("Cocoa");

function run(argv) {
    var app = $.NSApplication.sharedApplication;
    var statusBarItem = $.NSStatusBar.systemStatusBar.statusItemWithLength($.NSVariableStatusItemLength);
    var track = Application('Spotify').currentTrack

    ObjC.registerSubclass({
        name: "ProgramDelegate",
        methods: {
            "interval:": {
                types: ["void", ["id"]],
                implementation: function () {
                    statusBarItem.button.title = $("🎶 "+track.name()+" - "+track.artist());
                }
            }
        }
    });

    $.NSTimer.scheduledTimerWithTimeIntervalTargetSelectorUserInfoRepeats(0.5, $.ProgramDelegate.alloc.init, "interval:", null, true);
    app.run;
}

