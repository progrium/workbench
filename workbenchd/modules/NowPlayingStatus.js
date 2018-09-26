#!/usr/bin/env osascript -l JavaScript
ObjC.import("Cocoa");
 
function run(argv) {
    var app = $.NSApplication.sharedApplication;
    var statusBarItem = $.NSStatusBar.systemStatusBar.statusItemWithLength($.NSVariableStatusItemLength);

    ObjC.registerSubclass({
        name: "ProgramDelegate",
        methods: {
            "interval:": {
                types: ["void", ["id"]],
                implementation: function () {
                    try {
                        var spotify = Application("Spotify");
                        if (!spotify.running()) {
                            statusBarItem.button.title = "🎶";
                            return;
                        }
                        var state = spotify.playerState();
                        statusBarItem.button.appearsDisabled = (state != "playing");
                        statusBarItem.button.title = "🎶 "+spotify.currentTrack.name()+" - "+spotify.currentTrack.artist();
                    } catch (e) {
                        console.log(e);
                        statusBarItem.button.title = "🎶 (Error)";
                    }
                }
            }
        }
    });

    $.NSTimer.scheduledTimerWithTimeIntervalTargetSelectorUserInfoRepeats(0.5, $.ProgramDelegate.alloc.init, "interval:", null, true);
    app.run;
}

