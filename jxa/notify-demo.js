#!/usr/bin/env osascript -l JavaScript
// http://www.cocoawithlove.com/2010/09/minimalist-cocoa-programming.html
ObjC.import("Cocoa");

function run(argv) {

    var msg = $.NSUserNotification.alloc.init
    msg.setTitle("Wake up")

    $.NSUserNotificationCenter.defaultUserNotificationCenter.scheduleNotification(msg)

    var cocoaApp = $.NSApplication.sharedApplication;
    cocoaApp.run;
}

